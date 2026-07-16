package router

import (
	"GoAuthGateway/controllers"

	"github.com/go-chi/chi"
	// "github.com/go-chi/chi/v5"
)

func SetUpRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/ping", controllers.PingHandler)

	return router
}
