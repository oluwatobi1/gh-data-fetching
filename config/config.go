package config

import (
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	DB_URL       string `mapstructure:"DB_URL"`
	PORT         string `mapstructure:"PORT"`
	ENVIRONMENT  string `mapstructure:"ENVIRONMENT"`
	GITHUB_TOKEN string `mapstructure:"GITHUB_TOKEN"`
	START_DATE   string `mapstructure:"START_DATE"`
	END_DATE     string `mapstructure:"END_DATE"`
	DEFAULT_REPO string `mapstructure:"DEFAULT_REPO"`
}

var Env *Config = &Config{}

func LoadConfig() error {

	var (
		err error
	)
	viper.AddConfigPath("./config")
	viper.SetConfigName("dev.env")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		// If the local config file is not found, ignore the error and fall back to OS env
		log.Printf("Warning: %s\n", err)
	}

	viper.AutomaticEnv()
	rc := reflect.ValueOf(Env).Elem()
	count := 0
	for i := 0; i < rc.NumField(); i++ {
		pName := reflect.TypeOf(Config{}).Field(i).Name
		log.Println(pName, viper.GetString(pName))
		rc.FieldByName(pName).SetString(viper.GetString(pName))
		count += len(viper.GetString(pName))
	}

	if count < 1 {
		log.Fatalln("Error: Missing ENV configs")
	}
	return nil
}
