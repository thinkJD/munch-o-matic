package core

import (
	"munch-o-matic/client"

	"github.com/robfig/cron/v3"
)

// Daemon related
type Daemon struct {
	d     *Daemon
	chron *cron.Cron
	cli   *client.RestClient
}

type Job struct {
	Type     string
	Name     string
	Schedule string
	Params   map[string]interface{}
}

type DaemonConfiguration struct {
	Jobs []Job
}
