package config

import (
	"go.uber.org/zap"
	"os"
	"strconv"
	"webapp/logger"
)

type DatabaseConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     int
}

type ServerConfig struct {
	Host    string
	GinMode string
}

type DefaultUsers struct {
	Path string
}

type Config struct {
	DBConfig     DatabaseConfig
	ServerConfig ServerConfig
	DefaultUsers DefaultUsers
}

func GetConfigs() Config {
	logger.Info("Getting configs from environment")
	return Config{
		DBConfig: DatabaseConfig{
			DBUser:     getEnvVariable("DBUSER", ""),
			DBPassword: getEnvVariable("DBPASSWORD", ""),
			DBHost:     getEnvVariable("DBHOST", ""),
			DBName:     getEnvVariable("DBNAME", ""),
			DBPort:     getEnvVariableAsInt("DBPORT", 5000),
		},
		ServerConfig: ServerConfig{
			Host:    getEnvVariable("SERVERPORT", ":8080"),
			GinMode: getEnvVariable("GIN_MODE", "debug"),
		},
		DefaultUsers: DefaultUsers{
			Path: getEnvVariable("DEFAULTUSERS", ""),
		},
	}
}

func getEnvVariable(key, defaultValue string) string {
	logger.Info("Getting env", zap.Any("key", key))
	if ev, evExists := os.LookupEnv(key); evExists {
		return ev
	}
	return defaultValue
}

func getEnvVariableAsInt(key string, defaultValue int) int {
	logger.Info("Getting env", zap.Any("key", key))
	if ev, evExists := os.LookupEnv(key); evExists {
		var atoi, err = strconv.Atoi(ev)
		if err == nil {
			return atoi
		}
	}
	return defaultValue
}
