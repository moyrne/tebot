package configs

import "github.com/spf13/viper"

func LoadConfig(path ...string) error {
	dfPath := "configs/"
	if len(path) != 0 {
		dfPath = path[0]
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(dfPath)
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
