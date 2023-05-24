package main

import (
	"log"
	"redispat/router"

	"github.com/gin-gonic/gin"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
		}
	}()

	router.UseRoutes(gin.Default()).Run(":5783")
}
