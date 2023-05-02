package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
)

var appURL string

func doNew(appName string) {
	appName = strings.ToLower(appName)
	appURL = appName

	// sanitize the application name (convert url to single word)
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[(len(exploded) - 1)]
	}

	log.Println("App Name: ", appName)
	// git clone the skeleton application
	color.Green("\tCloning repository ...")
	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "https://github.com/PaulDong/golaravel-app.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		exitGracefully(err)
	}
	// remove .git directory
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a ready to go .env file
	color.Yellow("\tCreating .env file ...")
	currentPath, err := os.Getwd()
	if err != nil {
		exitGracefully(err)
	}

	// Get the parent directory of the current working directory
	parentPath := filepath.Dir(currentPath)
	sourcePath := filepath.Join(parentPath, "/cmd/cli/templates/env.txt")
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		exitGracefully(err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destPath := fmt.Sprintf("./%s/.env", appName)
	destFile, err := os.Create(destPath)
	if err != nil {
		exitGracefully(err)
	}
	defer destFile.Close()

	// Replace the string in the source file and write to the destination file
	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		line := scanner.Text()
		modifiedLine := strings.ReplaceAll(line, "${APP_NAME}", appName)
		modifiedLine = strings.ReplaceAll(modifiedLine, "${KEY}", gol.RandomString(32))
		_, err = destFile.WriteString(modifiedLine + "\n")
		if err != nil {
			exitGracefully(err)
		}
	}
	if err = scanner.Err(); err != nil {
		exitGracefully(err)
	}

	// create a makefile
	if runtime.GOOS == "windows" {
		source, err := os.Open(fmt.Sprintf("/%s/Makefile.windows", appName))
		if err != nil {
			exitGracefully(err)
		}
		defer source.Close()
		destination, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
		if err != nil {
			exitGracefully(err)
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)
		if err != nil {
			exitGracefully(err)
		}
	} else {
		source, err := os.Open(fmt.Sprintf("/%s/Makefile.mac", appName))
		if err != nil {
			exitGracefully(err)
		}
		defer source.Close()
		destination, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
		if err != nil {
			exitGracefully(err)
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)
		if err != nil {
			exitGracefully(err)
		}
	}
	_ = os.Remove("./" + appName + "/Makefile.mac")
	_ = os.Remove("./" + appName + "/Makefile.windows")
	// update the go.mod file
	color.Yellow("\tCreating go.mod file ...")
	_ = os.Remove("./" + appName + "/go.mod")

	data, err := templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		exitGracefully(err)
	}
	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)

	err = copyDataToFile([]byte(mod), "./" + appName + "/go.mod")
	if err != nil {
		exitGracefully(err)
	}

	// update the existing .go files with correct name/imports
	color.Yellow("\tUpdating source files ...")
	os.Chdir("./" + appName)
	updateSource()
	// run go mod tidy in the project dirctory
	color.Yellow("\tRunning go mod tidy ...")
	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Start()
	if err != nil {
		exitGracefully(err)
	}

	color.Green("Done building " + appURL)
	color.Green("Go build something awesome")
}
