package main

var version = "devel"

func main() {
	// wire app
	app, cleanup, err := wireApp(version)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Execute(); err != nil {
		panic(err)
	}
}
