package main

type Post struct {
	Id      string `json:"post_id"`
	Date    string `json:"post_date"`
	Author  string `json:"post_author"`
	Title   string `json:"post_title"`
	Content string `json:"post_content"`
}