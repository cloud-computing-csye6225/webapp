package services

import (
	"fmt"
	"webapp/config"
	"webapp/db"
)

type Services interface {
	LoadServices(configs config.Config)
}

type APIServices struct {
	Database          db.Database
	AssignmentService *AssignmentService
	AccountsService   *AccountsService
}

func (s *APIServices) LoadServices(configs config.Config) {
	d := &db.PostgresDB{}
	err := d.InitDatabase(configs.DBConfig)
	if err != nil {
		fmt.Printf("failed to initialize database: %s", err)
	}
	s.Database = d
	s.AccountsService = NewAccountService(d)
	s.AssignmentService = NewAssignmentService(d)
}
