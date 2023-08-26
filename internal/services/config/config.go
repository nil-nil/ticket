package config

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP struct {
		Port          int    `yaml:"port"`
		ListenAddress string `yaml:"listenAddress"`
	} `yaml:"httpServer"`
	Auth struct {
		JWT *struct {
			SigningMethod string `yaml:"signingMethod"`
			TokenLifetime uint64 `yaml:"tokenLifetime"`
			PublicKey     string `yaml:"publicKey"`
			PrivateKey    string `yaml:"privateKey"`
		} `yaml:"jwt"`
	} `yaml:"auth"`
}

func GetConfig(r io.Reader) (Config, error) {
	config := Config{}
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func ReadAndParseConfigFile(fileName string) (Config, error) {
	configFile, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, fmt.Errorf("unable to read config.yaml. %e", err)
	}
	configReader := bytes.NewBuffer(configFile)
	config, err := GetConfig(configReader)
	if err != nil {
		return Config{}, fmt.Errorf("unable to parse config.yaml. %e", err)
	}

	return config, nil
}
