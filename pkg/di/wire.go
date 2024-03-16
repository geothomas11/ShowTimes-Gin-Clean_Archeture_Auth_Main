//go:build wireinject
// +build wireinject

package di

import (
	http "ShowTimes/pkg/api"
	"ShowTimes/pkg/api/handler"
	"ShowTimes/pkg/config"
	"ShowTimes/pkg/db"
	"ShowTimes/pkg/helper"
	"ShowTimes/pkg/repository"
	"ShowTimes/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(

		db.ConectDatabse,

		repository.NewUserRepository,
		repository.NewAdminRepository,
		repository.NewOtpRepository,

		usecase.NewUserUseCase,
		usecase.NewAdminUseCase,
		usecase.NewOtpUsecase,

		handler.NewUserHandler,
		handler.NewAdminHandler,
		handler.NewOtpHandler,

		helper.NewHelper,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
