package models

import "time"

type Submission struct {
	ID                string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SubmissionURL     string     `json:"submission_url"`
	AccountID         string     `json:"-"`
	Account           Account    `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	AssignmentID      string     `gorm:"type:uuid;" json:"assignment_id"`
	Assignment        Assignment `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	SubmissionCreated time.Time  `gorm:"autoCreateTime" json:"submission_date"`
	SubmissionUpdated time.Time  `gorm:"autoUpdateTime" json:"submission_updated"`
}
