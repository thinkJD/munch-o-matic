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

type Notification struct {
	Error struct {
		Enabled bool
		Topic   string
	}
	Status struct {
		Enabled bool
		Topic   string
	}
}

type DaemonConfiguration struct {
	Jobs         []Job
	Notification Notification
}

type Config struct {
	DaemonConfiguration DaemonConfiguration
}

// Config Validation
// #################
func ValidateConfig(config Config) error {
	// Validate Notifications
	if err := validateNotifications(config.DaemonConfiguration.Notification); err != nil {
		return err
	}

	// Validate Jobs
	for _, job := range config.DaemonConfiguration.Jobs {
		if err := validateJob(job); err != nil {
			return err
		}
	}
	return nil
}

func validateNotifications(notification Notification) error {
	if notification.Error.Enabled {
		if notification.Error.Topic == "" {
			return errors.New("error notifications must have a topic")
		}
	}

	if notification.Status.Enabled {
		if notification.Status.Topic == "" {
			return errors.New("status notifications must have a topic")
		}
	}

	return nil
}

func validateJob(job Job) error {
	if len(job.Name) == 0 {
		return errors.New("a job is missing a Name")
	}
	if job.Schedule == "" {
		return errors.New("a job is missing a Schedule")
	}
	if len(job.Type) == 0 {
		return errors.New("a job is missing a Type")
	}

	switch job.Type {
	case "CheckBalance":
		// Check for specific parameters like 'minbalance' and 'template'
		if minBalance, ok := job.Params["minbalance"].(int); !ok || minBalance <= 0 {
			return fmt.Errorf("CheckBalance job '%v' is missing or has an invalid MinBalance", job.Name)
		}
	case "UpdateMetrics":
		// Specific validation for UpdateMetrics job if needed
		break
	// Add cases for other job types as necessary
	default:
		return fmt.Errorf("unknown job type '%v'", job.Type)
	}
	return nil
}
