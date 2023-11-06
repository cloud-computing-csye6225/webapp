package db

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"webapp/config"
	"webapp/logger"
	"webapp/models"
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
	//Create database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%d",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBPort,
	)
	initDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})

	logger.Info("Creating DB if not exists")
	result := initDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if result.Error != nil {
		logger.Error("Unable to create database", zap.Error(result.Error))
	}

	dsn = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err == nil {
		p.db = db
		automigrateError := p.db.AutoMigrate(&models.Account{}, &models.Assignment{})
		if automigrateError != nil {
			logger.Error("Unable to automigrate schemas", zap.Error(automigrateError))
			return automigrateError
		}
	}
	return err
}

func (p *PostgresDB) GetConnection() *gorm.DB {
	return p.db
}

func (p *PostgresDB) Ping() error {
	if p.db == nil {
		logger.Info("Database connection is not active, trying to reconnect to DB...")
		configs := config.GetConfigs()
		err := p.InitDatabase(configs.DBConfig)
		if err != nil {
			logger.Error("unable to ping database, database connection is not established", zap.Error(err))
			return err
		}
		logger.Info("Database connection established successfully")
	}

	sqlDB, err := p.db.DB()

	if err != nil {
		logger.Error("error getting generic SQL from Gorm...", zap.Error(err))
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		logger.Error("unable to ping database", zap.Error(err))
		return err
	}

	return err
}
