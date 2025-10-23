package main

import (
	routes "github.com/dis70rt/subpaper-backend/internal/Routes"
)

func main() {
	router := routes.Setup()
	router.Run(":8080")
}
