package models

import (
	// "conn/models"
	// "conn/models"
	"conn/models/postgres"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func InsertData() {
	db := postgres.Connection()
	smtm, err := db.Prepare("insert into albums(title, artist, price) values($1, $2, $3)")
	if err != nil {
		log.Fatal("error while preparing..", err)
	}

	for _, alb := range albums {
		_, err = smtm.Exec(alb.Title, alb.Artist, alb.Price)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Print("Successfully inserted data...")
	smtm.Close()

}

func GetAllAlbum(c *gin.Context) {
	db := postgres.Connection()

	rows, err := db.Query("SELECT * FROM albums")
	if err != nil {
		log.Fatal("Error while getting info from albums", err)
	}
	defer rows.Close()

	var allAlbums []album
	for rows.Next() {
		var al album
		if err := rows.Scan(&al.ID, &al.Title, &al.Artist, &al.Price); err != nil {
			log.Fatal("Error while scanning from database", err)
		}
		allAlbums = append(allAlbums, al)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if len(allAlbums) == 0 {
		c.IndentedJSON(http.StatusOK, []album{})
		return
	}
	c.IndentedJSON(http.StatusOK, allAlbums)
}

func GetById(c *gin.Context) {
	var check bool
	id := c.Param("id")

	db := postgres.Connection()

	rows, err := db.Query("select * from albums where id=$1", id)
	if err != nil {
		log.Fatal("error while getting from table...", err)
	}
	var justOne []album
	for rows.Next() {
		var a album
		if err = rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
			log.Fatal(err)
		}
		if id == a.ID {
			justOne = append(justOne, a)
			check = true 
			break
		}
	}
	if check{
		c.IndentedJSON(http.StatusOK, justOne)

	}
	if !check{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	}

	defer db.Close()
}

func CreateAlbum(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bro you did not create anything...."})
		log.Fatal("while binding json to struct error...", err)
	}

	db := postgres.Connection()
	_, err := db.Exec("insert into albums(title, artist, price) values($1, $2, $3)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		log.Fatal("wrong while insertng data to the table...", err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully created bro..."})

}

func UpdateAlbum(c *gin.Context) {
	id := c.Param("id")

	var change album
	if err := c.BindJSON(&change); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Something went wrong with JSON binding"})
		return
	}

	db := postgres.Connection()

	var count int
	err := db.QueryRow("select count(*) from albums where id=$1", id).Scan(&count)
	if err != nil {
		log.Fatal("smth wrong...:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Album not found"})
		return
	}

	_, err = db.Exec("UPDATE albums SET title=$1, artist=$2, price=$3 WHERE id=$4", change.Title, change.Artist, change.Price, id)
	if err != nil {
		log.Fatal("errrror updating album:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album updated successfully"})
}

func DeleteAlbum(c *gin.Context) {
	var check = false
	id := c.Param("id")
	var change album
	if err := c.BindJSON(&change); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "smth wrong with binding..."})
	}
	db := postgres.Connection()
	rows, err := db.Query("select * from albums")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var a album
		if err = rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "smth wrong with deleting..."})
			return
		}
		if a.ID == id {
			_, err := db.Exec("delete from albums where id=$1", id)
			if err != nil {
				log.Fatal("smth wrong with deleting brother...", err)
			}
			c.JSON(http.StatusOK, gin.H{"message": "deleted brooo..."})
			check = true
			return
		}

	}
	if !check {
		c.JSON(http.StatusBadRequest, gin.H{"message": "just double check which id you want to delete..."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "smth wrong with deleting..."})
		return
	}

}
