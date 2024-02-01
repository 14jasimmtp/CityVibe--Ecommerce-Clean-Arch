package main

import (
	"log"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/docs"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/di"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("/etc/secrets/.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	//	@title						Go + Gin CityViBe
	//	@description				E-Commerce API for clothes
	//	@contact.name				API Support
	//	@in							header
	//	@name						Authorization
	//	@in							header
	//	@BasePath					/
	//	@query.collection.format	multi
	// docs.SwaggerInfo.Version = "1.0"

	docs.SwaggerInfo.Title = "CiTyViBe"
	// // docs.SwaggerInfo.Host = "cityvibe.shop"
	// docs.SwaggerInfo.Host = "localhost:3000"

	di.InitialiseAPI()
}
