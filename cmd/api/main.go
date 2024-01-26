package main

import (
	"log"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/di"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("/home/jasim/CityViBe-Project-Ecommerce/.env")
	// err := godotenv.Load("/home/jasim/CityViBe-Project-Ecommerce/.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}


	di.InitialiseAPI()
}
