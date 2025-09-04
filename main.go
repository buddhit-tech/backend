package main

import (
	"backend/config"
	"backend/routes"
)

func main() {
	config.InitLogger()
	config.LoadEnv()
	// Connect to database
	pool := config.Connect()
	defer pool.Close()
	routes.InitServer(pool)
	select {}
}
