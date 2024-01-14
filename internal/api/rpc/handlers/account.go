package handlers

import (
	"context"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/pkg/runbotauthproto"
	"google.golang.org/grpc"
)

type IController interface {
	GetOne(ctx context.Context, uuid string) (*models.AccountGetResponse, error)
	Create(ctx context.Context, model *models.AccountCreate) (*models.Account, error)
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

func (h *Account) Add(ctx context.Context, model *runbotauthproto.GetAccount) (*runbotauthproto.GetAccountResponse, error) {
	modelout, err := h.controller.GetOne(ctx, model.UUID)
	if err != nil {
		return nil, err
	}
	response := h.convertAccountGetModelToResponse(modelout)
	return response, nil
}

func (h *Account) convertAccountGetModelToResponse(model *models.AccountGetResponse) *runbotauthproto.GetAccountResponse {
	return &runbotauthproto.GetAccountResponse{
		UUID:  model.UUID,
		Name:  model.Name,
		Email: model.Email,
	}
}
