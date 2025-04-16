package config

import (
	"github.com/spf13/viper"
	"strings"
)

var config *viper.Viper

func init() {
	viper.SetConfigName("config")

	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./../../configs")

	viper.AutomaticEnv()

	viper.SetEnvPrefix("TASK_MANAGEMENT_SERVICE")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		panic("Error initialising configs, " + err.Error())
	}
	config = viper.GetViper()
}

func AllSettings() map[string]interface{} {
	return config.AllSettings()
}

func GetInt32(key string) int32 {
	return config.GetInt32(key)
}

func GetStringSlice(key string) []string {
	return config.GetStringSlice(key)
}

func GetString(key string) string {
	return config.GetString(key)
}

func GetBool(key string) bool {
	return config.GetBool(key)
}
