package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// EnvPrefix is the default prefix for environment variables.
const DefaultEnvPrefix = "GZH"

// GetEnv returns the value of an environment variable with optional prefix.
func GetEnv(key string, prefix ...string) string {
	p := DefaultEnvPrefix
	if len(prefix) > 0 {
		p = prefix[0]
	}
	if p != "" {
		key = p + "_" + key
	}
	return os.Getenv(key)
}

// GetEnvOr returns the value of an environment variable or a default value.
func GetEnvOr(key, defaultValue string, prefix ...string) string {
	if v := GetEnv(key, prefix...); v != "" {
		return v
	}
	return defaultValue
}

// GetEnvBool returns the boolean value of an environment variable.
func GetEnvBool(key string, prefix ...string) bool {
	v := strings.ToLower(GetEnv(key, prefix...))
	return v == "true" || v == "1" || v == "yes" || v == "on"
}

// GetEnvBoolOr returns the boolean value or a default.
func GetEnvBoolOr(key string, defaultValue bool, prefix ...string) bool {
	v := GetEnv(key, prefix...)
	if v == "" {
		return defaultValue
	}
	v = strings.ToLower(v)
	return v == "true" || v == "1" || v == "yes" || v == "on"
}

// GetEnvInt returns the integer value of an environment variable.
func GetEnvInt(key string, prefix ...string) (int, bool) {
	v := GetEnv(key, prefix...)
	if v == "" {
		return 0, false
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return i, true
}

// GetEnvIntOr returns the integer value or a default.
func GetEnvIntOr(key string, defaultValue int, prefix ...string) int {
	if v, ok := GetEnvInt(key, prefix...); ok {
		return v
	}
	return defaultValue
}

// GetEnvDuration returns the duration value of an environment variable.
func GetEnvDuration(key string, prefix ...string) (time.Duration, bool) {
	v := GetEnv(key, prefix...)
	if v == "" {
		return 0, false
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, false
	}
	return d, true
}

// GetEnvDurationOr returns the duration value or a default.
func GetEnvDurationOr(key string, defaultValue time.Duration, prefix ...string) time.Duration {
	if v, ok := GetEnvDuration(key, prefix...); ok {
		return v
	}
	return defaultValue
}

// GetEnvList returns a list from a comma-separated environment variable.
func GetEnvList(key string, prefix ...string) []string {
	v := GetEnv(key, prefix...)
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// MustGetEnv returns the value of an environment variable or panics if not set.
func MustGetEnv(key string, prefix ...string) string {
	v := GetEnv(key, prefix...)
	if v == "" {
		p := DefaultEnvPrefix
		if len(prefix) > 0 {
			p = prefix[0]
		}
		fullKey := key
		if p != "" {
			fullKey = p + "_" + key
		}
		panic("required environment variable not set: " + fullKey)
	}
	return v
}

// LookupEnv returns the value and whether the environment variable is set.
func LookupEnv(key string, prefix ...string) (string, bool) {
	p := DefaultEnvPrefix
	if len(prefix) > 0 {
		p = prefix[0]
	}
	if p != "" {
		key = p + "_" + key
	}
	return os.LookupEnv(key)
}
