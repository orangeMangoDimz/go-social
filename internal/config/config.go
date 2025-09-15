package config

import (
	"time"

	"github.com/orangeMangoDimz/go-social/internal/ratelimiter"
)

type Config struct {
	Addr        string
	Db          DbConfig
	Env         string
	ApiURL      string
	Mail        MailConfig
	FrontendURL string
	Auth        AuthConfig
	RedisCfg    RedisConfig
	RateLimiter ratelimiter.Config
}

type RedisConfig struct {
	Addr     string
	Password string
	Db       int
	Enabled  bool
}

type AuthConfig struct {
	Basic BasicConfig
	Token TokenConfig
}

type BasicConfig struct {
	User string
	Pass string
}

type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}

type MailConfig struct {
	SendGrid  SendGridConfig
	MailTrap  MailTrapConfig
	FromEmail string
	Exp       time.Duration
}

type SendGridConfig struct {
	ApiKey string
}

type MailTrapConfig struct {
	ApiKey string
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
