package app

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	MainServerPort  int `env:"MAIN_SERVER_PORT" envDefault:"8080"`
	ProbeServerPort int `env:"PROBE_SERVER_PORT" envDefault:"8081"`

	Alphabet string `env:"ALPHABET" envDefault:"abcdefghijklmnopqrstuvwxyz0123456789"`

	TTL int `env:"TTL" envDefault:"10"`

	WorkerCount int    `env:"WORKER_COUNT" envDefault:"3"`
	WorkerURLs  string `env:"WORKER_URLS" envDefault:"worker1:8080,worker2:8080,worker3:8080"`
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
