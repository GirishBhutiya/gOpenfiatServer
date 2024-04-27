package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/datastax-ext/astra-go-sdk"
	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

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
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	log.Println("LoadConfig")
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		log.Println("No config file found", err)
		log.Println(err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
func GetAstraDBSession(config Config) (*gocql.Session, error) {
	log.Println("GetAstraDBSession")
	cluster, err := gocqlastra.NewClusterFromURL("https://api.astra.datastax.com", config.AstraDBId, config.AstraDBApplicationToken, 10*time.Second)

	if err != nil {
		return nil, err
	}

	cluster.Timeout = 30 * time.Second

	return gocql.NewSession(*cluster)

}
func GetAstraDBClient(config Config) (*astra.Client, error) {
	log.Println("GetAstraDBClient")
	c, err := astra.NewStaticTokenClient(
		config.AstraDBApplicationToken,
		astra.WithSecureConnectBundle(config.BundlePath),
		astra.WithDefaultKeyspace(config.DatabaseKeySpace),
	)
	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	}

	return c, nil
}
