package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr        string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN    string `env:"DATABASE_URI"`
}

func ParseFlags() *Config {
	runAddr := flag.String("a", "localhost:8000", "...")
	accrualAddress := flag.String("r", "http://localhost:8080", "...")
	databaseDSN := flag.String("d", "", "...")

	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		runAddr = &envRunAddr
	}
	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		databaseDSN = &envDatabaseDSN
	}
	if accrualAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrualAddress != "" {
		databaseDSN = &accrualAddress
	}
	return &Config{
		RunAddr:        *runAddr,
		AccrualAddress: *accrualAddress,
		DatabaseDSN:    *databaseDSN,
	}
}

func NewConfig() *Config {
	return ParseFlags()
}
