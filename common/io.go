package common

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// GetEnvBool returns a named boolean environment variable, or a fallback value if it doesn't exist
func GetEnvBool(name string, fallback bool) bool {
	if val, ok := os.LookupEnv(name); ok {
		switch val {
		case "true", "yes", "1", "":
			return true
		case "false", "no", "0":
			return false
		}
	}

	return fallback
}

// GetEnvDuration returns a named duration environment variable, or a fallback value if it doesn't exist
func GetEnvDuration(name string, fallback time.Duration) time.Duration {
	if val, ok := os.LookupEnv(name); ok {
		if valUint, err := strconv.ParseUint(val, 10, 64); err == nil {
			return time.Duration(valUint) * time.Second
		}
	}

	return fallback
}

// GetEnvString returns a named string environment variable, or a fallback value if it doesn't exist
func GetEnvString(name, fallback string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}

	return fallback
}

// IsFileExists checks if a file exists and is not a directory.
func IsFileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ReadEntireFile reads the entire file and returns its contents
func ReadEntireFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	return ioutil.ReadAll(file)
}
