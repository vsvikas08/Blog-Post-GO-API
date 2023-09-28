package main

import (
	"context"
	"encoding/json"
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
		err = data.Scan(&post.Id, &post.Title, &post.Author, &post.Date, &post.Content)
		if err != nil {
			fmt.Println("Error in reading post : ",err)
			continue
		}
		posts = append(posts, post)
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func GetPostByID(c *gin.Context) {
	db := getDBInstance().db
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID can not be empty"})
		return
	}
	var post Post
	err := db.QueryRow("SELECT * FROM blog WHERE post_id = ?",id).Scan(&post.Id,&post.Title,&post.Author,&post.Date,&post.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message" : "Post doesn't exist."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data" : post})
}
func CreateNewPost(c *gin.Context) {
	reqBody, _ := c.GetRawData()
	var post Post
	if err := json.Unmarshal(reqBody,&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message" : "Post data is not valid"})
		return
	}
	if post.Title == "" || post.Author == "" || post.Content == "" || post.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message" : "Field is missing"})
		return
	}
	db := getDBInstance().db
	query := "INSERT INTO `blog` (`post_title`,`post_author`,`post_date`,`post_content`) VALUES (?,?,?,?)"
	insert, err := db.ExecContext(context.Background(), query,post.Title,post.Author,post.Date,post.Content)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message" : "Server error"})
		return
	}
	var data Post
	id,_ := insert.LastInsertId()
	err = db.QueryRow("SELECT * FROM blog where post_id = ?",id).Scan(&data.Id,&data.Title,&data.Author,&data.Date,&data.Content)
	if err != nil {
		fmt.Println("Error in reading data : ",err)
	}
	c.JSON(http.StatusCreated, gin.H{"data" : data})
}
func UpdatePost(c *gin.Context) {
	reqBody, _ := c.GetRawData()
	var data map[string]interface{}
	if err := json.Unmarshal(reqBody,&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message" : "Request data is not correct"})
		return
	}
	if data["post_id"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message" : "ID does not exist to update post"})
		return
	}
	db := getDBInstance().db

	// ID exist or not
	var post Post
	err := db.QueryRow("SELECT * FROM blog where post_id = ?",data["post_id"]).Scan(&post.Id,&post.Title,&post.Author,&post.Date,&post.Content)
	if err != nil {
		fmt.Println("=>> Scan error", err)
	}
	if post.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message" : "ID does not exist"})
		return
	}
	
	// update post data
	query := "UPDATE blog SET "
	flag := true
	for key,val := range data {
		if key == "post_id" {
			continue
		}
		if !flag {
			query += " , "
		}
		query += fmt.Sprintf("%s='%s' ",key,val)
		if flag {
			flag = false
		}
	}
	query += "WHERE post_id = " + post.Id

	_, err = db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message" : "Server error"})
		return
	}
	err = db.QueryRow("SELECT * FROM blog WHERE post_id = ?",post.Id).Scan(&post.Id,&post.Title,&post.Author,&post.Date,&post.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message" : "Data updated but not fetched."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data" : post})
}
func DeletePostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message" : "No ID to delete record."})
		return
	}
	db := getDBInstance().db
	// check data exist or not
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM blog WHERE post_id = %s)",id)
	row := db.QueryRow(query)
	var exists bool
	row.Scan((&exists))
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message" : "ID does not exist"})
		return
	}
	// delete data
	_, err := db.Query("DELETE FROM blog WHERE post_id = ?",id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message" : "Id does not exist."})
		return
	}
	
	msg := fmt.Sprintf("id = %s deleted",id)
	c.JSON(http.StatusOK, gin.H{"message" : msg})
}
func main() {
	db := getDBInstance().db
	defer db.Close()

	r := gin.Default()

	r.GET("/post", GetAllPost)
	r.GET("/post/:id", GetPostByID)
	r.POST("/post", CreateNewPost)
	r.PUT("/update_post", UpdatePost)
	r.DELETE("/post/:id", DeletePostByID)

	r.Run(":8000")
}