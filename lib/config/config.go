package config

import (
	// stdlib
	"io/ioutil"
	"path/filepath"
	// 3rd-party
	"gopkg.in/yaml.v2"
	"lab.pztrn.name/golibs/mogrus"
)

// HTTPServerConfiguration handles HTTP server configuration in config file
type HTTPServerConfiguration struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// HTTPClientConfiguration handles HTTP client configuration in config file
type HTTPClientConfiguration struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Signature string `yaml:"signature"`
}

// Config is a struct which represents config file structure
type Config struct {
	HTTPListener HTTPServerConfiguration `yaml:"receiver"`
	HTTPSender   HTTPClientConfiguration `yaml:"sender"`
}

// Init is a configuration initializer
func (c *Config) Init(log *mogrus.LoggerHandler, configPath string) {
	log.Info("Config file path: " + configPath)
	fname, _ := filepath.Abs(configPath)
	yamlFile, yerr := ioutil.ReadFile(fname)
	if yerr != nil {
		log.Fatal("Can't read config file")
	}

	yperr := yaml.Unmarshal(yamlFile, c)
	if yperr != nil {
		log.Fatal("Can't parse config file")
	}
}
