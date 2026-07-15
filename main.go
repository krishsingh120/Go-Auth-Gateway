package main

import (
	"GoAuthGateway/app"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")

	cfg := app.NewConfig(":8080")

	app := app.NewApplication(cfg)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
