package main

import (
	"dash/database"
	routeController "dash/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database.InitMongoDb()

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/upload", routeController.UploadVideo)

	router.GET("/mpd", routeController.ServeMpd)

	router.GET("/get/:id/:name", routeController.GetMpdById)

	router.Run()
}
