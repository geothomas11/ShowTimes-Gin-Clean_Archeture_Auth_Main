package main

import (
	"ShowTimes/pkg/config"
	"ShowTimes/pkg/di"

	// "ShowTimes/pg/di"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

// @title Go + Gin E-Commerce API Watch Hive
// @version 1.0.0
// @description ShowTimes is an E-commerce platform to purchase Watch
// @contact.name API Support
// @securityDefinitions.apikey BearerTokenAuth
// @in header
// @name Authorization
// @host localhost:7000
// @BasePath /
// @query.collection.format multi
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading the env file")
	}
	config, Err := config.LoadConfig()
	if Err != nil {
		log.Fatal("cannot load config : ", Err)
	}
	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
