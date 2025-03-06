package config

import "github.com/spf13/viper"

type Config struct {
	Environment        string `mapstructure:"ENVIRONMENT"`
	LogLevel           string `mapstructure:"LOG_LEVEL"`
	HTTPServerAddress  string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCUserAddress    string `mapstructure:"GRPC_USER_ADDRESS"`
	GRPCCatalogAddress string `mapstructure:"GRPC_CATALOG_ADDRESS"`
	GRPCOrderAddress   string `mapstructure:"GRPC_ORDER_ADDRESS"`
	GRPCPaymentAddress string `mapstructure:"GRPC_PAYMENT_ADDRESS"`
}

// LoadConfig reads configuration from file or enviroment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
