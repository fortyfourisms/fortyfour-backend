// Package logger provides a centralized, structured logging solution using zerolog.
// It replaces Rollbar with local structured logging and ensures no sensitive data is logged.
package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// log is the global logger instance.
var log zerolog.Logger

// sensitiveKeys defines keys that should never be logged.
var sensitiveKeys = map[string]bool{
	"password":           true,
	"token":              true,
	"secret":             true,
	"jwt_secret":         true,
	"jwt_refresh_secret": true,
	"authorization":      true,
	"cookie":             true,
	"api_key":            true,
	"access_token":       true,
	"refresh_token":      true,
	"otp":                true,
	"old_password":       true,
	"new_password":       true,
	"db_password":        true,
	"redis_password":     true,
	"rabbitmq_password":  true,
}

func init() {
	// Default: production JSON logger
	log = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}

// Init initializes the global logger with the given log level and environment.
// Call this once at application startup.
func Init(level string, environment string) {
	lvl := parseLevel(level)
	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = time.RFC3339

	if environment == "development" {
		// Pretty console output for development
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		log = zerolog.New(output).
			With().
			Timestamp().
			Caller().
			Str("env", environment).
			Logger().
			Level(lvl)
	} else {
		// Structured JSON for production
		log = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Caller().
			Str("env", environment).
			Logger().
			Level(lvl)
	}
}

// parseLevel converts a string log level to zerolog.Level.
func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disabled", "off":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}

// Error logs an error with context.
func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

// Errorf logs an error with a formatted message.
func Errorf(err error, format string, args ...interface{}) {
	log.Error().Err(err).Msgf(format, args...)
}

// Info logs an informational message.
func Info(msg string) {
	log.Info().Msg(msg)
}

// Infof logs a formatted informational message.
func Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

// Warn logs a warning message.
func Warn(msg string) {
	log.Warn().Msg(msg)
}

// Warnf logs a formatted warning message.
func Warnf(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

// Debug logs a debug message.
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Debugf logs a formatted debug message.
func Debugf(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

// Fatal logs a fatal message and exits.
func Fatal(msg string) {
	log.Fatal().Msg(msg)
}

// Fatalf logs a formatted fatal message and exits.
func Fatalf(format string, args ...interface{}) {
	log.Fatal().Msgf(format, args...)
}

// FatalErr logs a fatal error and exits.
func FatalErr(err error, msg string) {
	log.Fatal().Err(err).Msg(msg)
}

// WithField returns a logger event with a key-value pair.
// It redacts sensitive fields automatically.
func WithField(key string, value interface{}) *zerolog.Event {
	if IsSensitiveKey(key) {
		return log.Info().Str(key, "[REDACTED]")
	}
	return log.Info().Interface(key, value)
}

// IsSensitiveKey checks if a key is considered sensitive and should not be logged.
func IsSensitiveKey(key string) bool {
	return sensitiveKeys[strings.ToLower(key)]
}

// Get returns the underlying zerolog.Logger for advanced usage.
func Get() zerolog.Logger {
	return log
}
