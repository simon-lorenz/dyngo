package config

import (
	"dyngo/logger"
	"errors"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type DyngoConfiguration struct {
	Cron      string                 `yaml:"cron" validate:"required"`
	Services  ServicesConfiguration  `yaml:"services" validate:"required"`
	Detection DetectionConfiguration `yaml:"detection" validate:"required"`
	Log       LogConfiguration       `yaml:"log" validate:"required"`
}

type ServicesConfiguration struct {
	Desec   ServiceConfiguration `yaml:"desec" validate:"required_without_all"`
	Porkbun ServiceConfiguration `yaml:"porkbun" validate:"required_without_all"`
}

type ServiceConfiguration struct {
	Username string                `yaml:"username" validate:"required"`
	Password string                `yaml:"password" validate:"required"`
	Domains  []DomainConfiguration `yaml:"domains" validate:"required,dive"`
}

type DomainConfiguration struct {
	Domain string `yaml:"domain" validate:"required,hostname"`
	V4     bool   `yaml:"v4"`
	V6     bool   `yaml:"v6"`
}

type DetectionConfiguration struct {
	V4 AddressDetectionConfiguration `yaml:"v4"`
	V6 AddressDetectionConfiguration `yaml:"v6"`
}

type AddressDetectionConfiguration struct {
	Web string `yaml:"web" validate:"omitempty,url"`
	Cmd string `yaml:"cmd"`
}

type LogConfiguration struct {
	Level string `yaml:"level" validate:"oneof=trace debug info warning error fatal"`
}

var Cron string
var Services ServicesConfiguration
var Detection DetectionConfiguration
var Log LogConfiguration

func getConfigurationFileAsBytes(path string) []byte {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		logger.Fatal.Println("Configuration file " + path + " missing!")
		os.Exit(1)
	}

	yamlFile, err := os.Open(path)

	if err != nil {
		logger.Fatal.Println("Error when reading " + path + ": " + err.Error())
		os.Exit(1)
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

		logger.Fatal.Println("Configuration file is invalid")
		os.Exit(1)
	}

	Cron = config.Cron
	Services = config.Services
	Detection = config.Detection
	Log = config.Log
}
