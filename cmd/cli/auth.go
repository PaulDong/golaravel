package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func doAuth() error {
	// migrations
	dbType := gol.DB.DataType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := gol.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := gol.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens cascade;"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	// run migrate
	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}
	// copy file
	err = copyFileFromTemplate("templates/data/user.go.txt", gol.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/data/token.go.txt", gol.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/middleware/user.go.txt", gol.RootPath+"/middleware/user.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/middleware/token.go.txt", gol.RootPath+"/middleware/token.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("  - users, tokens, and remember_bokens migrations created and executed")
	color.Yellow("  - user and token models created")
	color.Yellow("  - auth middleware created")
	color.Yellow("")
	color.Yellow("Don't forget to add user and token models in data/models.go, and to add appropriate middleware to your routes")

	return nil
}
