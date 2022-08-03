package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "technicaltestdb"
)

type comment struct {
	Id          string `json:"id"`
	TextFR      string `json:"textfr"`
	TextEn      string `json:"texten"`
	PublishedAt string `json:"publishedat"`
	AuthorID    string `json:"authorid"`
	TargetId    string `json:"targetid"`
}

func main() {
	router := gin.Default()
	router.GET("/target/:targetID/comment", getComment)
	router.POST("/target/:targetID/comment", postComment)
	router.Run(":8080")
	return
}

func getComment(c *gin.Context) {
	db := DBConnection()
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM comments WHERE targetid = '%s'", c.Param("targetID")))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var comments []comment

	comments, err = rowsToSlice(rows)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, comments)
}

func rowsToSlice(rows *sql.Rows) ([]comment, error) {
	var comments []comment
	for rows.Next() {
		var cmnt comment
		if err := rows.Scan(&cmnt.Id, &cmnt.TextFR, &cmnt.TextEn, &cmnt.PublishedAt, &cmnt.AuthorID, &cmnt.TargetId); err != nil {
			return comments, err
		}
		comments = append(comments, cmnt)
	}
	if err := rows.Err(); err != nil {
		return comments, err
	}
	return comments, nil
}

func postComment(c *gin.Context) {
	var newComment comment
	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	db := DBConnection()
	_, err := db.Query(fmt.Sprintf("INSERT INTO comments VALUES ('%s', '%s', '%s', '%s', '%s', '%s')", newComment.Id, newComment.TextFR, newComment.TextEn, newComment.PublishedAt, newComment.AuthorID, newComment.TargetId))
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusCreated, newComment)
	defer db.Close()
}

func DBConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}