package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doMake(arg2, arg3 string) error {
	switch arg2 {
	case "key":
		rnd := gol.RandomString(32)
		color.Yellow("32 character encryption key: %s", rnd)
	case "migration":
		dbType := gol.DB.DataType
		if arg3 == "" {
			exitGracefully(errors.New("you must give the migration a name"))
		}
		var tableName = arg3
    plur := pluralize.NewClient()
		if plur.IsPlural(arg3) {
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(plur.Plural(arg3))
		}
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), tableName)
		if os.Getenv("TABLE_PREFIX") != "" {
			tableName = os.Getenv("TABLE_PREFIX") + tableName
		}
		upFile := gol.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := gol.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"
		// err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		content, err := templateFS.ReadFile("templates/migrations/migration." + dbType + ".up.sql")
		if err != nil {
			exitGracefully(err)
		}
		mContent := string(content)
		mContent = strings.ReplaceAll(mContent, "$TABLENAME$", tableName)
		err = os.WriteFile(upFile, []byte(mContent), 0644)
		if err != nil {
			exitGracefully(err)
		}
		content, err = templateFS.ReadFile("templates/migrations/migration." + dbType + ".down.sql")
		if err != nil {
			exitGracefully(err)
		}
		mContent = string(content)
		mContent = strings.ReplaceAll(mContent, "$TABLENAME$", tableName)
		err = os.WriteFile(downFile, []byte(mContent), 0644)
		// err = copyFileFromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}
	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}
	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the model a name"))
		}
		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
		}
		model := string(data)
		plur := pluralize.NewClient()
		var modelName = arg3
		var tableName = arg3
		if plur.IsPlural(arg3) {
			modelName = plur.Singular(arg3)
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(plur.Plural(arg3))
		}
		if os.Getenv("TABLE_PREFIX") != "" {
			tableName = os.Getenv("TABLE_PREFIX") + tableName
		}
		fileName := gol.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}
		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)
		err = os.WriteFile(fileName, []byte(model), 0644)
		if err != nil {
			exitGracefully(err)
		}
	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the handler a name"))
		}
		fileName := gol.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}
		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}
		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))
		err = os.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			exitGracefully(err)
		}
	case "mail":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the mail template a name"))
		}
		htmlMail := gol.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
		plainMail := gol.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"

		err := copyFileFromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
		if err != nil {
			exitGracefully(err)
		}
		err = copyFileFromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
		if err != nil {
			exitGracefully(err)
		}
	case "session":
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}
	default:

	}
	return nil
}
