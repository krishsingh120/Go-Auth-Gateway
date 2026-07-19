package main

import (
	"GoAuthGateway/app"
	dbConfig "GoAuthGateway/config/db"
	envConfig "GoAuthGateway/config/env"
	"fmt"
	"log"
)

func main() {
	envConfig.LoadEnv()
	fmt.Println("Hello World!")

	cfg := app.NewConfig()
	application := app.NewApplication(cfg)
	dbConfig.SetUpDB()

	err := application.Run()
	if err != nil {
		log.Fatalf("Error during server crash: %v", err)
	}
}
