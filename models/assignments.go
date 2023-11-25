package models

import "time"

type Assignment struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name              string    `json:"name"`
	Points            int       `json:"points"`
	NumOfAttempts     int       `json:"num_of_attempts"`
	Deadline          time.Time `json:"deadline"`
	AccountID         string    `json:"-"`
	Account           Account   `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	AssignmentCreated time.Time `gorm:"autoCreateTime" json:"assignment_created"`
	AssignmentUpdated time.Time `gorm:"autoUpdateTime" json:"assignment_updated"`
}
