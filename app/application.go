package app

import (
	dbConfig "GoAuthGateway/config/db"
	config "GoAuthGateway/config/env"
	"GoAuthGateway/controllers"
	repo "GoAuthGateway/db/repositories"
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
	Store  repo.Storage
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
		Store:  *repo.NewStorage(),
	}
}

func (app *Application) Run() error {

	db, err := dbConfig.SetUpDB()

	if err != nil {
		fmt.Println("Error Setting up databases", err)
		return err
	}

	ur := repo.NewUserRepository(db)
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
