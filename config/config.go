package config

import (
	"dyngo/logger"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/tidwall/jsonc"
)

var IPv4CheckUrl string
var IPv4Enabled bool
var IPv6CheckUrl string
var IPv6Enabled bool
var Username string
var Token string
var Domains []string
var Cron string

func getConfigurationFileAsBytes() []byte {
	var pathToConfiguration = "/etc/dyngo/config.json"

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

func init() {
	type ConfigurationJson struct {
		IPv4CheckUrl string
		IPv6CheckUrl string
		Username     string
		Token        string
		Domains      []string
		Cron         string
	}

	var config ConfigurationJson

	json.Unmarshal(jsonc.ToJSON(getConfigurationFileAsBytes()), &config)

	// TODO: Validation

	IPv4CheckUrl = config.IPv4CheckUrl
	IPv4Enabled = config.IPv4CheckUrl != ""
	IPv6CheckUrl = config.IPv6CheckUrl
	IPv6Enabled = config.IPv6CheckUrl != ""
	Username = config.Username
	Token = config.Token
	Domains = config.Domains
	Cron = config.Cron
}
