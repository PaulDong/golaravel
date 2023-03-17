package main

import (
	"fmt"
	"time"
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

	return nil
}
