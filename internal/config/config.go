package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address"`
}

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	Storage    string `yaml:"storage" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

// function to load configuration file data. It must execute successfully else the process must terminate
func MustLoad() *Config {
	// get config file path from env variable
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		// check if config path is passed through command line argument
		flags := flag.String("config", "", "configuration file path")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatalf("config path is not set")
		}
	}

	// check if file exists at configPath
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("File doesn't exist at path: %s\n", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Unable to read config file: %s\n", configPath)
	}

	fmt.Printf("config details: %+v\n", cfg)

	return &cfg

}
