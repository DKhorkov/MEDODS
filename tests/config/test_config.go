package testconfig

import (
	"fmt"
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	"github.com/DKhorkov/medods/internal/config"
)

type TestDatabaseConfig struct {
	Driver        string
	DSN           string
	MigrationsDir string
}

type TestRefreshTokenConfig struct {
	GUID  string
	Value string
}

type TestConfig struct {
	Database     TestDatabaseConfig
	RefreshToken TestRefreshTokenConfig
	SMTP         config.SMTPConfig
	JWT          config.JWTConfig
	Logging      config.LoggingConfig
	IP           string
	HashCost     int
}

func New() *TestConfig {
	return &TestConfig{
		Database: TestDatabaseConfig{
			Driver:        "sqlite3",
			DSN:           "file::memory:?cache=shared", // "test.db" can be also used
			MigrationsDir: "/internal/database/migrations",
		},
		RefreshToken: TestRefreshTokenConfig{
			GUID:  "42385b4f-d5cd-4543-acef-229fb60fe35g",
			Value: "testValue",
		},
		SMTP: config.SMTPConfig{
			Host:     "smtp.freesmtpservers.com",
			Port:     25,
			Login:    "smtp",
			Password: "smtp",
		},
		JWT: config.JWTConfig{
			Algorithm:       "HS256",
			SecretKey:       "testSecret",
			RefreshTokenTTL: time.Minute * 5,
			AccessTokenTTL:  time.Minute * 1,
		},
		IP:       "127.0.0.1",
		HashCost: 4,
		Logging: config.LoggingConfig{
			Level:       logging.LogLevels.DEBUG,
			LogFilePath: fmt.Sprintf("logs/%s.log", time.Now().Format("02-01-2006")),
		},
	}
}
