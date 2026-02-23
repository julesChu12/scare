package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 日志配置
type Config struct {
	Level       string // 日志级别: debug/info/warn/error
	Format      string // 输出格式: json/console
	OutputPath  string // 日志文件路径
	ErrorPath   string // 错误日志路径（分级存储）
	MaxSize     int    // 单文件最大大小(MB)
	MaxBackups  int    // 保留文件数
	MaxAge      int    // 保留天数
	Compress    bool   // 是否压缩
	EnableFile  bool   // 是否输出到文件
	EnableStdout bool  // 是否输出到控制台
}

var (
	defaultLogger *zap.SugaredLogger
	zapLogger     *zap.Logger
)

// Init 初始化日志系统
func Init(cfg Config) error {
	// 设置默认值
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.Format == "" {
		cfg.Format = "console"
	}
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 7
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30
	}

	// 解析日志级别
	level := parseLevel(cfg.Level)

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 根据格式选择编码器
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var cores []zapcore.Core

	// 控制台输出
	if cfg.EnableStdout || !cfg.EnableFile {
		consoleCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if cfg.EnableFile && cfg.OutputPath != "" {
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(cfg.OutputPath), 0755); err != nil {
			return err
		}

		// 文件编码器配置（无颜色）
		fileEncoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder, // 无颜色
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// 普通日志文件
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.OutputPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true,
		}
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncoderConfig), // 文件始终用 JSON，无颜色
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)

		// 错误日志单独文件（分级存储）
		if cfg.ErrorPath != "" {
			if err := os.MkdirAll(filepath.Dir(cfg.ErrorPath), 0755); err != nil {
				return err
			}
			errorWriter := &lumberjack.Logger{
				Filename:   cfg.ErrorPath,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
				LocalTime:  true,
			}
			errorCore := zapcore.NewCore(
				zapcore.NewJSONEncoder(fileEncoderConfig), // 无颜色
				zapcore.AddSync(errorWriter),
				zapcore.ErrorLevel, // 只记录 error 及以上
			)
			cores = append(cores, errorCore)
		}
	}

	// 创建 logger
	core := zapcore.NewTee(cores...)
	zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defaultLogger = zapLogger.Sugar()

	return nil
}

// InitSimple 简单初始化（仅控制台输出）
func InitSimple(level string) {
	Init(Config{
		Level:        level,
		Format:       "console",
		EnableStdout: true,
		EnableFile:   false,
	})
}

// parseLevel 解析日志级别
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Info 信息日志
func Info(format string, args ...any) {
	if defaultLogger == nil {
		InitSimple("info")
	}
	defaultLogger.Infof(format, args...)
}

// Error 错误日志
func Error(format string, args ...any) {
	if defaultLogger == nil {
		InitSimple("info")
	}
	defaultLogger.Errorf(format, args...)
}

// Sync 刷新日志缓冲
func Sync() {
	if zapLogger != nil {
		zapLogger.Sync()
	}
}
