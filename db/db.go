package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"webapp/config"
)

type Database interface {
	InitDatabase(config.DatabaseConfig) error
	Ping() error
	GetConnection() *gorm.DB
}

type PostgresDB struct {
	db *gorm.DB
}

func (p *PostgresDB) InitDatabase(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err == nil {
		p.db = db
	}
	return err
}

func (p *PostgresDB) GetConnection() *gorm.DB {
	return p.db
}

func (p *PostgresDB) Ping() error {
	if p.db == nil {
		fmt.Printf("Database connection is not active, trying to reconnect to DB....\n")
		configs := config.GetConfigs()
		err := p.InitDatabase(configs.DBConfig)
		if err != nil {
			fmt.Printf("unable to ping database, database connection is not established\n%s\n", err)
			return err
		}
		fmt.Printf("Database connection established successfully\n")
	}

	sqlDB, err := p.db.DB()

	if err != nil {
		fmt.Printf("error getting generic SQL from Gorm...\n %s", err)
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		fmt.Printf("unable to ping database %s\n", err)
		return err
	}

	return err
}
