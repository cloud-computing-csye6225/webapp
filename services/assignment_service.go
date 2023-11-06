package services

import (
	"go.uber.org/zap"
	"webapp/db"
	"webapp/logger"
	"webapp/models"
)

type AssignmentService struct {
	db *db.PostgresDB
}

func NewAssignmentService(db *db.PostgresDB) *AssignmentService {
	return &AssignmentService{db}
}

func (as AssignmentService) AddAssignment(assignment *models.Assignment) error {
	err := as.db.GetConnection().Create(&assignment).Error
	if err != nil {
		logger.Error("failed to create an Assignment", zap.Error(err))
		return err
	}
	return nil
}

func (as AssignmentService) GetAssignment() ([]models.Assignment, error) {
	var assignments []models.Assignment

	if err := as.db.GetConnection().Find(&assignments).Error; err != nil {
		logger.Error("failed to get the Assignment", zap.Error(err))
		return assignments, err
	}
	return assignments, nil
}

func (as AssignmentService) GetAssignmentByID(assignmentID string) (models.Assignment, error) {
	var assignment models.Assignment

	if err := as.db.GetConnection().Where("id= ?", assignmentID).First(&assignment).Error; err != nil {
		logger.Error("failed to get the Assignment", zap.Error(err))
		return assignment, err
	}
	return assignment, nil
}

func (as AssignmentService) DeleteAssignment(assignment models.Assignment) error {
	if err := as.db.GetConnection().Delete(&assignment).Error; err != nil {
		logger.Error("failed to delete the Assignment", zap.Error(err))
		return err
	}
	return nil
}

func (as AssignmentService) UpdateAssignment(assignment models.Assignment) (models.Assignment, error) {
	oldAssignment, err := as.GetAssignmentByID(assignment.ID)
	if err != nil {
		logger.Error("failed to get the Assignment", zap.Error(err))
		return oldAssignment, err
	} else {
		updateError := as.db.GetConnection().Model(&oldAssignment).Updates(models.Assignment{
			Name:          assignment.Name,
			Points:        assignment.Points,
			NumOfAttempts: assignment.NumOfAttempts,
			Deadline:      assignment.Deadline,
			AccountID:     assignment.AccountID,
		}).Error

		if updateError != nil {
			logger.Error("failed to update the Assignment", zap.Error(err))
			return oldAssignment, updateError
		}
	}
	return oldAssignment, err
}
