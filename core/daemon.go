package core

import (
	"fmt"
	"math/rand"
	"munch-o-matic/client"
	"time"

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

func (d *Daemon) AddJob(StatusChan chan string, Job Job) error {
	switch Job.Type {
	case "CheckBalance":
		email, ok1 := Job.Params["email"].(string)
		minBalance, ok2 := Job.Params["minbalance"].(int)
		if !ok1 || !ok2 {
			return fmt.Errorf("invalid parameter types for CheckBalance")
		}

		_, err := d.chron.AddFunc(Job.Schedule, func() {
			d.sendLowBalanceEmail(StatusChan, minBalance, email)
		})
		if err != nil {
			return fmt.Errorf("error adding job: %w", err)
		}

	case "Order":
		strategy, ok1 := Job.Params["strategy"].(string)
		weeks, ok2 := Job.Params["weeks"].(int)
		if !ok1 || !ok2 {
			return fmt.Errorf("invalid parameter types for Order")
		}
		_, err := d.chron.AddFunc(Job.Schedule, func() {
			d.orderFood(StatusChan, strategy, weeks)
		})
		if err != nil {
			return fmt.Errorf("error adding job: %w", err)
		}
	default:
		return fmt.Errorf("%v is not a valid type", Job.Type)
	}
	return nil
}

func (d Daemon) Run() error {

	statusChan := make(chan string)

	for _, job := range d.cfg.DaemonConfiguration.Jobs {
		d.AddJob(statusChan, job)
	}

	go func() {
		for msg := range statusChan {
			fmt.Println(msg)
		}
	}()

	d.chron.Start()
	fmt.Println("Next job execution: ", d.chron.Entries()[0].Next)
	// Let it run for 2 minutes to see a couple of executions
	time.Sleep(2 * time.Minute)
	d.chron.Stop()

	return nil
}

func (d Daemon) orderFood(ch chan string, Strategy string, WeeksInAdvance int) func() {
	return func() {
		menu, err := d.cli.GetMenuWeeks(WeeksInAdvance)
		if err != nil {
			ch <- fmt.Sprintf("auto-order error: %w", err)
		}

		dishes, err := ChooseDishesByStrategy(Strategy, menu)
		if err != nil {
			ch <- fmt.Sprintf("auto-order error: %w", err)
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
		ch <- fmt.Sprintf("Food ordered with %v strategy, for %v weeks", Strategy, WeeksInAdvance)
	}
}

func (d Daemon) sendLowBalanceEmail(ch chan string, MinBalance int, Email string) func() {
	return func() {
		// Simulating checking balance and sending email
		balance := 50 // let's say
		if balance < 100 {
			ch <- fmt.Sprintf("Balance < %v; Email sent to %v", MinBalance, Email)
		} else {
			ch <- "Balance is okay."
		}
	}
}

// Pick dishes automatically based on a few strategies
func ChooseDishesByStrategy(Strategy string, UpcomingDishes map[string][]client.UpcomingDish) (map[int]client.UpcomingDish, error) {
	retVal := map[int]client.UpcomingDish{}

	// Helper function to decide if menu should be skipped
	shouldSkipMenu := func(menu []client.UpcomingDish) bool {
		for _, dish := range menu {
			if dish.Booked || dish.Dummy {
				return true
			}
		}
		return false
	}

	// Iterate over the dishes of the day
	for _, menu := range UpcomingDishes {
		if shouldSkipMenu(menu) {
			continue
		}
		// Choose dish based on the strategy
		switch Strategy {

		case "SchoolFav":
			var maxPos, maxVal int
			for i, dish := range menu {
				if dish.Orders > maxVal {
					maxPos = i
					maxVal = dish.Orders
				}
			}
			retVal[menu[maxPos].OrderId] = menu[maxPos]

		case "Random":
			randomInt := rand.Intn(len(menu))
			retVal[menu[randomInt].OrderId] = menu[randomInt]

		case "PersonalFav":
			/* TODO: Add personal order count in getMenuWeek or structure the code better.
			var maxPos, maxVal int
			for i, dish := range menu {
				GetOrderCount()
			}
			*/
			return map[int]client.UpcomingDish{}, fmt.Errorf("PersonalFav is not implemented, sorry")

		default:
			return map[int]client.UpcomingDish{}, fmt.Errorf("%v is not a valid strategy", Strategy)
		}
	}
	return retVal, nil
}
