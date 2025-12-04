package main

import (
	"log"

	"github.com/youbuwei/doeot-go/internal/modules"
	"github.com/youbuwei/doeot-go/pkg/boot"
)

func main() {
	app := boot.New("json-rpc")

	for _, m := range modules.All(app.DB()) {
		app.RegisterModule(m)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
