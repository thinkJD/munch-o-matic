package client

import "errors"

// Config Global configuration struct
type Config struct {
	LoginCredentials   LoginCredentials
	SessionCredentials SessionCredentials
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

func ValidateConfig(Cfg Config) error {
	if len(Cfg.LoginCredentials.User) == 0 {
		return errors.New("logincredentials.user is required")
	}
	if len(Cfg.LoginCredentials.Password) == 0 {
		return errors.New("logincredentials.password is required")
	}
	return nil
}
