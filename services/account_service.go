package services

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"webapp/db"
	"webapp/logger"
	"webapp/models"
)

type AccountsService struct {
	db *db.PostgresDB
}

func NewAccountService(db *db.PostgresDB) *AccountsService {
	return &AccountsService{db}
}

func (as AccountsService) AddAccount(account models.Account) error {
	err := as.db.GetConnection().Create(&account).Error
	if err != nil {
		logger.Error("failed to create an Account", zap.Error(err))
		return err
	}
	return nil
}

func (as AccountsService) GetAccountByEmail(email string) (models.Account, error) {
	var account models.Account

	if err := as.db.GetConnection().Where("email= ?", email).First(&account).Error; err != nil {
		logger.Error("failed to get the Account", zap.Error(err))
		return account, err
	}
	return account, nil
}

func (as AccountsService) GetAccountByID(accountID string) (models.Account, error) {
	var account models.Account

	if err := as.db.GetConnection().Where("id= ?", accountID).First(&account).Error; err != nil {
		logger.Error("failed to get the Account", zap.Error(err))
		return account, err
	}
	return account, nil
}

func (as AccountsService) GetAccounts() ([]models.Account, error) {
	var accounts []models.Account

	if err := as.db.GetConnection().Find(&accounts).Error; err != nil {
		logger.Error("failed to get Accounts", zap.Error(err))
		return accounts, err
	}
	return accounts, nil
}

func (as AccountsService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		logger.Error("failed to hash the password", zap.Error(err))
	}
	return string(bytes), err
}

func (as AccountsService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		logger.Error("unable to validate password", zap.Error(err))
	}
	return err == nil
}
