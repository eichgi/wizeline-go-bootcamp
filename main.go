package main

import (
	"academy-go-q12021/routes"
	"github.com/gofiber/fiber"
	"log"
)

func main() {

	app := fiber.New()

	routes.Setup(app)

	//app.Listen(":3000")
	log.Fatal(app.Listen("localhost:3000"), nil)
}
