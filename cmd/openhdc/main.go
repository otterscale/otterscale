package main

import (
	_ "github.com/openhdc/openhdc/internal/migrations"
)

var version = "devel"

func main() {
	// wire app
	app, cleanup, err := wireApp(version)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Start(); err != nil {
		panic(err)
	}
}
