package app

import (
	config "GoAuthGateway/config/env"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Addr string
}

type Application struct {
	Config Config
}

// constructor for creating newConfig
func NewConfig() Config {
	port := config.GetString("PORT", ":8080")

	return Config{
		Addr: port,
	}
}

// constructor for creating newApplication
func NewApplication(cfg Config) *Application {
	return &Application{
		Config: cfg,
	}
}

func (app *Application) Run() error {

	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      nil,              // TODO: Setup chi router and put it here
		ReadTimeout:  10 * time.Second, // Set read timeout to 10 sec
		WriteTimeout: 10 * time.Second, // Set write timeout to 10 sec
	}

	fmt.Println("🚀 Server is listening!")
	fmt.Println("Listening PORT is", app.Config.Addr)

	return server.ListenAndServe()
}
