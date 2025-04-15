package configs

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

func InitLogger(config *Config) *zerolog.Logger {
	// Parse log level
	level, err := zerolog.ParseLevel(config.Logger.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel // Default to INFO if parsing fails
	}
	zerolog.SetGlobalLevel(level) // Set global log level

	// Declare output as io.Writer to support both os.Stdout and ConsoleWriter
	var output io.Writer = os.Stdout

	// Pretty-print logs for non-production environments
	if config.AppEnv != "production" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("[%s]", i))
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("| %s |", i)
			},
			PartsExclude: []string{zerolog.TimestampFieldName},
		}
	}

	// Initialize and assign the global logger
	loggerBuilder := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()

	// if config.AppEnv == "production" {
	// 	loggerBuilder = loggerBuilder.Hook(&hook.LoggerHook{})
	// 	fmt.Println("✅ Logger hook enabled in production mode")
	// }

	// Assign the global logger
	Logger = &loggerBuilder
	return Logger
}
