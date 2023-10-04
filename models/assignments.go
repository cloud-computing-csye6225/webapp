package models

import "time"

type Assignment struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name              string    `json:"name" binding:"required"`
	Points            int       `json:"points" binding:"required,gte=1,lte=100"`
	NumOfAttempts     int       `json:"num_of_attempts" binding:"required,gte=1,lte=100"`
	Deadline          time.Time `json:"deadline" binding:"required"`
	AccountID         string    `json:"-"`
	Account           Account   `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	AssignmentCreated time.Time `gorm:"autoCreateTime" json:"assignment_created"`
	AssignmentUpdated time.Time `gorm:"autoUpdateTime" json:"assignment_updated"`
}
