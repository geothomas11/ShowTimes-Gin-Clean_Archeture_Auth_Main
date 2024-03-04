package di

import (
	"ShowTimes/pkg/config"
	"net/http"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(

		db.ConnectDatabase,

		repository.NewUserRepository,
		repository.NewAdminRepository,

		usecase.NewUserUseCase,
		usecase.NewAdminUseCase,

		handler.NewUserHandler,
		handler.NewAdminHandler,

		helper.NewHelper,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
