package main

import (
	"log"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/docs"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/config"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/di"
)

func main() {
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	//	@title						Go + Gin CityViBe
	//	@description				E-Commerce API for clothes
	//	@contact.name				API Support
	//	@in							header
	//	@name						Authorization
	//	@in							header
	//  @host 						cityvibe.jasim.online
	//	@BasePath					/
	//	@query.collection.format	multi
	// 	docs.SwaggerInfo.Version = "1.0"

	docs.SwaggerInfo.Title = "CiTyViBe"
	docs.SwaggerInfo.Host = "cityvibe.jasim.online"

	di.InitialiseAPI(config)
}
