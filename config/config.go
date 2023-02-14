package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var instantiated *Config
var once sync.Once

type Config struct {
	CollectionName   string `env:"COLLECTION_NAME" envDefault:""`
	DatabaseName     string `env:"DATABASE_NAME" envDefault:""`
	ConnectionString string `env:"CONNECTION_STR" envDefault:""`
}

// Read and parse the configuration file
func read() *Config {
	config := Config{}
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}
	return &config
}

func GetInstance() *Config {
	once.Do(func() {
		instantiated = read()
	})
	return instantiated
}
