package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/datastax-ext/astra-go-sdk"
	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"
)

var BundlePath = "secure-connect-openfiat-test.zip"

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variable
type Config struct {
	ListenPort              string        `mapstructure:"LISTEN_PORT"`
	AccessTokenDuration     time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	AstraDBApplicationToken string        `mapstructure:"ASTRA_DB_APPLICATION_TOKEN"`
	AstraDBAPIEndpoint      string        `mapstructure:"ASTRA_DB_API_ENDPOINT"`
	AstraDBId               string        `mapstructure:"ASTRA_DB_ID"`
	TokenSymmetricKey       string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	RefreshTokenDuration    time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	BundlePath              string        `mapstructure:"BUNDLE_PATH"`
	DatabaseKeySpace        string        `mapstructure:"DATABASE_KEYSPACE"`
	DBMigrateUp             bool          `mapstructure:"DB_MIGRATE_UP"`
	OTPAPIURL               string        `mapstructure:"OTP_API_URL"`
	Production              bool          `mapstructure:"PRODUCTION"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {

	accessTokenDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		log.Println("can not parse ACCESS_TOKEN_DURATION", err)
		return
	}
	refreshTokenDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		log.Println("can not parse REFRESH_TOKEN_DURATION", err)
		return
	}
	dbMigrateUp, err := strconv.ParseBool(os.Getenv("DB_MIGRATE_UP"))
	if err != nil {
		log.Println("can not parse DB_MIGRATE_UP", err)
		return
	}
	production, err := strconv.ParseBool(os.Getenv("PRODUCTION"))
	if err != nil {
		log.Println("can not parse PRODUCTION", err)
		return
	}
	conf := Config{
		ListenPort:              os.Getenv("LISTEN_PORT"),
		AccessTokenDuration:     accessTokenDuration,
		AstraDBApplicationToken: os.Getenv("ASTRA_DB_APPLICATION_TOKEN"),
		AstraDBAPIEndpoint:      os.Getenv("ASTRA_DB_API_ENDPOINT"),
		AstraDBId:               os.Getenv("ASTRA_DB_ID"),
		TokenSymmetricKey:       os.Getenv("TOKEN_SYMMETRIC_KEY"),
		RefreshTokenDuration:    refreshTokenDuration,
		BundlePath:              os.Getenv("BUNDLE_PATH"),
		DatabaseKeySpace:        os.Getenv("DATABASE_KEYSPACE"),
		DBMigrateUp:             dbMigrateUp,
		OTPAPIURL:               os.Getenv("OTP_API_URL"),
		Production:              production,
	}

	return conf, nil
}
func GetAstraDBSession(config Config) (*gocql.Session, error) {
	cluster, err := gocqlastra.NewClusterFromURL("https://api.astra.datastax.com", config.AstraDBId, config.AstraDBApplicationToken, 10*time.Second)

	if err != nil {
		return nil, err
	}

	cluster.Timeout = 30 * time.Second

	return gocql.NewSession(*cluster)

}
func GetAstraDBClient(config Config) (*astra.Client, error) {
	c, err := astra.NewStaticTokenClient(
		config.AstraDBApplicationToken,
		astra.WithSecureConnectBundle(BundlePath),
		astra.WithDefaultKeyspace(config.DatabaseKeySpace),
	)
	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	}

	return c, nil
}
