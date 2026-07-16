package app

import (
	config "GoAuthGateway/config/env"
	"GoAuthGateway/controllers"
	db "GoAuthGateway/db/repositories"
	"GoAuthGateway/router"
	"GoAuthGateway/services"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Addr string
}

type Application struct {
	Config Config
	Store  db.Storage
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
		Store:  *db.NewStorage(),
	}
}

func (app *Application) Run() error {

	ur := db.NewUserRepository()
	us := services.NewUserService(ur)
	uc := controllers.NewUserController(us)
	uRouter := router.NewUserRouter(uc)

	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      router.SetUpRouter(uRouter),
		ReadTimeout:  10 * time.Second, // Set read timeout to 10 sec
		WriteTimeout: 10 * time.Second, // Set write timeout to 10 sec
	}

	fmt.Println("🚀 Server is listening!")
	fmt.Println("Listening PORT is", app.Config.Addr)

	return server.ListenAndServe()
}
