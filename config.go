package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	// self host details
	host string
	port int

	// k8s configs
	namespace string
}

func (c *Config) Init() {
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 9090)

	viper.SetDefault("namespace", "default")
}

func (c *Config) Load() {
	c.host = viper.GetString("host")
	c.port = viper.GetInt("port")
	c.namespace = viper.GetString("namespace")
}

func (c *Config) AddConfigs() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/atlas/")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("Fatal error config file %s \n", err)
	}
}

func GetConfig() *Config {
	c := Config{}
	c.Init()
	c.AddConfigs()
	c.Load()
	return &c
}
