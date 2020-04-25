package main

import (
	"os"
	"fmt"

	router "github.com/michaelchandrag/gobiz-spy/module/router"
)

func main() {
	r := router.SetupRouter()
	currentPort := os.Getenv("PORT")
	if currentPort == "" {
		os.Setenv("PORT", "8080")
		os.Setenv("BASE_URL", "http://localhost:8080")
	}
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}