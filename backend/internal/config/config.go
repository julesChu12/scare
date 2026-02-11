package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Mail     MailConfig     `mapstructure:"mail"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	Name string `mapstructure:"name"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	ExpiresIn        int    `mapstructure:"expires_in"`
	RefreshExpiresIn int    `mapstructure:"refresh_expires_in"`
}

type LogConfig struct {
	Level        string `mapstructure:"level"`
	Format       string `mapstructure:"format"`
	OutputPath   string `mapstructure:"output_path"`
	ErrorPath    string `mapstructure:"error_path"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups"`
	MaxAge       int    `mapstructure:"max_age"`
	Compress     bool   `mapstructure:"compress"`
	EnableFile   bool   `mapstructure:"enable_file"`
	EnableStdout bool   `mapstructure:"enable_stdout"`
}

type MailConfig struct {
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
	SMTPUser string `mapstructure:"smtp_user"`
	SMTPPass string `mapstructure:"smtp_pass"`
	SMTPTLS  bool   `mapstructure:"smtp_tls"`
}

type StorageConfig struct {
	Driver                 string             `mapstructure:"driver"`
	SignedURLExpireSeconds int                `mapstructure:"signed_url_expire_seconds"`
	Local                  LocalStorageConfig `mapstructure:"local"`
	OSS                    OSSStorageConfig   `mapstructure:"oss"`
}

type LocalStorageConfig struct {
	BasePath string `mapstructure:"base_path"`
	BaseURL  string `mapstructure:"base_url"`
}

type OSSStorageConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Region          string `mapstructure:"region"`
	BaseURL         string `mapstructure:"base_url"`
}

func Load() (*Config, error) {
	if err := loadEnvFile(".env"); err != nil {
		return nil, fmt.Errorf("load .env failed: %w", err)
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnv(v)

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("read config failed: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	applyDefaults(&cfg)

	return &cfg, nil
}

func bindEnv(v *viper.Viper) {
	bindings := map[string]string{
		"server.port": "SERVER_PORT",
		"server.mode": "SERVER_MODE",
		"server.name": "SERVER_NAME",

		"database.host":           "DB_HOST",
		"database.port":           "DB_PORT",
		"database.user":           "DB_USER",
		"database.password":       "DB_PASSWORD",
		"database.dbname":         "DB_NAME",
		"database.charset":        "DB_CHARSET",
		"database.max_idle_conns": "DB_MAX_IDLE_CONNS",
		"database.max_open_conns": "DB_MAX_OPEN_CONNS",

		"redis.host":      "REDIS_HOST",
		"redis.port":      "REDIS_PORT",
		"redis.password":  "REDIS_PASSWORD",
		"redis.db":        "REDIS_DB",
		"redis.pool_size": "REDIS_POOL_SIZE",

		"jwt.secret":             "JWT_SECRET",
		"jwt.expires_in":         "JWT_EXPIRES_IN",
		"jwt.refresh_expires_in": "JWT_REFRESH_EXPIRES_IN",

		"log.level":         "LOG_LEVEL",
		"log.format":        "LOG_FORMAT",
		"log.output_path":   "LOG_OUTPUT_PATH",
		"log.error_path":    "LOG_ERROR_PATH",
		"log.max_size":      "LOG_MAX_SIZE",
		"log.max_backups":   "LOG_MAX_BACKUPS",
		"log.max_age":       "LOG_MAX_AGE",
		"log.compress":      "LOG_COMPRESS",
		"log.enable_file":   "LOG_ENABLE_FILE",
		"log.enable_stdout": "LOG_ENABLE_STDOUT",

		"mail.smtp_host": "SMTP_HOST",
		"mail.smtp_port": "SMTP_PORT",
		"mail.smtp_user": "SMTP_USER",
		"mail.smtp_pass": "SMTP_PASS",
		"mail.smtp_tls":  "SMTP_TLS",

		"storage.driver":                    "STORAGE_DRIVER",
		"storage.signed_url_expire_seconds": "STORAGE_SIGNED_URL_EXPIRE_SECONDS",
		"storage.local.base_path":           "STORAGE_LOCAL_BASE_PATH",
		"storage.local.base_url":            "STORAGE_LOCAL_BASE_URL",
		"storage.oss.endpoint":              "STORAGE_OSS_ENDPOINT",
		"storage.oss.bucket":                "STORAGE_OSS_BUCKET",
		"storage.oss.access_key_id":         "STORAGE_OSS_ACCESS_KEY_ID",
		"storage.oss.access_key_secret":     "STORAGE_OSS_ACCESS_KEY_SECRET",
		"storage.oss.region":                "STORAGE_OSS_REGION",
		"storage.oss.base_url":              "STORAGE_OSS_BASE_URL",
	}

	for key, env := range bindings {
		_ = v.BindEnv(key, env)
	}
}

func applyDefaults(cfg *Config) {
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "debug"
	}
	if cfg.Server.Name == "" {
		cfg.Server.Name = "scare-backend"
	}

	if cfg.Database.Port == 0 {
		cfg.Database.Port = 3306
	}
	if cfg.Database.Charset == "" {
		cfg.Database.Charset = "utf8mb4"
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 10
	}
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 100
	}

	if cfg.Redis.Port == 0 {
		cfg.Redis.Port = 6379
	}
	if cfg.Redis.PoolSize == 0 {
		cfg.Redis.PoolSize = 100
	}

	if cfg.JWT.ExpiresIn == 0 {
		cfg.JWT.ExpiresIn = 24
	}
	if cfg.JWT.RefreshExpiresIn == 0 {
		cfg.JWT.RefreshExpiresIn = 168
	}

	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "console"
	}
	if cfg.Log.OutputPath == "" {
		cfg.Log.OutputPath = "logs/app.log"
	}
	if cfg.Log.ErrorPath == "" {
		cfg.Log.ErrorPath = "logs/error.log"
	}
	if cfg.Log.MaxSize == 0 {
		cfg.Log.MaxSize = 100
	}
	if cfg.Log.MaxBackups == 0 {
		cfg.Log.MaxBackups = 7
	}
	if cfg.Log.MaxAge == 0 {
		cfg.Log.MaxAge = 30
	}

	if cfg.Storage.SignedURLExpireSeconds == 0 {
		cfg.Storage.SignedURLExpireSeconds = 3600
	}
	if cfg.Storage.Driver == "" {
		cfg.Storage.Driver = "local"
	}
	if cfg.Storage.Local.BasePath == "" {
		cfg.Storage.Local.BasePath = "./storage"
	}
	if cfg.Storage.Local.BaseURL == "" {
		cfg.Storage.Local.BaseURL = "http://localhost:8080/static"
	}
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		idx := strings.Index(line, "=")
		if idx == -1 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		if key == "" {
			continue
		}

		value = trimQuotes(value)
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, value)
	}

	return scanner.Err()
}

func trimQuotes(value string) string {
	if len(value) < 2 {
		return value
	}
	if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
		return value[1 : len(value)-1]
	}
	return value
}
