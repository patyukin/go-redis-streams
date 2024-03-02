package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func LoadYamlConfig() (*Config, error) {
	configFilePath := os.Getenv(ConfigFilePath)

	cfg := &Config{}

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w, path: %s", err, configFilePath)
	}

	defer func() {
		errFileClose := file.Close()
		if errFileClose != nil {
			log.Fatal(errFileClose)
		}
	}()

	d := yaml.NewDecoder(file)

	if err = d.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
