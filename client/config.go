package client

type LoginCredentials struct {
	User     string
	Password string
}

type Config struct {
	LoginCredentials LoginCredentials
}
