package core

import (
	"errors"
	"fmt"
)

type Job struct {
	Type     string
	Name     string
	Schedule string
	Params   map[string]interface{}
}

type DaemonConfiguration struct {
	Jobs []Job
}

type Config struct {
	DaemonConfiguration DaemonConfiguration
}

func ValidateConfig(Config Config) error {
	for _, job := range Config.DaemonConfiguration.Jobs {
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
			if topic, ok := job.Params["topic"].(string); !ok || topic == "" {
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
