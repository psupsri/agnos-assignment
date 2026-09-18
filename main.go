package main

import (
	"agnos-assignment/config"
	"agnos-assignment/routes"
)

func main() {
	db := config.ConnectDatabaseAndMigrate()

	r := routes.SetUpRouter(db)

	r.Run(":8080")
}
