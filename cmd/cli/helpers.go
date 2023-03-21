package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup() {
	err := godotenv.Load()
	if err != nil {
		exitGracefully(err)
	}

	path, err := os.Getwd()
	if err != nil {
		exitGracefully(err)
	}

	gol.RootPath = path
	gol.DB.DataType = os.Getenv("DATABASE_TYPE")
}

func getDSN() string {
	dbType := gol.DB.DataType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
		return dsn
	}
	return "mysql://" + gol.BuildDSN()
}

func showHelp() {
	color.Yellow(`Available commands:

  help                  - show the help commands
  version               - print application version
  migrate               - runs all up migrations that have not been run previously
  migrate up            - runs all up migrations that have not been run previously
  migrate down          - reverses the most recent migration
  migrate reset         - runs all down migration in reverse order, and then all up migrations
  migrate force         - runs all down migration 1 
  make migration <name> - create two new up and down migration in the migration folder
  make auth             - creates and runs migrations for authentication tables, and creates models and middleware
  make handler <name>   - creates a stub handler in the handlers directory
  make model <name>     - creates a new model in the models directory
  `)
}
