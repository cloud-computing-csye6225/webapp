package services

import (
	"go.uber.org/zap"
	"webapp/db"
	"webapp/logger"
	"webapp/models"
)

type SubmissionsService struct {
	db *db.PostgresDB
}

func NewSubmissionsService(db *db.PostgresDB) *SubmissionsService {
	return &SubmissionsService{db}
}

func (as AccountsService) CreateSubmission(submission models.Submission) error {
	err := as.db.GetConnection().Create(&submission).Error
	if err != nil {
		logger.Error("failed to create a Submission", zap.Error(err))
		return err
	}
	return nil
}
