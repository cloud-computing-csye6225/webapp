package services

import (
	"fmt"
	"gorm.io/gorm"
	"webapp/db"
	"webapp/models"
)

type AssignmentService struct {
	db *gorm.DB
}

func NewAssignmentService(db *db.PostgresDB) *AssignmentService {
	return &AssignmentService{db.GetConnection()}
}

func (as AssignmentService) AddAssignment(assignment *models.Assignment) error {
	err := as.db.Create(&assignment).Error
	if err != nil {
		fmt.Printf("Failed to create an Assignment, %s\n", err)
		return err
	}
	return nil
}

func (as AssignmentService) GetAssignment() ([]models.Assignment, error) {
	var assignments []models.Assignment

	if err := as.db.Find(&assignments).Error; err != nil {
		fmt.Printf("Failed to get an Accounts, %s\n", err)
		return assignments, err
	}
	return assignments, nil
}

func (as AssignmentService) GetAssignmentByID(assignmentID string) (models.Assignment, error) {
	var assignment models.Assignment

	if err := as.db.Where("id= ?", assignmentID).First(&assignment).Error; err != nil {
		fmt.Printf("Failed to get an Assignment, %s\n", err)
		return assignment, err
	}
	return assignment, nil
}

func (as AssignmentService) DeleteAssignment(assignment models.Assignment) error {
	if err := as.db.Delete(&assignment).Error; err != nil {
		fmt.Printf("Failed to delete the assignment, %s\n", err)
		return err
	}
	return nil
}

func (as AssignmentService) UpdateAssignment(assignment models.Assignment) (models.Assignment, error) {
	oldAssignment, err := as.GetAssignmentByID(assignment.ID)
	if err != nil {
		fmt.Printf("Failed to get assignment with ID %s to delete, %s\n", assignment.ID, err)
		return oldAssignment, err
	} else {
		updateError := as.db.Model(&oldAssignment).Updates(models.Assignment{
			Name:          assignment.Name,
			Points:        assignment.Points,
			NumOfAttempts: assignment.NumOfAttempts,
			Deadline:      assignment.Deadline,
			AccountID:     assignment.AccountID,
		}).Error

		if updateError != nil {
			fmt.Printf("Failed to update the assignment, %s\n", err)
			return oldAssignment, updateError
		}
	}
	return oldAssignment, err
}
