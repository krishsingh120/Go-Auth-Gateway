package services

import (
	db "GoAuthGateway/db/repositories"
	"fmt"
)

type UserService interface {
	GetUserById() error
}

type UserServiceImpl struct {
	userRepository db.UserRepository // Repo interface
}

func NewUserService(_userRepository db.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: _userRepository,
	}
}

func (us *UserServiceImpl) GetUserById() error {
	fmt.Println("fetching user with userService")
	us.userRepository.GetAll()
	return nil
}
