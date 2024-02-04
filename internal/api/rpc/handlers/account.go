package handlers

import (
	"context"
	"errors"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/internal/logapp"
	"github.com/alexsibrin/runbot-auth/pkg/runbotauthproto"
	"google.golang.org/grpc"
)

const (
	accountKey = "accounts"
)

var (
	ErrDependenciesAreNil = errors.New("dependencies are nil")
	ErrControllerIsNil    = errors.New("controller is nil")
	ErrLoggerIsNil        = errors.New("logger is nil")
)

type IController interface {
	GetOneByEmail(ctx context.Context, email string) (*models.AccountGetModel, error)
	GetOneByUUID(ctx context.Context, uuid string) (*models.AccountGetModel, error)
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
	if d == nil {
		return nil, ErrDependenciesAreNil
	}
	if d.Controller == nil {
		return nil, ErrControllerIsNil
	}
	if d.Logger == nil {
		return nil, ErrLoggerIsNil
	}

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
	modelout, err := h.controller.GetOneByUUID(ctx, model.UUID)
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

func (h *Account) convertAccountGetModelToResponse(model *models.AccountGetModel) *runbotauthproto.GetAccountResponse {
	return &runbotauthproto.GetAccountResponse{
		UUID:  model.UUID,
		Name:  model.Name,
		Email: model.Email,
	}
}
