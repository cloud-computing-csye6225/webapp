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

func (ss SubmissionsService) CreateSubmission(submission *models.Submission) error {
	err := ss.db.GetConnection().Create(&submission).Error
	if err != nil {
		logger.Error("failed to create a Submission", zap.Error(err))
		return err
	}
	return nil
}

func (ss SubmissionsService) GetSubmissionByID(submissionId string) (models.Submission, error) {
	var submission models.Submission

	if err := ss.db.GetConnection().Where("id= ?", submissionId).First(&submission).Error; err != nil {
		logger.Error("failed to get the Assignment", zap.Error(err))
		return submission, err
	}
	return submission, nil
}
