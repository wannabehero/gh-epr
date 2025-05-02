package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func LoadConfig() {
	setDefaults()

	globalViper := viper.New()
	globalViper.SetConfigName("config")
	globalViper.SetConfigType("yaml")
	globalViper.AddConfigPath("$HOME/.config/gh-aipr")
	_ = globalViper.ReadInConfig()

	for _, key := range globalViper.AllKeys() {
		viper.Set(key, globalViper.Get(key))
	}

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	localConfigPath := filepath.Join(cwd, ".github", "gh-aipr", "config.yaml")
	if _, err := os.Stat(localConfigPath); err == nil {
		localViper := viper.New()
		localViper.SetConfigFile(localConfigPath)
		if err := localViper.ReadInConfig(); err == nil {
			for _, key := range localViper.AllKeys() {
				viper.Set(key, localViper.Get(key))
			}
		}
	}
}
