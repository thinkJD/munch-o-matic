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

type Daemon struct {
	d     *Daemon
	chron *cron.Cron
	cli   *client.RestClient
	cfg   Config
}

func NewDaemon(Cfg Config, Cli *client.RestClient) (*Daemon, error) {
	retVal := Daemon{}

	err := ValidateConfig(Cfg)
	if err != nil {
		return &retVal, err
	}

	retVal.cli = Cli
	retVal.cfg = Cfg
	retVal.chron = cron.New(cron.WithSeconds())

	return &retVal, nil
}

func (d *Daemon) AddJob(StatusChan chan jobStatus, Job Job) error {
	switch Job.Type {
	case "CheckBalance":
		topic, ok1 := Job.Params["topic"].(string)
		minBalance, ok2 := Job.Params["minbalance"].(int)
		template, ok3 := Job.Params["template"].(string)
		if !ok1 || !ok2 || !ok3 {
			return fmt.Errorf("invalid parameter types for CheckBalance")
		}

		_, err := d.chron.AddFunc(Job.Schedule, d.sendLowBalanceNotification(StatusChan, minBalance, topic, template))
		if err != nil {
			return fmt.Errorf("error adding job: %w", err)
		}

	case "Order":
		strategy, ok1 := Job.Params["strategy"].(string)
		weeks, ok2 := Job.Params["weeks"].(int)
		if !ok1 || !ok2 {
			return fmt.Errorf("invalid parameter types for Order")
		}
		_, err := d.chron.AddFunc(Job.Schedule, d.orderFood(StatusChan, strategy, weeks))
		if err != nil {
			return fmt.Errorf("error adding job: %w", err)
		}

	case "UpdateMetrics":
		_, err := d.chron.AddFunc(Job.Schedule, d.updateMetrics(StatusChan))
		if err != nil {
			return fmt.Errorf("error adding job: %w", err)
		}

	default:
		return fmt.Errorf("%v is not a valid type", Job.Type)
	}
	return nil
}

type jobStatus struct {
	JobId string
	Msg   string
	Err   error
}

func (d Daemon) Run() error {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	statusChan := make(chan jobStatus)

	for _, job := range d.cfg.DaemonConfiguration.Jobs {
		d.AddJob(statusChan, job)
	}

	// Handle status updates and errors
	go func() {
		for msg := range statusChan {
			if msg.Err != nil {
				fmt.Errorf("error in job=%s\t%w\n", msg.JobId, msg.Err)
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
		close(statusChan)

		fmt.Println("Shutdown complete")
		os.Exit(0)
	}()

	// Start cron jobs
	d.chron.Start()
	fmt.Println("Next job execution: ", d.chron.Entries()[0].Next)

	// Block here until interrupted
	select {}

}

func (d Daemon) orderFood(ch chan jobStatus, Strategy string, WeeksInAdvance int) func() {
	return func() {
		var jobId = "orderFood"
		menu, err := d.cli.GetMenuWeeks(WeeksInAdvance)
		if err != nil {
			ch <- jobStatus{JobId: jobId, Err: fmt.Errorf("trouble getting menu: %v", err.Error())}
		}

		dishes, err := ChooseDishesByStrategy(Strategy, menu)
		if err != nil {
			ch <- jobStatus{JobId: jobId, Err: fmt.Errorf("Trouble choosing dish: %v", err.Error())}
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
		ch <- jobStatus{JobId: jobId, Msg: fmt.Sprintf("Food ordered with %v strategy, for %v weeks", Strategy, WeeksInAdvance)}
	}
}

func (d Daemon) sendLowBalanceNotification(ch chan jobStatus, MinBalance int, Topic string, Template string) func() {
	return func() {
		var jobId = "sendLowBalanceNotification"
		ch <- jobStatus{JobId: jobId, Msg: "Checking account balance"}

		user, err := d.cli.GetUser()
		if err != nil {
			ch <- jobStatus{JobId: jobId, Err: fmt.Errorf("trouble getting user details: %v", err.Error())}
		}

		if user.User.Customer.AccountBalance.Amount <= MinBalance {
			ch <- jobStatus{JobId: jobId, Msg: "Account balance below minimum"}
			SendTemplateNotification("thinkjd_munch_o_matic", "Low balance notification", Template, user)
		}

	}
}

func (d Daemon) updateMetrics(ch chan jobStatus) func() {
	return func() {
		var jobId = "updateMetrics"
		ch <- jobStatus{JobId: jobId, Msg: "Update metrics"}
		//upcommingDishes, err := d.cli.GetMenuWeeks(4)
		//if err != nil {

		//}
		//UpdateOrdersPlaced()
	}
}
