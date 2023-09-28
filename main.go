package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql","root:@tcp(localhost:3306)/blog_post")
	if err != nil {
		fmt.Println("Error", err)
	}
	defer db.Close()

	data, err := db.Query("select * from blog")
	if err != nil {
		fmt.Println("Error",err)
	}
	fmt.Println(data)
}