package utils

import (
	"errors"
	"fmt"
	. "munch-o-matic/client/types"
	"time"

	"github.com/robfig/cron/v3"
)

func Run(config Config) error {
	c := cron.New(cron.WithSeconds())
	statusChan := make(chan string)

	for _, job := range config.Daemon.Jobs {
		switch job.Type {

		case "CheckBalance":
			email, ok1 := job.Params["email"].(string)
			minBalance, ok2 := job.Params["minbalance"].(int)
			if !ok1 || !ok2 {
				return fmt.Errorf("invalid parameter types for CheckBalance")
			}

			_, err := c.AddFunc(job.Schedule, func() {
				sendLowBalanceEmail(statusChan, minBalance, email)
			})
			if err != nil {
				return fmt.Errorf("error adding job: %w", err)
			}

		case "Order":
			strategy, ok1 := job.Params["strategy"].(string)
			weeks, ok2 := job.Params["weeksinadvance"].(int)
			if !ok1 || !ok2 {
				return fmt.Errorf("invalid parameter types for Order")
			}
			_, err := c.AddFunc(job.Schedule, func() {
				orderFood(statusChan, strategy, weeks)
			})
			if err != nil {
				return fmt.Errorf("error adding job: %w", err)
			}
		default:
			return fmt.Errorf("%v is not a valid type", job.Type)
		}
	}

	go func() {
		for msg := range statusChan {
			fmt.Println(msg)
		}
	}()

	c.Start()
	// Let it run for 2 minutes to see a couple of executions
	time.Sleep(2 * time.Minute)
	c.Stop()

	return nil
}

func orderFood(ch chan string, Strategy string, WeeksInAdvance int) {
	// Simulating ordering food
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
			if email, ok := job.Params["Email"].(string); !ok || email == "" {
				return fmt.Errorf("CheckBalance job '%v' is missing or has an invalid Email", job.Name)
			}
			if minBalance, ok := job.Params["MinBalance"].(int); !ok || minBalance <= 0 {
				return fmt.Errorf("CheckBalance job '%v' is missing or has an invalid MinBalance", job.Name)
			}
		case "Order":
			if strategy, ok := job.Params["Strategy"].(string); !ok || strategy == "" {
				return fmt.Errorf("Order job '%v' is missing or has an invalid Strategy", job.Name)
			}
			if weeks, ok := job.Params["WeeksInAdvance"].(int); !ok || weeks <= 0 {
				return fmt.Errorf("Order job '%v' is missing or has an invalid WeeksInAdvance", job.Name)
			}
		default:
			return fmt.Errorf("unknown job type '%v'", job.Type)
		}
	}
	return nil
}
