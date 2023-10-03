package config

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     int
}

type ServerConfig struct {
	Host string
}

type Config struct {
	DBConfig     DatabaseConfig
	ServerConfig ServerConfig
}

func GetConfigs() *Config {
	return &Config{
		DBConfig: DatabaseConfig{
			DBUser:     getEnvVariable("DBUSER", ""),
			DBPassword: getEnvVariable("DBPASSWORD", ""),
			DBHost:     getEnvVariable("DBHOST", ""),
			DBName:     getEnvVariable("DBNAME", ""),
			DBPort:     getEnvVariableAsInt("DBPORT", 5000),
		},
		ServerConfig: ServerConfig{
			Host: getEnvVariable("SERVERPORT", ":8080"),
		},
	}
}

func getEnvVariable(key, defaultValue string) string {
	if ev, evExists := os.LookupEnv(key); evExists {
		return ev
	}
	return defaultValue
}

func getEnvVariableAsInt(key string, defaultValue int) int {
	if ev, evExists := os.LookupEnv(key); evExists {
		var atoi, err = strconv.Atoi(ev)
		if err == nil {
			return atoi
		}
	}
	return defaultValue
}
