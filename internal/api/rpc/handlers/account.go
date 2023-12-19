package handlers

import "runbot-auth/internal/api/models"

type IController interface {
	GetAccount(uuid string) (*models.Account, error)
	UpdateAccount(update *models.UpdateAccount) error
}

type AccountService struct {
}
