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

	fmt.Println(d.chron.Entries())

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
			ch <- fmt.Sprintf("trouble getting menu: %v", err.Error())
		}

		dishes, err := ChooseDishesByStrategy(Strategy, menu)
		if err != nil {
			ch <- fmt.Sprintf("Trouble choosing dish: %v", err.Error())
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

func (d Daemon) sendLowBalanceNotification(ch chan string, MinBalance int, Topic string, Template string) func() {
	return func() {
		ch <- fmt.Sprint("Checking account balance")
		user, err := d.cli.GetUser()
		if err != nil {
			ch <- fmt.Sprintf("trouble getting user details: %v", err.Error())
		}

		if user.User.Customer.AccountBalance.Amount <= MinBalance {
			ch <- fmt.Sprint("Account balance below minimum")
			template := "Hello, your balance: {{.User.Customer.AccountBalance.Amount}}"
			SendTemplateNotification("thinkjd_munch_o_matic", "Low balance notification", template, user)
		}

	}
}
