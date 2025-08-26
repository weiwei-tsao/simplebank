package util

import "github.com/spf13/viper"

// Config stores all configurations of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // look for .env file

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	err = viper.Unmarshal(&config)
	return
}
