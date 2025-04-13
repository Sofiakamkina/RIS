package app

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	MainServerPort  int `env:"MAIN_SERVER_PORT" envDefault:"8080"`
	ProbeServerPort int `env:"PROBE_SERVER_PORT" envDefault:"8081"`

	ManagerURL string `env:"MANAGER_URL" envDefault:"manager:8080"`
}

var (
	configInstance *Config
	configReadOnce sync.Once
	configError    error
)

func NewConfig() (*Config, error) {
	configReadOnce.Do(func() {
		configInstance = &Config{}
		configError = cleanenv.ReadEnv(configInstance)
	})

	return configInstance, configError
}
