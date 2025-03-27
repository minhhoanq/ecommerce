package config

import "github.com/spf13/viper"

type Config struct {
	Environment          string `mapstructure:"ENVIRONMENT"`
	GRPCServerAddress    string `mapstructure:"GRPC_SERVER_ADDRESS"`
	GRPCUserAddress      string `mapstructure:"GRPC_USER_ADDRESS"`
	GRPCCatalogAddress   string `mapstructure:"GRPC_CATALOG_ADDRESS"`
	LogLevel             string `mapstructure:"LOG_LEVEL"`
	DBDriver             string `mapstructure:"DB_DRIVER"`
	DBSource             string `mapstructure:"DB_SOURCE"`
	DBHost               string `mapstructure:"DB_HOST"`
	DBPort               int    `mapstructure:"DB_PORT"`
	DBUser               string `mapstructure:"DB_USER"`
	DBPassword           string `mapstructure:"DB_PASSWORD"`
	DBName               string `mapstructure:"DB_NAME"`
	RedisAddress         string `mapstructure:"REDIS_ADDRESS"`
	EmailSenderName      string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string `mapstructure:"EMAIL_SENDER_PASSWORD"`
	KafkaBrokers         string `mapstructure:"KAFKA_BROKERS"`
	KafkaClientId        string `mapstructure:"KAFKA_CLIENT_ID"`
	KafkaConsumerGroupId string `mapstructure:"KAFKA_CONSUMER_GROUP_ID"`
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
