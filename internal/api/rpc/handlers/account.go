package handlers

import (
	"context"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/internal/logapp"
	"github.com/alexsibrin/runbot-auth/pkg/runbotauthproto"
	"google.golang.org/grpc"
)

const (
	accountKey = "accounts"
)

type IController interface {
	GetOne(ctx context.Context, uuid string) (*models.AccountGetResponse, error)
	Create(ctx context.Context, model *models.AccountCreate) (*models.Account, error)
}

type AccountDependencies struct {
	Controller IController
	Logger     logapp.ILogger
}

type Account struct {
	controller IController
	logger     logapp.ILogger
	runbotauthproto.UnimplementedAccountServer
}

func NewAccount(d *AccountDependencies) (*Account, error) {
	// TODO: Add checking
	l := d.Logger.WithField(handlersKey, accountKey)
	return &Account{
		controller: d.Controller,
		logger:     l,
	}, nil
}

func (h *Account) Register(s *grpc.Server) {
	runbotauthproto.RegisterAccountServer(s, h)
}

func (h *Account) Get(ctx context.Context, model *runbotauthproto.GetAccount) (*runbotauthproto.GetAccountResponse, error) {
	modelout, err := h.controller.GetOne(ctx, model.UUID)
	if err != nil {
		h.handlerError(err)
		return nil, err
	}
	response := h.convertAccountGetModelToResponse(modelout)
	return response, nil
}

func (h *Account) Add(ctx context.Context, model *runbotauthproto.GetAccount) (*runbotauthproto.GetAccountResponse, error) {
	modelout, err := h.controller.GetOne(ctx, model.UUID)
	if err != nil {
		h.handlerError(err)
		return nil, err
	}
	response := h.convertAccountGetModelToResponse(modelout)
	return response, nil
}

func (h *Account) handlerError(err error) {
	h.logger.Error(err)
}

func (h *Account) convertAccountGetModelToResponse(model *models.AccountGetResponse) *runbotauthproto.GetAccountResponse {
	return &runbotauthproto.GetAccountResponse{
		UUID:  model.UUID,
		Name:  model.Name,
		Email: model.Email,
	}
}
