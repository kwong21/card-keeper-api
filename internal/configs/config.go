package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Configuration holds the the config file values
type Configuration interface {
	DBConfigs() DBConfiguration
	APIAllowedOrigin() string
	APIListenPort() string
	AuthConfiguration() AuthConfiguration
}

// DBConfiguration holds the database configuration values
type DBConfiguration struct {
	Type       string
	Host       string
	User       string
	Password   string
	Database   string
	ReplicaSet string
}

// AuthConfiguration holds the authentication configuration values
type AuthConfiguration struct {
	Audience string
	Issuer   string
	JWKS     string
}

type configuration struct {
	DB  DBConfiguration
	API struct {
		AllowedOrigin string
		ListenPort    string
	}
	Auth AuthConfiguration
}

// Default configuration for the API server
func Default() Configuration {
	c := new(configuration)

	c.DB.Type = "in-memory"
	c.API.AllowedOrigin = "http://localhost:4200"
	c.API.ListenPort = "8080"

	return c
}

// NewFromFile reads JSON file to get config values. Returns DEFAULT if Can't parse file
func NewFromFile(path string) Configuration {
	jsonFile, err := os.Open(path)

	if err != nil {
		return Default()
	}
	defer jsonFile.Close()

	c := new(configuration)
	b, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(b, &c)

	return c
}

func (c *configuration) DBConfigs() DBConfiguration {
	return c.DB
}

func (c *configuration) APIAllowedOrigin() string {
	return c.API.AllowedOrigin
}

func (c *configuration) APIListenPort() string {
	return c.API.ListenPort
}

func (c *configuration) AuthConfiguration() AuthConfiguration {
	return c.Auth
}
