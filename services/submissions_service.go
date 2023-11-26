package services

import (
	"go.uber.org/zap"
	"time"
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

func (ss SubmissionsService) CheckSubmissionAttemptValidity(submission models.Submission, allowedAttempts int) (bool, error) {
	var submissions []models.Submission

	if err := ss.db.GetConnection().Where("assignment_id=?", submission.AssignmentID).Find(&submissions).Error; err != nil {
		logger.Error("Failed to check attempt validity", zap.Any("error", err))
		return false, err
	}

	if len(submissions)+1 > allowedAttempts {
		return false, nil
	} else {
		return true, nil
	}
}

func (ss SubmissionsService) CheckForLateSubmission(deadline time.Time) bool {
	currentTime := time.Now()
	return deadline.After(currentTime)
}
