package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OwlintTechnicalTestBenjamin/pkg"
	"github.com/OwlintTechnicalTestBenjamin/translate"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// password in clear !!
const (
	host           = "db"
	port           = 5432
	user           = "postgres"
	password       = "docker"
	dbname         = "commentsdb"
	faulty_backend = "https://faulty-backend.herokuapp.com/on_comment"
)

// main function
// launch the router
func main() {
	router := gin.Default()
	router.GET("/target/:targetID/comment", getComment)
	router.POST("/target/:targetID/comment", postComment)
	router.Run(":8080")
	return
}

// function trigger by GET /target/:targetId/comments
// return all comments whos match with the targetid sent in parameters
func getComment(c *gin.Context) {
	db := DBConnection()
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM comments WHERE targetid = '%s'", c.Param("targetID")))
	if err != nil {
		c.String(http.StatusInternalServerError, "Query error")
		return
		//panic(err)
	}
	defer rows.Close()

	var comments []pkg.Comment
	comments, err = rowsToSlice(rows)
	if err != nil {
		c.String(http.StatusInternalServerError, "error during transform *Rows in Slice")
		return
		//panic(err)
	}
	if len(comments) == 0 {
		c.String(http.StatusNotFound, "Comment not found")
		return
	}

	c.IndentedJSON(http.StatusOK, comments)
	return
}

// function who extract data from *sql.Rows to []comment
// []comment can be easely manipulate
func rowsToSlice(rows *sql.Rows) ([]pkg.Comment, error) {
	var comments []pkg.Comment
	for rows.Next() {
		var cmnt pkg.Comment
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

// function trigger by POST /target/:targetId/comments
// get parameters from request to save a new comment in database
// If one or more parameters miss, do we must return a 400 Bad request ?
func postComment(c *gin.Context) {
	db := DBConnection()
	defer db.Close()
	var newComment pkg.Comment
	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	//Traduction
	if newComment.TextFR != "" {
		newComment.TextEn = translate.DeepLTranslate("FR", "EN", newComment.TextFR)
	} else if newComment.TextEn != "" {
		newComment.TextFR = translate.DeepLTranslate("EN", "FR", newComment.TextEn)
	} else {
		return
	}

	DbInsertNewComment(db, newComment)

	//Post sur FaultyBackend
	//extraction du message et de l'auteur et conversion en JSON
	var faulty pkg.Faultymessage
	faulty.Message = newComment.TextFR
	faulty.Author = newComment.AuthorID
	json, err := json.Marshal(faulty)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(faulty_backend, "application/json", bytes.NewReader(json))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//str, err := ioutil.ReadAll(resp.Body)
	//fmt.Printf("response: %s\n", string(str))

	//reponse
	c.IndentedJSON(http.StatusCreated, newComment)
}

// function to connect with the database
// Warning the pointer must be closed outside from this function
func DBConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("sqlOpen error\n")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("db.Ping error")
		panic(err)
	}
	return db
}

func DbInsertNewComment(db *sql.DB, newComment pkg.Comment) *sql.Rows {
	/*
		tx, err := db.Begin()
		rows, err := tx.Exec()
	*/
	rows, err := db.Query(fmt.Sprintf("INSERT INTO comments VALUES ('%s', '%s', '%s', '%s', '%s', '%s')", newComment.Id, newComment.TextFR, newComment.TextEn, newComment.PublishedAt, newComment.AuthorID, newComment.TargetId))
	if err != nil {
		fmt.Printf("DbInsertNewComment: an error as occured\n")
		panic(err)
	}
	return rows
}
