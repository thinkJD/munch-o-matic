package client

type LoginCredentials struct {
	User     string
	Password string
}

type SessionCredentials struct {
	SessionID  string
	UserId     int
	CustomerId int
}

type Config struct {
	LoginCredentials   LoginCredentials
	SessionCredentials SessionCredentials
}
