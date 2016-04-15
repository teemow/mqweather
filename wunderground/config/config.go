package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type errRequiredProperties []string

func (err errRequiredProperties) Error() string {
	return fmt.Sprintf("Required: %s\n", strings.Join(err, ", "))
}

type Config struct {
	ApiKey  string
	ApiHost string
}

func FromBytes(data []byte) (*Config, error) {
	conf := &Config{}
	err := json.Unmarshal(data, conf)

	if err != nil {
		return conf, err
	}

	return conf, err

}

func Configuration(configFilePath string) *Config {

	if configFilePath == "" {
		panic("Missing configuration file path.")
	}
	configFileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("Error loading file: %s", err))
	}

	config, err := FromBytes(configFileBytes)
	if err != nil {
		panic(fmt.Sprintf("Error validating configuration: %s", err))
	}

	return config
}

func (c *Config) Validate() error {
	requiredProp := errRequiredProperties{}

	if len(c.ApiKey) == 0 {
		requiredProp = append(requiredProp, "API Key")
	}

	if len(c.ApiHost) == 0 {
		requiredProp = append(requiredProp, "API Host")
	}

	if len(requiredProp) > 0 {
		return requiredProp
	}

	return nil
}
