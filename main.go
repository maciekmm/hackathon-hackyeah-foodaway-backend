package main

import (
	"log"

	"os"

	"github.com/maciekmm/HackYeah/app"
)

func main() {
	logger := log.New(os.Stdout, "HackYeah!", log.Ldate|log.Lshortfile)
	app := &app.Application{Logger: logger}

	err := app.Init()
	if err != nil {
		logger.Fatal(err)
	}

	err = app.Serve()
	if err != nil {
		logger.Fatal(err)
	}
}
