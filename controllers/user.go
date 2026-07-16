package controllers

import (
	"GoAuthGateway/services"
	"fmt"
	"net/http"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(_userService services.UserService) *UserController {
	return &UserController{
		UserService: _userService,
	}
}

func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating user in userController")
	uc.UserService.CreateUser()
	w.Write([]byte("Register user controller"))
}
