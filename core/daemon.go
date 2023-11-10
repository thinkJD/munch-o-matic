package core

import (
	"fmt"
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
			sendLowBalanceEmail(StatusChan, minBalance, email)
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
			orderFood(StatusChan, strategy, weeks)
		})
		if err != nil {
			return fmt.Errorf("error adding job: %w", err)
		}
	default:
		return fmt.Errorf("%v is not a valid type", Job.Type)
	}
	return nil
}

func (d *Daemon) Run() error {

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

func orderFood(ch chan string, Strategy string, WeeksInAdvance int) {
	/*
		c, err := client.NewClient(ClientConfig)
		if err != nil {
			ch <- fmt.Sprintf("auto-order error: %w", err)
		}

		menu, err := c.GetMenuWeeks(WeeksInAdvance)
		if err != nil {
			ch <- fmt.Sprintf("auto-order error: %w", err)
		}

		dishes, err := client.ChooseDishesByStrategy(Strategy, menu)
		if err != nil {
			ch <- fmt.Sprintf("auto-order error: %w", err)
		}
		fmt.Println(dishes)
			for _, dish := range dishes {
				err := c.OrderDish(dish.OrderId, false)
				if err != nil {
					ch <- fmt.Sprintf("Could not order")
				}
			}
	*/
	ch <- fmt.Sprintf("Food ordered with %v strategy, for %v weeks", Strategy, WeeksInAdvance)
}

func sendLowBalanceEmail(ch chan string, MinBalance int, Email string) {
	// Simulating checking balance and sending email
	balance := 50 // let's say
	if balance < 100 {
		ch <- fmt.Sprintf("Balance < %v; Email sent to %v", MinBalance, Email)
	} else {
		ch <- "Balance is okay."
	}
}
