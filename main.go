package main

import (
	"GoAuthGateway/app"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")

	cfg := app.Config{
		Addr: ":5000",
	}

	app := app.Application{
		Config: cfg,
	}

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
