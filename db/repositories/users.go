package db

import (
	// "database/sql"
	"fmt"
)

type UserRepository interface {
	Create() error
}

type UserRepositoryImpl struct {
	// db *sql.DB
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		// db: _db,
	}
}

func (u *UserRepositoryImpl) Create() error {
	fmt.Println("Creating user in userRepository")
	return nil
}
