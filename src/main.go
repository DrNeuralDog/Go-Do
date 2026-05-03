package main

import (
	"log"

	"godo/src/app"
)

func main() {
	instanceLock, locked := app.CheckSingleInstance()

	if !locked {
		return
	}

	defer instanceLock.Unlock()

	application, err := app.New()

	if err != nil {
		log.Fatal(err)
	}

	if err := application.Initialize(); err != nil {
		log.Fatal(err)
	}

	// Миграция старых TXT-файлов задач в актуальный YAML не должен блокировать запуск!
	if err := application.RunMigration(); err != nil {
		log.Printf("Warning: Migration failed: %v", err)
	}

	application.CreateMainUI()

	application.Run()
}
