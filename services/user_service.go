package services

import (
	db "GoAuthGateway/db/repositories"
	"fmt"
)

type UserService interface {
	CreateUser() error
}

type UserServiceImpl struct {
	userRepository db.UserRepository // Repo interface
}

func NewUserService(_userRepository db.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: _userRepository,
	}
}

func (us *UserServiceImpl) CreateUser() error {
	fmt.Println("Creating user with userService")
	us.userRepository.Create()
	return nil
}
