package config

import (
	"os"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	//Database info
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	AUTHTOKEN  string `mapstructure:"DB_AUTHTOKEN"`
	ACCOUNTSID string `mapstructure:"DB_ACCOUNTSID"`
	SERVICESID string `mapstructure:"DB_SERVICESID"`

	Admin_AccessKey  string `mapstructure:"AdminAccessKey"`
	Admin_RefreshKey string `mapstructure:"AdminRefreshKey"`

	User_AccessKey  string `mapstructure:"UserAccessKey"`
	User_RefreshKey string `mapstructure:"UserRefreshKey"`
	//AWS S3 bucket
	AWSRegion          string `mapstructure:"AWSRegion"`
	AWSAccesskeyID     string `mapstructure:"AWSAccesskeyID"`
	AWSSecretaccesskey string `mapstructure:"AWSSecretaccesskey"`
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD",
	"DB_AUTHTOKEN", "DB_ACCOUNTSID", "DB_SERVICESID",
	"AdminAccessKey", "AdminRefreshKey",
	"UserAccessKey", "UserRefreshKey",
	"AWSRegion", "AWSAccesskeyID", "AWSSecretaccesskey",
}

type ConfigAuth struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig ConfigAuth

func LoadConfig() (Config, error) {
	var config Config
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}
	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}
	GoogleConfig()
	return config, nil

}

//Google loggin config

func GoogleConfig() oauth2.Config {
	AppConfig.GoogleLoginConfig = oauth2.Config{
		ClientID:     os.Getenv("Auth2ClientID"),
		ClientSecret: os.Getenv("Auth2ClientSecret"),
		RedirectURL:  "http://localhost:7000/user/google_callback",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig
}
