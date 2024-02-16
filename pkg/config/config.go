package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type MongoDBConfig struct {
	URI string `yaml:"uri" mapstructure:"uri"`
}

type RedisDBConfig struct {
	URI string `yaml:"uri" mapstructure:"uri"`
}

type RabbitMQConfig struct {
	URI string `yaml:"uri" mapstructure:"uri"`
}

type ServerConfig struct {
	Port int    `yaml:"port" mapstructure:"port"`
	Host string `yaml:"host" mapstructure:"host"`
}

type ExpressionServiceConfig struct {
	GourutinesCount  int    `yaml:"gourutines-count" mapstructure:"gourutines-count"`
	WorkerInfoUpdate int    `yaml:"worker-info-update" mapstructure:"worker-info-update"`
	ExpressionQueue  string `yaml:"expr-queue" mapstructure:"expr-queue"`
	ResultQueue      string `yaml:"res-queue" mapstructure:"res-queue"`
	Exchange         string `yaml:"exchange" mapstructure:"exchange"`
	RouteKey         string `yaml:"route-key" mapstructure:"route-key"`
}

type LoggerConfig struct {
	Level      string `yaml:"level" mapstructure:"level"`
	Path       string `yaml:"path" mapstructure:"path"`
	MaxSize    int    `yaml:"max-size" mapstructure:"max-size"`
	MaxAge     int    `yaml:"max-age" mapstructure:"max-age"`
	MaxBackups int    `yaml:"max-backups" mapstructure:"max-backups"`
}

type AppConfig struct {
	PasswordMinLength int    `mapstructure:"password-min-length" yaml:"password-min-length"`
	AccessTokenTTL    int    `mapstructure:"access-token-ttl" yaml:"access-token-ttl"`   // in minutes
	RefreshTokenTTL   int    `mapstructure:"refresh-token-ttl" yaml:"refresh-token-ttl"` // in minutes
	TokenSecret       string `mapstructure:"token-secret" yaml:"token-secret"`
	CacheTTL          int    `mapstructure:"cache-ttl" yaml:"cache-ttl"`
}

type Config struct {
	Logger *LoggerConfig            `yaml:"logger" mapstructure:"logger"`
	Expr   *ExpressionServiceConfig `yaml:"expression-service" mapstructure:"expression-service"`
	Server *ServerConfig            `yaml:"server" mapstructure:"server"`
	Mongo  *MongoDBConfig           `yaml:"mongodb" mapstructure:"mongodb"`
	Redis  *RedisDBConfig           `yaml:"redis" mapstructure:"redis"`
	App    *AppConfig               `yaml:"app" mapstructure:"app"`
	Rabbit *RabbitMQConfig          `yaml:"rabbit" mapstructure:"rabbit"`
}

func NewConfig(configPath string) (*Config, error) {
	v := viper.New()
	setDefaults(v)

	var config = new(Config)
	err := readConfig(v, configPath)

	if err := v.Unmarshal(config); err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	return config, err
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")

	v.SetDefault("logger.level", "debug")
	v.SetDefault("logger.path", "./log/logs")
	v.SetDefault("logger.max-size", 100)
	v.SetDefault("logger.max-age", 30)
	v.SetDefault("logger.max-backups", 3)

	v.SetDefault("mongo.uri", "localhost:27017")

	v.SetDefault("redis.uri", "localhost:6379")

	v.SetDefault("app.password-min-length", 8)
	v.SetDefault("app.access-token-ttl", 15)
	v.SetDefault("app.refresh-token-ttl", 60)
	v.SetDefault("app.token-secret", "secret")
	v.SetDefault("app.cache-ttl", 10)

	v.SetDefault("expression-service.expr-queue", "input")
	v.SetDefault("expression-service.res-queue", "output")
	v.SetDefault("expression-service.gourutines-count", 10)
	v.SetDefault("expression-service.exchange", "exprService")
	v.SetDefault("expression-service.route-key", "rpc")
	v.SetDefault("expression-service.worker-info-update", 5)

	v.SetDefault("rabbit.uri", "amqp://guest:guest@localhost:5672/")
}

func readConfig(v *viper.Viper, configPath string) error {
	v.AddConfigPath(configPath)
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
