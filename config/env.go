package config

import (
	"log"

	"github.com/spf13/viper"
)

type EnvConfigs struct {
	LocalServerPort string `mapstructure:"PORT"`
	DbConnect       string `mapstructure:"DSN"`
	JWT             string `mapstructure:"SECRET"`
	Twillio1        string `mapstructure:"TWILIO_ACCOUNT_SID"`
	Twillio2        string `mapstructure:"TWILIO_AUTH_TOKEN"`
	Twillio3        string `mapstructure:"VERIFY_SERVICE_SID"`
	StripeKey       string `mapstructure:"STRIPE_KEY"`
	RzpKey          string `mapstructure:"KEY_ID"`
	RzpSecret       string `mapstructure:"KEY_SECRET"`
}

var EnvConf *EnvConfigs

func InitEnvConfigs() {
	EnvConf = LoadEnvVariables()
}

func LoadEnvVariables() (configs *EnvConfigs) {
	viper.AddConfigPath(".")

	viper.SetConfigName("app")

	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error loading env env variables")
	}

	if err := viper.Unmarshal(&configs); err != nil {
		log.Fatal("Error while unmarshalling loaded variables into struct")
	}
	return
}
