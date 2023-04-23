package config

import (
	"dyngo/logger"
	"errors"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type DomainConfiguration struct {
	Domain string `yaml:"domain" validate:"required,hostname"`
	V4     bool   `yaml:"v4"`
	V6     bool   `yaml:"v6"`
}

type ServiceConfiguration struct {
	Username string                `yaml:"username" validate:"required"`
	Password string                `yaml:"password" validate:"required"`
	Domains  []DomainConfiguration `yaml:"domains" validate:"required,dive"`
}

type ServicesConfiguration struct {
	Desec ServiceConfiguration `yaml:"desec" validate:"required_without_all"`
}

type IPv4AddressDetectionConfiguration struct {
	Web string `yaml:"web" validate:"required,url"`
}

type IPv6AddressDetectionConfiguration struct {
	Web string `yaml:"web" validate:"required,url"`
}

type DyngoConfiguration struct {
	Cron                 string                            `yaml:"cron" validate:"required"`
	Services             ServicesConfiguration             `yaml:"services" validate:"required,dive"`
	IPv4AddressDetection IPv4AddressDetectionConfiguration `yaml:"v4AddressDetection" validate:"required"`
	IPv6AddressDetection IPv6AddressDetectionConfiguration `yaml:"v6AddressDetection" validate:"required"`
	LogLevel             string                            `yaml:"logLevel" validate:"oneof=trace debug info warning error fatal"`
}

var Cron string
var Services ServicesConfiguration
var IPv4AddressDetection IPv4AddressDetectionConfiguration
var IPv6AddressDetection IPv6AddressDetectionConfiguration
var LogLevel int

func getConfigurationFileAsBytes(path string) []byte {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		logger.Error.Println("Configuration file " + path + " missing!")
		os.Exit(1)
	}

	yamlFile, err := os.Open(path)

	if err != nil {
		logger.Error.Println("Error when reading " + path + ": " + err.Error())
		os.Exit(1)
	} else {
		logger.Info.Println("Using configuration file " + path)
	}

	// defer the closing of our yamlFile so that we can parse it later on
	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)

	return byteValue
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

		logger.Error.Println("Configuration file is invalid")
		os.Exit(1)
	}

	Cron = config.Cron
	Services = config.Services
	IPv4AddressDetection = config.IPv4AddressDetection
	IPv6AddressDetection = config.IPv6AddressDetection
	LogLevel, _ = logger.LogLevelByString(config.LogLevel)
}
