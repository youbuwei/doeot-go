package main

import (
	"log"

	"github.com/youbuwei/doeot-go/internal/order"
	"github.com/youbuwei/doeot-go/internal/user"
	"github.com/youbuwei/doeot-go/pkg/boot"
)

func main() {
	app := boot.New("json-rpc")

	// Wire the user module with shared DB from app.
	userModule := user.NewModule(app.DB())
	app.RegisterModule(userModule)

	orderModule := order.NewModule(app.DB())
	app.RegisterModule(orderModule)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
