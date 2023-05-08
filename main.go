package main

import (
	"github.com/DenrianWeiss/anvilEstimate/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/api/v1/simulation/run", handler.HandleSimulationRequest)
	r.GET("/api/v1/simulation/:entry", handler.GetSimulationResult)

	err := r.Run()
	if err != nil {
		panic(err)
	}
	select {}
}
