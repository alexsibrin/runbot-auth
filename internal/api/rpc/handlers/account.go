package handlers

import (
	"context"
	"google.golang.org/grpc"
	"runbot-auth/internal/api/models"
	"runbot-auth/pkg/runbotauthproto"
)

type IController interface {
	GetOne(ctx context.Context, uuid string) (*models.AccountGet, error)
}

type AccountDependencies struct {
	Controller IController
}

type Account struct {
	controller IController
	runbotauthproto.UnimplementedAccountServer
}

func NewAccount(d *AccountDependencies) (*Account, error) {
	return &Account{
		controller: d.Controller,
	}, nil
}

func (h *Account) Register(s *grpc.Server) {
	runbotauthproto.RegisterAccountServer(s, h)
}

func (h *Account) Get(ctx context.Context, model *runbotauthproto.GetAccount) (*runbotauthproto.GetAccountResponse, error) {
	modelout, err := h.controller.GetOne(ctx, model.UUID)
	if err != nil {
		return nil, err
	}
	response := h.convertAccountGetModelToResponse(modelout)
	return response, nil
}

func (h *Account) convertAccountGetModelToResponse(model *models.AccountGet) *runbotauthproto.GetAccountResponse {
	return &runbotauthproto.GetAccountResponse{
		UUID:  model.UUID,
		Name:  model.Name,
		Email: model.Email,
	}
}
