package util

import "github.com/spf13/viper"

type Config struct {
	DriverSource  string `mapstructure:"driver_source"`
	DriverName    string `mapstructure:"driver_name"`
	ServerAddress string `mapstructure:"server_address"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return

}
