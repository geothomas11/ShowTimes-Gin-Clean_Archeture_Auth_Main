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
		repository.NewCategoryRepository,
		repository.NewInventoryRepository,
		repository.NewCartRepository,
		repository.NewOrderRepository,

		usecase.NewUserUseCase,
		usecase.NewAdminUseCase,
		usecase.NewOtpUsecase,
		usecase.NewCategoryUseCase,
		usecase.NewInventoryUseCase,
		usecase.NewCartUseCase,
		usecase.NewOrderUseCase,

		handler.NewUserHandler,
		handler.NewAdminHandler,
		handler.NewOtpHandler,
		handler.NewCategoryHandler,
		handler.NewProductHandler,
		handler.NewCartHandler,
		handler.NewOrderHandler,

		helper.NewHelper,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
