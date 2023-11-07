package types

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
