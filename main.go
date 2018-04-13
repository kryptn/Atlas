package main

import "log"

func main() {

	config := GetConfig()

	app := Application{
		config: config,
		server: Server(config),
	}

	log.Print("Serving")
	app.Serve()
}
