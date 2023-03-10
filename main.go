package main

import (
	"log"

	"enigmacamp.com/fine_dms/delivery"
)

func main() {
	app := delivery.NewAppServer()

	if err := app.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
