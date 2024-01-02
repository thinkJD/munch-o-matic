package core

import (
	"fmt"
	"munch-o-matic/client"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
)

type jobStatus struct {
	JobId string
	Msg   string
	Err   error
}

type jobNotification struct {
	JobId    string
	Title    string
	Template string
	Data     interface{}
}

type Daemon struct {
	chron            *cron.Cron
	cli              *client.RestClient
	cfg              Config
	statusChan       chan jobStatus
	notificationChan chan jobNotification
}

func NewDaemon(Cfg Config, Cli *client.RestClient) (*Daemon, error) {
	retVal := Daemon{
		statusChan:       make(chan jobStatus),
		notificationChan: make(chan jobNotification),
	}

	err := ValidateConfig(Cfg)
	if err != nil {
		return &retVal, err
	}

	retVal.cli = Cli
	retVal.cfg = Cfg
	retVal.chron = cron.New(cron.WithSeconds())

	return &retVal, nil
}

func (d Daemon) Run() error {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	for _, job := range d.cfg.DaemonConfiguration.Jobs {
		d.AddJob(job)
	}

	// Handle notifications
	go func() {
		for msg := range d.notificationChan {
			// Skip if the status notification is disabled
			if !d.cfg.DaemonConfiguration.Notification.Status.Enabled {
				continue
			}

			SendTemplateNotification(
				d.cfg.DaemonConfiguration.Notification.Status.Topic,
				msg.Title, msg.Template, msg.Data)
		}
	}()

	// Handle status updates and errors
	go func() {
		for msg := range d.statusChan {
			if msg.Err != nil {
				err := fmt.Errorf("error in job=%s\t%w", msg.JobId, msg.Err)
				if d.cfg.DaemonConfiguration.Notification.Error.Enabled {
					SendNotification(d.cfg.DaemonConfiguration.Notification.Error.Topic,
						"Error", err.Error())
				}
			}
			fmt.Printf("%s:\t%s\n", msg.JobId, msg.Msg)
		}
	}()

	// Run metrics server
	go func() {
		http.Handle("/metrics", GetPrometheusHandler())
		fmt.Println("Metrics server is running on :9090")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			fmt.Println("Failed to start the metrics server:", err)
		}
	}()

	// Handle signals and clean up
	go func() {
		<-stopChan // wait for interrupt signal
		fmt.Println("Shutting down...")

		d.chron.Stop()
		close(d.statusChan)

		fmt.Println("Shutdown complete")
		os.Exit(0)
	}()

	// Start cron jobs
	d.chron.Start()
	fmt.Println("Next job execution: ", d.chron.Entries()[0].Next)

	// Block here until interrupted
	select {}
}

func (d *Daemon) AddJob(Job Job) error {
	switch Job.Type {
	case "CheckBalance":
		minBalance, ok1 := Job.Params["minbalance"].(int)
		template, ok2 := Job.Params["template"].(string)
		if !ok1 || !ok2 {
			return fmt.Errorf("invalid parameter types for CheckBalance")
		}

		_, err := d.chron.AddFunc(Job.Schedule, d.sendLowBalanceNotification(minBalance, template))
		if err != nil {
			return fmt.Errorf("error adding job: %w", err)
		}

	case "Order":
		strategy, ok1 := Job.Params["strategy"].(string)
		weeks, ok2 := Job.Params["weeks"].(int)
		if !ok1 || !ok2 {
			return fmt.Errorf("invalid parameter types for Order")
		}
		_, err := d.chron.AddFunc(Job.Schedule, d.orderFood(strategy, weeks))
		if err != nil {
			return fmt.Errorf("error adding job: %w", err)
		}

	case "UpdateMetrics":
		_, err := d.chron.AddFunc(Job.Schedule, d.updateMetrics())
		if err != nil {
			return fmt.Errorf("error adding job: %w", err)
		}

	default:
		return fmt.Errorf("%v is not a valid type", Job.Type)
	}
	return nil
}

// Jobs
// ##########

func (d Daemon) orderFood(Strategy string, WeeksInAdvance int) func() {
	return func() {
		var jobId = "orderFood"
		menu, err := d.cli.GetMenuWeeks(WeeksInAdvance)
		if err != nil {
			d.statusChan <- jobStatus{JobId: jobId, Err: fmt.Errorf("trouble getting menu: %v", err.Error())}
		}

		dishes, err := ChooseDishesByStrategy(Strategy, menu)
		if err != nil {
			d.statusChan <- jobStatus{JobId: jobId, Err: fmt.Errorf("Trouble choosing dish: %v", err.Error())}
		}
		fmt.Println(dishes)
		/*
			for _, dish := range dishes {
				err := Cli.OrderDish(dish.OrderId, false)
				if err != nil {
					ch <- fmt.Sprintf("Could not order")
				}
			}
		*/
		d.statusChan <- jobStatus{JobId: jobId, Msg: fmt.Sprintf("Food ordered with %v strategy, for %v weeks", Strategy, WeeksInAdvance)}
	}
}

func (d Daemon) sendLowBalanceNotification(MinBalance int, Template string) func() {
	return func() {
		var jobId = "sendLowBalanceNotification"
		d.statusChan <- jobStatus{JobId: jobId, Msg: "Checking account balance"}

		user, err := d.cli.GetUser()
		if err != nil {
			d.statusChan <- jobStatus{JobId: jobId, Err: fmt.Errorf("trouble getting user details: %v", err.Error())}
		}

		if user.User.Customer.AccountBalance.Amount <= MinBalance {
			d.statusChan <- jobStatus{JobId: jobId, Msg: "Account balance below minimum"}
			// Send notification
			d.notificationChan <- jobNotification{
				JobId:    jobId,
				Title:    "Low Balance",
				Template: Template,
				Data:     user}
		}

	}
}

func (d Daemon) updateMetrics() func() {
	return func() {
		var jobId = "updateMetrics"
		d.statusChan <- jobStatus{JobId: jobId, Msg: "Update metrics"}
		menuWeeks, err := d.cli.GetMenuWeeks(4)
		if err != nil {
			d.statusChan <- jobStatus{JobId: jobId, Err: fmt.Errorf("could not load dishes: %w", err)}
			return
		}

		for _, dishes := range menuWeeks {
			for _, dish := range dishes {
				if dish.Dummy {
					continue
				}
				UpdateOrdersPlaced(dish.OrderId, dish.Dish.Name, dish.Orders)
			}
		}

		user, err := d.cli.GetUser()
		if err != nil {
			d.statusChan <- jobStatus{JobId: jobId, Err: fmt.Errorf("could not load user: %w", err)}
			return
		}
		UpdateAccountBalance(user.User.ID, user.User.FirstName, user.User.Customer.AccountBalance.Amount)

		totalPayed := 0
		for _, i := range user.User.Customer.Bookings {
			totalPayed += i.BookingPrice
		}
		UpdatePaymentsTotal(user.User.ID, user.User.FirstName, totalPayed)
	}
}
