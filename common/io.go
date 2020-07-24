package common

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

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

func GetEnvDuration(name string, fallback time.Duration) time.Duration {
	if val, ok := os.LookupEnv(name); ok {
		if valUint, err := strconv.ParseUint(val, 10, 64); err == nil {
			return time.Duration(valUint) * time.Second
		}
	}

	return fallback
}

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
