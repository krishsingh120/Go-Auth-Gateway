package app

import (
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

func (app *Application) Run() error {

	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      nil,              // TODO: Setup chi router and put it here
		ReadTimeout:  10 * time.Second, // Set read timeout to 10 sec
		WriteTimeout: 10 * time.Second, // Set write timeout to 10 sec
	}

	fmt.Println("Server is listening!")
	fmt.Println("Listening PORT is",app.Config.Addr)

	return server.ListenAndServe()
}
