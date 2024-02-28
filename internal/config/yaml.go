package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func LoadYamlConfig() (*Config, error) {
	configFilePath := os.Getenv(ConfigFilePath)

	cfg := &Config{}

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
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
