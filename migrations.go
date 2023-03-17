package golaravel

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (g *Golaravel) MigrateUp(dsn string) error {
	m, err :=  migrate.New("file://" + g.RootPath + "/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running Up migration: ", err)
		return err
	}
	return nil
}

func (g *Golaravel) MigrateDownAll(dsn string) error {
	m, err :=  migrate.New("file://" + g.RootPath + "/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Down(); err != nil {
		log.Println("Error running Down migration: ", err)
		return err
	}
	return nil
}

func (g *Golaravel) Steps(n int, dsn string) error {
	m, err :=  migrate.New("file://" + g.RootPath + "/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Steps(n); err != nil {
		log.Println("Error running Steps migration: ", err)
		return err
	}
	return nil
}

func (g *Golaravel) MigrateForce(dsn string) error {
	m, err :=  migrate.New("file://" + g.RootPath + "/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Force(-1); err != nil {
		log.Println("Error running Force migration: ", err)
		return err
	}
	return nil
}