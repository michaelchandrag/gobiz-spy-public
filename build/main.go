package main

import (
	"os"
	"fmt"

	router "github.com/michaelchandrag/go-my-skeleton/module/router"
	database "github.com/michaelchandrag/go-my-skeleton/database"
)

func main() {
	err := database.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	fmt.Println("Database Connected")

	r := router.SetupRouter()
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}