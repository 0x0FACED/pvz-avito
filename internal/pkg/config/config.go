package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Database DatabaseConfig
	Server   ServerConfig
	Logger   LoggerConfig
}

type DatabaseConfig struct {
	DSN               string        `env:"DATABASE_DSN,required"`
	MaxOpenConns      int32         `env:"DATABASE_MAX_OPEN_CONNS"`
	MaxIdleConns      int32         `env:"DATABASE_MAX_IDLE_CONNS"`
	ConnMaxLifetime   time.Duration `env:"DATABASE_CONN_MAX_LIFETIME"`
	ConnMaxIdleTime   time.Duration `env:"DATABASE_CONN_MAX_IDLE_LIFETIME"`
	ConnectionTimeout time.Duration `env:"DATABASE_CONNECTION_TIMEOUT"`
	HealthCheckPeriod time.Duration `env:"DATABASE_HEALTH_CHECK_PERIOD"`
}

type ServerConfig struct {
	// Server options
	Host string `env:"SERVER_HOST" envDefault:"localhost"`
	Port string `env:"SERVER_PORT" envDefault:"8080"`

	// http
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT"`

	DebugMode bool `env:"SERVER_DEBUG_MODE"`

	JWTSecret string `env:"SERVER_JWT_SECRET" envDefault:"test-secret-key"`
}

type LoggerConfig struct {
	// Main options
	LogLevel string `env:"LOGGER_LEVEL" envDefault:"debug"`
	NoColor  bool   `env:"LOGGER_NO_COLOR" envDefault:"false"`

	// Time options
	TimeFormat   string `env:"LOGGER_TIME_FORMAT" envDefault:"2006-01-02T15:04:05Z"`
	TimeLocation string `env:"LOGGER_TIME_LOCATION" envDefault:"UTC"`

	// Output options
	PartsOrder    string `env:"LOGGER_PARTS_ORDER" envDefault:"time,level,logger,message"`
	PartsExclude  string `env:"LOGGER_PARTS_EXCLUDE" envDefault:""`
	FieldsOrder   string `env:"LOGGER_FIELDS_ORDER" envDefault:""`
	FieldsExclude string `env:"LOGGER_FIELDS_EXCLUDE" envDefault:""`

	// Dir for all log files
	LogsDir string `env:"LOGGER_LOGS_DIR" envDefault:"./logs"`
}

// MustLoad loads config from .env file and parse it to CodexConig.
// Panics if err != nil
func MustLoad() *AppConfig {
	if err := godotenv.Load(); err != nil {
		panic("failed to load config, err: " + err.Error())
	}

	cfg := &AppConfig{}

	if err := env.Parse(&cfg.Database); err != nil {
		panic("failed to parse database config, err: " + err.Error())
	}

	if err := env.Parse(&cfg.Server); err != nil {
		panic("failed to parse server config, err: " + err.Error())
	}

	if err := env.Parse(&cfg.Logger); err != nil {
		panic("failed to parse logger config, err: " + err.Error())
	}

	return cfg
}

// cfg for integration tests
func LoadTest() *AppConfig {
	return &AppConfig{
		Database: DatabaseConfig{
			DSN:               "postgres://test_user:test_pass@localhost:5432/pvz_avito_test_db?sslmode=disable",
			MaxOpenConns:      10,
			MaxIdleConns:      5,
			ConnMaxLifetime:   30 * time.Minute,
			ConnMaxIdleTime:   5 * time.Minute,
			ConnectionTimeout: 10 * time.Second,
			HealthCheckPeriod: 15 * time.Second,
		},
		Server: ServerConfig{
			Host:         "127.0.0.1",
			Port:         "8080",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  30 * time.Second,
			DebugMode:    true,
			JWTSecret:    "test-jwt-secret",
		},
		Logger: LoggerConfig{
			LogLevel:      "debug",
			NoColor:       false,
			TimeFormat:    "2006-01-02T15:04:05Z",
			TimeLocation:  "UTC",
			PartsOrder:    "time,level,logger,message",
			PartsExclude:  "",
			FieldsOrder:   "",
			FieldsExclude: "",
			LogsDir:       "./test_logs",
		},
	}
}
