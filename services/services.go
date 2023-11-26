package services

import (
	"go.uber.org/zap"
	"webapp/config"
	"webapp/db"
	"webapp/logger"
)

type Services interface {
	LoadServices(configs config.Config)
}

type APIServices struct {
	Database           db.Database
	AssignmentService  *AssignmentService
	AccountsService    *AccountsService
	SubmissionsService *SubmissionsService
	AWSService         *AWSService
}

func (s *APIServices) LoadServices(configs config.Config) {
	logger.Info("Initializing DB...")
	d := &db.PostgresDB{}
	err := d.InitDatabase(configs.DBConfig)
	if err != nil {
		logger.Error("failed to initialize database", zap.Error(err))
	}
	s.Database = d
	s.AccountsService = NewAccountService(d)
	s.AssignmentService = NewAssignmentService(d)
	s.SubmissionsService = NewSubmissionsService(d)
	s.AWSService = NewAWSService(&configs.AWSConfig)
}
