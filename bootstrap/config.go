package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBUser           string `mapstructure:"DB_USER"`
	DBPass           string `mapstructure:"DB_PASS"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           string `mapstructure:"DB_PORT"`
	DBName           string `mapstructure:"DB_NAME"`
	ServerAddr       string `mapstructure:"SERVER_ADDRESS"`
	JWTSecret        string `mapstructure:"JWT_SECRET"`
	JWTExpire        int    `mapstructure:"JWT_EXPIRE"`
	JWTRefreshExpire int    `mapstructure:"JWT_REFRESH_EXPIRE"`
	AppInitialState  string `mapstructure:"APP_INITIAL_STATE"`

	TuitionType string `mapstructure:"TUITION_TYPE"`
	TuitionCost int    `mapstructure:"TUITION_COST"`
}

func NewConfig() *Config {
	conf := Config{}

	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalln(err)
	}


	return &conf
}
