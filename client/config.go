package client

// Config Global configuration struct
type Config struct {
	LoginCredentials   LoginCredentials
	SessionCredentials SessionCredentials
	Daemon             Daemon
}

// Client related
type LoginCredentials struct {
	User     string
	Password string
}

type SessionCredentials struct {
	SessionID  string
	UserId     int
	CustomerId int
}

// Daemon related
type Job struct {
	Name     string
	Schedule string
}

type Daemon struct {
	Jobs []Job
}
