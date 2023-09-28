package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func GetAllPost(c *gin.Context) {
	db := getDBInstance().db
	data, err := db.Query("SELECT * FROM blog")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message" : "Database error"})
		return
	}
	var posts []Post
	for data.Next() {
		var post Post
		err = data.Scan(&post.Id, &post.Author, &post.Date, &post.Title, &post.Content)
		if err != nil {
			fmt.Println("Error in reading post : ",err)
			continue
		}
		posts = append(posts, post)
	}
	c.JSON(http.StatusAccepted, gin.H{"data": posts})
}
func main() {
	db := getDBInstance().db
	defer db.Close()

	r := gin.Default()

	r.GET("/post", GetAllPost)

	r.Run(":8000")
}