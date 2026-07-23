package db

import (
	"GoAuthGateway/models"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	GetById() (*models.User, error)
	Create() error
	GetAll() ([]*models.User, error)
	DeleteById(id int64) error
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(_db *sql.DB) UserRepository {
	return &UserRepositoryImpl{
		db: _db,
	}
}

func (u *UserRepositoryImpl) GetAll() ([]*models.User, error) {
	query := "SELECT id, username, email FROM users"

	rows, err := u.db.Query(query)
	if err != nil {
		fmt.Println("Error during fetching all row", err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.Id, &user.Username, &user.Email)
		if err != nil {
			fmt.Println("Error during Making arrays of row", err)
			return nil, err
		}

		users = append(users, &user)
	}

	// check for iteration error
	if err = rows.Err(); err != nil {
		fmt.Println("Error during iteration on all row", err)
		return nil, err
	}


	fmt.Println("All Users: ", &users)

	return users, nil
}

func (u *UserRepositoryImpl) DeleteById(_id int64) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := u.db.Exec(query, _id)

	if err != nil {
		fmt.Println("failed to execute delete: %w", err)
		return err
	}

	rowsAffected, rowErr := result.RowsAffected()

	if rowErr != nil {
		fmt.Println("failed to read rows affected: %w", rowErr)
		return rowErr
	}

	if rowsAffected == 0 {
		fmt.Println("Warning: No user was found with that ID. Zero rows deleted.")
		return nil
	}

	fmt.Printf("Successfully deleted %d row(s)\n", rowsAffected)

	return nil
}

func (u *UserRepositoryImpl) Create() error {
	query := "INSERT into users (username, email, password) VALUES (?, ?, ?)"

	result, err := u.db.Exec(query, "testUser", "test123@gamil.com", "password@123")

	if err != nil {
		fmt.Println("Error inserting user", err)
		return err
	}

	rowsAffected, rowErr := result.RowsAffected()

	if rowErr != nil {
		fmt.Println("Error getting row affected", rowErr)
		return rowErr
	}

	if rowsAffected == 0 {
		fmt.Println("No row were affected, user not created")
		return nil
	}

	fmt.Println("User created successfully, rows affected: ", rowsAffected)

	return nil
}

func (u *UserRepositoryImpl) GetById() (*models.User, error) {
	fmt.Println("fetching user in userRepository")

	// Step1: Prepare a query
	query := "SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?"
	fmt.Println(query)

	// Step2: Execute the query
	// QueryRow responsible for this "?" inside query string
	row := u.db.QueryRow(query, 1)

	// Step3: process the result
	user := &models.User{}

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No user found with this given ID")
			return nil, err
		} else {
			fmt.Println("Error scanning user: ", err)
			return nil, err
		}
	}

	// Step4: Print the user details
	fmt.Println("User fetched successfully ", *user)

	return user, nil
}
