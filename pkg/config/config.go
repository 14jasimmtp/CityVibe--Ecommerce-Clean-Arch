package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	AUTHTOKEN   string `mapstructure:"TWILIO_AUTH_TOKEN"`
	ACCOUNTSID  string `mapstructure:"TWILIO_ACCOUNT_SID"`
	SERVICESSID string `mapstructure:"VERIFY_SERVICE_SID"`

	KEY       string `mapstructure:"SECRETKEY"`

	Aws_region string `mapstructure:"AwsRegion"`
	Aws_access string `mapstructure:"AwsAccessKey"`
	Aws_secret string `mapstructure:"AwsSecretKey"`

	KEY_ID_FOR_PAY     string `mapstructure:"KEY_ID_PAY"`
	SECRET_KEY_FOR_PAY string `mapstructure:"KEY_SECRET_PAY"`
}

var envs = []string{
	 "DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", "TWILIO_AUTH_TOKEN", "TWILIO_ACCOUNT_SID", "VERIFY_SERVICE_SID", "SECRETKEY", "KEY_ID_PAY", "KEY_SECRET_PAY","AwsRegion","AwsAccessKey", "AwsSecretKey",
}

func LoadConfig() (Config, error) {
	var confg Config
	viper.AddConfigPath("/home/ubuntu/CityVibe-Project-Clean-Architecture/.env")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return confg, err
		}
	}
	if err := viper.Unmarshal(&confg); err != nil {
		return confg, err
	}
	if err := validator.New().Struct(&confg); err != nil {
		return confg, err
	}
	return confg, nil
}