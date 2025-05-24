package main

import (
	"API_VentasGO/api"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.InitRoutes(r)

	if err := r.Run(":9090"); err != nil {
		panic(fmt.Errorf("error trying to start server: %v", err))
	}

}
