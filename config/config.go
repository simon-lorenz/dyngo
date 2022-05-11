package config

import (
	"dyngo/logger"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/tidwall/jsonc"

	"github.com/go-playground/validator/v10"
)

type DomainConfiguration struct {
	Domain string `json:"domain" validate:"required,hostname"`
	V4     bool   `json:"v4"`
	V6     bool   `json:"v6"`
}

type ServiceConfiguration struct {
	Username string                `json:"username" validate:"required"`
	Password string                `json:"password" validate:"required"`
	Domains  []DomainConfiguration `json:"domains" validate:"required,dive"`
}

type ServicesConfiguration struct {
	Desec ServiceConfiguration `json:"desec" validate:"required_without_all"`
}

type IPv4AddressDetectionConfiguration struct {
	Web string `json:"web" validate:"required,url"`
}

type IPv6AddressDetectionConfiguration struct {
	Web string `json:"web" validate:"required,url"`
}

type DyngoConfiguration struct {
	Cron                 string                            `json:"cron" validate:"required"`
	Services             ServicesConfiguration             `json:"services" validate:"required,dive"`
	IPv4AddressDetection IPv4AddressDetectionConfiguration `json:"v4AddressDetection" validate:"required"`
	IPv6AddressDetection IPv6AddressDetectionConfiguration `json:"v6AddressDetection" validate:"required"`
	LogLevel             string                            `json:"logLevel"`
}

var Cron string
var Services ServicesConfiguration
var IPv4AddressDetection IPv4AddressDetectionConfiguration
var IPv6AddressDetection IPv6AddressDetectionConfiguration
var LogLevel int

func getConfigurationFileAsBytes() []byte {
	var pathToConfiguration = "/etc/dyngo/config.jsonc"

	if _, err := os.Stat(pathToConfiguration); errors.Is(err, os.ErrNotExist) {
		logger.Error.Println("Configuration file " + pathToConfiguration + " missing!")
		os.Exit(1)
	}

	jsonFile, err := os.Open(pathToConfiguration)

	if err != nil {
		logger.Error.Println("Error when reading " + pathToConfiguration + ": " + err.Error())
		os.Exit(1)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	return byteValue
}

func Parse() {
	var config DyngoConfiguration

	json.Unmarshal(jsonc.ToJSON(getConfigurationFileAsBytes()), &config)

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

	switch strings.ToLower(config.LogLevel) {
	case "trace":
		LogLevel = logger.LogLevelTrace
		break
	case "debug":
		LogLevel = logger.LogLevelDebug
		break
	case "info":
		LogLevel = logger.LogLevelInfo
		break
	case "warning":
		LogLevel = logger.LogLevelWarning
		break
	case "error":
		LogLevel = logger.LogLevelError
		break
	case "fatal":
		LogLevel = logger.LogLevelFatal
		break
	default:
		if config.LogLevel != "" {
			logger.Warn.Printf("[Configuration] Log Level '%v' is invalid", config.LogLevel)
		}

		logger.Info.Println("[Configuration] Falling back to log level 'info'")
	}
}
