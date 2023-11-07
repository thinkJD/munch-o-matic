package core

import (
	"errors"
	"fmt"
	"munch-o-matic/client"
	. "munch-o-matic/client/types"

	"github.com/robfig/cron/v3"
)

func New(Cfg DaemonConfiguration, Cli *client.RestClient) (*Daemon, error) {
	retVal := Daemon{}
	retVal.cli = Cli
	retVal.chron = cron.New(cron.WithSeconds())

	return &Daemon{}, nil
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

}

/*
func (d *Daemon) Run(config DaemonConfiguration) error {

		statusChan := make(chan string)

		for _, job := range config.Daemon.Jobs {


		}

		go func() {
			for msg := range statusChan {
				fmt.Println(msg)
			}
		}()

		c.Start()
		fmt.Println("Next job execution: ", c.Entries()[0].Next)
		// Let it run for 2 minutes to see a couple of executions
		time.Sleep(2 * time.Minute)
		c.Stop()

		return nil
	}
*/
func orderFood(ch chan string, Strategy string, WeeksInAdvance int, ClientConfig Config) {
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
	/*
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

func ValidateConfig(config Daemon) error {
	for _, job := range config.Jobs {
		if len(job.Type) == 0 {
			return errors.New("a job is missing a Type")
		}
		if len(job.Name) == 0 {
			return errors.New("a job is missing a Name")
		}
		if job.Schedule == "" {
			return errors.New("a job is missing a Schedule")
		}

		switch job.Type {
		case "CheckBalance":
			if email, ok := job.Params["email"].(string); !ok || email == "" {
				return fmt.Errorf("CheckBalance job '%v' is missing or has an invalid Email", job.Name)
			}
			if minBalance, ok := job.Params["minbalance"].(int); !ok || minBalance <= 0 {
				return fmt.Errorf("CheckBalance job '%v' is missing or has an invalid MinBalance", job.Name)
			}
		case "Order":
			if strategy, ok := job.Params["strategy"].(string); !ok || strategy == "" {
				return fmt.Errorf("Order job '%v' is missing or has an invalid Strategy", job.Name)
			}
			if weeks, ok := job.Params["weeks"].(int); !ok || weeks <= 0 {
				return fmt.Errorf("Order job '%v' is missing or has an invalid WeeksInAdvance", job.Name)
			}
		default:
			return fmt.Errorf("unknown job type '%v'", job.Type)
		}
	}
	return nil
}
