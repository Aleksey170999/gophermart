package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr        string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN    string `env:"DATABASE_URI"`
	JWTSecret      string `env:"JWT_SECRET"`
}

func ParseFlags() *Config {
	runAddr := flag.String("a", "localhost:8000", "Server address and port to listen on")
	accrualAddress := flag.String("r", "http://localhost:8080", "Accrual system address")
	databaseDSN := flag.String("d", "", "Database connection string")
	jwtSecret := flag.String("j", "secret", "JWT secret key")

	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		runAddr = &envRunAddr
	}
	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		databaseDSN = &envDatabaseDSN
	}
	if envAccrualAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualAddress != "" {
		accrualAddress = &envAccrualAddress
	}
	if envJWTSecret := os.Getenv("JWT_SECRET"); envJWTSecret != "" {
		jwtSecret = &envJWTSecret
	}
	if *jwtSecret == "" {
		panic("JWT_SECRET environment variable or -jwt-secret flag must be provided")
	}

	return &Config{
		RunAddr:        *runAddr,
		AccrualAddress: *accrualAddress,
		DatabaseDSN:    *databaseDSN,
		JWTSecret:      *jwtSecret,
	}
}

func NewConfig() *Config {
	return ParseFlags()
}
