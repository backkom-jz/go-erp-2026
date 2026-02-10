package bootstrap

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
	DB     DBConfig     `mapstructure:"db"`
	Redis  RedisConfig  `mapstructure:"redis"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

type ServerConfig struct {
	Addr string `mapstructure:"addr"`
	Mode string `mapstructure:"mode"`
}

type LogConfig struct {
	Level    string   `mapstructure:"level"`
	Encoding string   `mapstructure:"encoding"`
	Output   []string `mapstructure:"output"`
}

type DBConfig struct {
	Driver                string `mapstructure:"driver"`
	DSN                   string `mapstructure:"dsn"`
	MaxIdleConns          int    `mapstructure:"max_idle_conns"`
	MaxOpenConns          int    `mapstructure:"max_open_conns"`
	ConnMaxLifetimeSecond int    `mapstructure:"conn_max_lifetime_seconds"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Enabled  bool   `mapstructure:"enabled"`
}

type JWTConfig struct {
	Secret            string `mapstructure:"secret"`
	AccessTTLMinutes  int    `mapstructure:"access_ttl_minutes"`
	RefreshTTLMinutes int    `mapstructure:"refresh_ttl_minutes"`
}

func LoadConfig() (*Config, error) {
	env := strings.TrimSpace(os.Getenv("APP_ENV"))
	if env == "" {
		env = "dev"
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	if err := v.ReadInConfig(); err != nil {
		v.SetConfigName("config")
		if err2 := v.ReadInConfig(); err2 != nil {
			return nil, fmt.Errorf("read config failed: %w", err)
		}
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}
	return &cfg, nil
}
