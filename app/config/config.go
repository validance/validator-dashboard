package config

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config *Config
var once sync.Once

type Config struct {
	App          App                       `mapstructure:"app"`
	CoingeckoIds []string                  `mapstructure:"coingeckoIds"`
	Chains       []string                  `mapstructure:"chains"`
	Database     Database                  `mapstructure:"database"`
	Cosmos       map[string]CosmosAppchain `mapstructure:"cosmos"`
	Aptos        Aptos                     `mapstructure:"aptos"`
	Polygon      Polygon                   `mapstructure:"polygon"`
}

type App struct {
	Host         string   `mapstructure:"host"`
	Port         string   `mapstructure:"port"`
	Cron         string   `mapstructure:"cron"`
	AllowOrigins []string `mapstructure:"allowOrigins"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
}

type CosmosAppchain struct {
	GrpcUrl               string   `mapstructure:"grpcUrl"`
	Denom                 string   `mapstructure:"denom"`
	ValidatorOperatorAddr string   `mapstructure:"validatorOperatorAddr"`
	ValidatorAddr         string   `mapstructure:"validatorAddr"`
	Exponent              int      `mapstructure:"exponent"`
	GrantAddrs            []string `mapstructure:"grantAddrs"`
}

type Aptos struct {
	validatorAddr string
}

type Polygon struct {
	ValidatorIndex int    `mapstructure:"validatorIndex"`
	SignerAddr     string `mapstructure:"signerAddr"`
	OwnerAddr      string `mapstructure:"ownerAddr"`
	Denom          string `mapstructure:"denom"`
	Exponent       int    `mapstructure:"exponent"`
	EndpointUrl    string `mapstructure:"endpointUrl"`
}

func GetConfig() *Config {
	if config == nil {
		once.Do(newConfig)
	}
	return config
}

func newConfig() {
	config = new(Config)

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("cannot read config file\n")
	}

	if err := v.Unmarshal(config); err != nil {
		log.Printf("cannot parse config file\n")
	}
}
