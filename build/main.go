package main

import (
	"os"
	"fmt"

	router "github.com/michaelchandrag/gobiz-spy/module/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}