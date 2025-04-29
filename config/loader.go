package config

import (
	"github.com/spf13/viper"
)

func LoadConfig() {
	setDefaults()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/gh-aipr")
	_ = viper.ReadInConfig()
}
