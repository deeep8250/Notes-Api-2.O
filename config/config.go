package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port   int
	DBDSN  string
	AppEnv string
}

func Load() *Config {
	//for port
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("invalid port : ", err.Error())
	}

	// for dsn
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB DSN is required")
	}

	// APP_ENV (default "dev")
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	return &Config{
		Port:   port,
		DBDSN:  dsn,
		AppEnv: env,
	}

}
