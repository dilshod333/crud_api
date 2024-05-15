package main

import (
	"conn/models"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
)



func main(){

	router := gin.Default()
	router.GET("/albums", models.GetAllAlbum)
	router.GET("/albums/:id", models.GetById)
	router.POST("/albums/", models.CreateAlbum)
	router.PUT("/albums/:id", models.UpdateAlbum)
	router.DELETE("/albums/:id", models.DeleteAlbum)
	router.Run(":8080")
}

