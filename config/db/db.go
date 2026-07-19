package config

import (
	env "GoAuthGateway/config/env"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	// _ "github.com/go-sql-driver/mysql" // req but not use this way to handle in go.
)

func SetUpDB() (*sql.DB, error) {

	cfg := mysql.NewConfig()

	// DB config setup
	cfg.DBName = env.GetString("DB_NAME", "auth_gateway_dev")
	cfg.User = env.GetString("DB_USER", "root")
	cfg.Passwd = env.GetString("DB_PASSWORD", "282005@singh")
	cfg.Addr = env.GetString("DB_ADDR", "127.0.0.1:3306")
	cfg.Net = env.GetString("DB_NET", "tcp")

	fmt.Print("Connecting to Database: ", cfg.DBName, " ", cfg.FormatDSN())

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("Error connecting to DB", err)
		return nil, err
	}

	fmt.Println("Trying to connect DB...")

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println("Error ping to DB", err)
		return nil, pingErr
	}

	fmt.Println("connecting to db successfully🚀")

	return db, nil
}
