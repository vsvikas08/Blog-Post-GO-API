package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql","root:@tcp(localhost:3306)/go_blog")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer db.Close()

	data, err := db.Query("select * from blog")
	if err != nil {
		fmt.Println("Error",err)
		return
	}
	for data.Next() {
		var post Post
		err = data.Scan(&post.Id, &post.Title, &post.Date, &post.Author, &post.Content)
		if err != nil {
			fmt.Println("Error : ",err)
			return
		}
		fmt.Println(post)
	}
}