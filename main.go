package main

import (
	"backend/config"
	"backend/routes"
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	config.LoadEnv()
	config.InitLogger()
	// Connect to database
	pool := config.Connect()
	defer pool.Close()
	routes.InitServer(pool)
	select {}
}
