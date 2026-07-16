package main

import (
	"GoAuthGateway/app"
	config "GoAuthGateway/config/env"
	"fmt"
	"log"
)

func main() {
	config.LoadEnv()
	fmt.Println("Hello World!")

	cfg := app.NewConfig()
	application := app.NewApplication(cfg)

	err := application.Run()
	if err != nil {
		log.Fatalf("Error during server crash: %v", err)
	}
}
