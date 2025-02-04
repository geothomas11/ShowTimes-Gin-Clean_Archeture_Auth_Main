// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"ShowTimes/pkg/api"
	"ShowTimes/pkg/api/handler"
	"ShowTimes/pkg/config"
	"ShowTimes/pkg/db"
	"ShowTimes/pkg/helper"
	"ShowTimes/pkg/repository"
	"ShowTimes/pkg/usecase"
)

// Injectors from wire.go:

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	gormDB, err := db.ConectDatabse(cfg)
	if err != nil {
		return nil, err
	}
	adminRepository := repository.NewAdminRepository(gormDB)
	interfacesHelper := helper.NewHelper(cfg)
	adminUseCase := usecase.NewAdminUseCase(adminRepository, interfacesHelper)
	adminHandler := handler.NewAdminHandler(adminUseCase)
	userRepository := repository.NewUserRepository(gormDB)
	userUseCase := usecase.NewUserUseCase(userRepository, cfg, interfacesHelper)
	userHandler := handler.NewUserHandler(userUseCase)
	otpRepository := repository.NewOtpRepository(gormDB)
	otpUseCase := usecase.NewOtpUsecase(cfg, otpRepository, interfacesHelper)
	otpHandler := handler.NewOtpHandler(otpUseCase)
	categoryRepository := repository.NewCategoryRepository(gormDB)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	productRepository := repository.NewInventoryRepository(gormDB)
	productUseCase := usecase.NewInventoryUseCase(productRepository, interfacesHelper)
	productHandler := handler.NewProductHandler(productUseCase)
	cartRepository := repository.NewCartRepository(gormDB)
	cartUseCase := usecase.NewCartUseCase(cartRepository, productRepository)
	cartHandler := handler.NewCartHandler(cartUseCase)
	orderRepository := repository.NewOrderRepository(gormDB)
	paymentRepository := repository.NewPaymentRepository(gormDB)
	orderUseCase := usecase.NewOrderUseCase(orderRepository, cartRepository, userRepository, paymentRepository)
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepository, orderRepository, cfg)
	orderHandler := handler.NewOrderHandler(orderUseCase, paymentUseCase)
	paymentHandler := handler.NewPaymentHandler(paymentUseCase)
	serverHTTP := http.NewServerHTTP(adminHandler, userHandler, otpHandler, categoryHandler, productHandler, cartHandler, orderHandler, paymentHandler)
	return serverHTTP, nil
}
