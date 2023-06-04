package config

import (
	"dyngo/helpers/dns"
	"dyngo/logger"
	"io/ioutil"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type DyngoConfiguration struct {
	Services  ServicesConfiguration  `yaml:"services" validate:"required"`
	Detection DetectionConfiguration `yaml:"detection" validate:"required"`
	Log       LogConfiguration       `yaml:"log" validate:"required"`
}

type ServicesConfiguration struct {
	Generic []*GenericServiceConfiguration `yaml:"generic" validate:"dive"`

	Desec   *ServiceConfiguration `yaml:"desec"`
	Porkbun *ServiceConfiguration `yaml:"porkbun"`
}

type GenericServiceConfiguration struct {
	Name                 string           `yaml:"name" validate:"required"`
	Protocol             string           `yaml:"protocol" validate:"required,oneof=dyndns2"`
	URL                  string           `yaml:"url" validate:"required,url"`
	ServiceConfiguration `yaml:",inline"` // https://github.com/go-yaml/yaml/issues/63
}

type ServiceConfiguration struct {
	Username string                `yaml:"username" validate:"required"`
	Password string                `yaml:"password" validate:"required"`
	Domains  []DomainConfiguration `yaml:"domains" validate:"required,dive"`
}

type DomainConfiguration struct {
	Name    string       `yaml:"name" validate:"required,hostname"`
	Records []dns.Record `yaml:"records" validate:"required,dive,oneof=A AAAA"`
}

type DetectionConfiguration struct {
	Triggers *struct {
		Cron    string `yaml:"cron"`
		Startup bool   `yaml:"startup"`
	} `yaml:"triggers" validate:"required"`
	Strategies *struct {
		V4 DetectionStrategies `yaml:"v4"`
		V6 DetectionStrategies `yaml:"v6"`
	} `yaml:"strategies" validate:"required"`
}

type DetectionStrategies struct {
	Web string `yaml:"web" validate:"omitempty,url"`
	Cmd string `yaml:"cmd"`
}

type LogConfiguration struct {
	Level string `yaml:"level" validate:"oneof=trace debug info warning error fatal"`
}

var Services ServicesConfiguration
var Detection DetectionConfiguration
var Log LogConfiguration

func getConfigurationFileAsBytes(path string) []byte {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		logger.Fatal.Println("Error when reading " + path + ": " + err.Error())
		os.Exit(1)
	}

	return file
}

func Parse(path string) {
	var config DyngoConfiguration

	yaml.Unmarshal(getConfigurationFileAsBytes(path), &config)

	v := validator.New()
	err := v.Struct(config)

	if err != nil {
		errors := err.(validator.ValidationErrors)

		for _, e := range errors {
			logger.Warn.Println("Validation error in configuration: " + e.Error())
		}

		logger.Fatal.Println("Configuration file is invalid")
		os.Exit(1)
	}

	Services = config.Services
	Detection = config.Detection
	Log = config.Log
}
