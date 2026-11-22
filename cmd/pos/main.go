package main

import (
	"log"

	"pos/internal/adapters/console"
	"pos/internal/adapters/csv"
	"pos/internal/adapters/memory"
	"pos/internal/application"
	"pos/internal/domain"
)

func main() {
	course, err := domain.NewCourse(1)
	if err != nil {
		log.Fatal(err)
	}

	repository := memory.NewRepository(course)
	exporter := csv.NewExporter()

	service := application.NewService(repository, exporter)
	menu := console.NewMenu(service)

	if err := menu.Run(); err != nil {
		log.Fatal(err)
	}
}
