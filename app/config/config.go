package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	App      App                       `mapstructure:"app"`
	Database Database                  `mapstructure:"database"`
	Cosmos   map[string]CosmosAppchain `mapstructure:"cosmos"`
	Aptos    Aptos                     `mapstructure:"aptos"`
	Polygon  Polygon                   `mapstructure:"polygon"`
}

type App struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Cron string `mapstructure:"cron"`
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
	CoingeckoId           string   `mapstructure:"coingeckoId"`
	ValidatorOperatorAddr string   `mapstructure:"validatorOperatorAddr"`
	ValidatorAddr         string   `mapstructure:"validatorAddr"`
	Exponent              int      `mapstructure:"exponent"`
	GrantAddrs            []string `mapstructure:"grantAddrs"`
}

type Aptos struct {
	validatorAddr string
}

type Polygon struct {
	validatorAddr string
}

func GetConfig() *Config {
	if config == nil {
		config = newConfig()
	}
	return config
}

func newConfig() *Config {
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

	return config
}
