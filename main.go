package main

import (
	"github.com/Alitindrawan24/go-mrt-schedules/modules/station"
	"github.com/gin-gonic/gin"
)

func main() {
	Init()
}

func Init() {
	router := gin.Default()
	api := router.Group("/api/v1")

	station.Init(api)

	router.Run(":8080")
}
