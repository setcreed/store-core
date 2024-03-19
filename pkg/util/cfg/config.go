package cfg

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

const (
	yamlConfig = "yaml"
)

type Config struct {
	configFile string
	configType string

	data []byte
}

func New() *Config {
	return &Config{}
}

func (c *Config) SetConfigFile(configFile string) {
	c.configFile = configFile
}

func (c *Config) SetConfigType(in string) {
	c.configType = in
}

func (c *Config) readInConfig() error {
	var err error
	c.data, err = ioutil.ReadFile(c.configFile)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Binding(out interface{}) error {
	if err := c.readInConfig(); err != nil {
		return err
	}
	switch c.configType {
	case yamlConfig:
		if err := yaml.Unmarshal(c.data, out); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported config type %s", c.configType)
	}

	return nil
}
