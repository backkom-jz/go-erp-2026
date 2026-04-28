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
	MQ     MQConfig     `mapstructure:"mq"`
	AI     AIConfig     `mapstructure:"ai"`
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

type MQConfig struct {
	Enabled             bool   `mapstructure:"enabled"`
	URL                 string `mapstructure:"url"`
	Exchange            string `mapstructure:"exchange"`
	RoutingKey          string `mapstructure:"routing_key"`
	OrderTimeoutMinutes int    `mapstructure:"order_timeout_minutes"`
	OutboxMaxRetry      int    `mapstructure:"outbox_max_retry"`
	OutboxBaseDelaySeconds int `mapstructure:"outbox_base_delay_seconds"`
}

type AIConfig struct {
	Enabled        bool    `mapstructure:"enabled"`
	BaseURL        string  `mapstructure:"base_url"`
	APIKey         string  `mapstructure:"api_key"`
	Model          string  `mapstructure:"model"`
	TimeoutSeconds int     `mapstructure:"timeout_seconds"`
	Temperature    float64 `mapstructure:"temperature"`
	MaxTokens      int     `mapstructure:"max_tokens"`
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
	if strings.TrimSpace(cfg.JWT.Secret) == "" {
		return nil, fmt.Errorf("jwt.secret must not be empty")
	}
	if cfg.MQ.Exchange == "" {
		cfg.MQ.Exchange = "erp.events"
	}
	if cfg.MQ.RoutingKey == "" {
		cfg.MQ.RoutingKey = "erp.default"
	}
	if cfg.MQ.OrderTimeoutMinutes <= 0 {
		cfg.MQ.OrderTimeoutMinutes = 30
	}
	if cfg.MQ.OutboxMaxRetry <= 0 {
		cfg.MQ.OutboxMaxRetry = 6
	}
	if cfg.MQ.OutboxBaseDelaySeconds <= 0 {
		cfg.MQ.OutboxBaseDelaySeconds = 3
	}
	if strings.TrimSpace(cfg.AI.BaseURL) == "" {
		cfg.AI.BaseURL = "https://api.deepseek.com/chat/completions"
	}
	if strings.TrimSpace(cfg.AI.Model) == "" {
		cfg.AI.Model = "DeepSeekV4-pro"
	}
	if cfg.AI.TimeoutSeconds <= 0 {
		cfg.AI.TimeoutSeconds = 60
	}
	if cfg.AI.Temperature <= 0 {
		cfg.AI.Temperature = 0.7
	}
	if cfg.AI.MaxTokens <= 0 {
		cfg.AI.MaxTokens = 1024
	}
	return &cfg, nil
}
