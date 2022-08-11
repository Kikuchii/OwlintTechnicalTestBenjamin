package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

// password in clear !!
const (
	host           = "0.0.0.0"
	port           = 5432
	user           = "postgres"
	password       = "docker"
	dbname         = "commentsdb"
	faulty_backend = "https://faulty-backend.herokuapp.com/on_comment"
	deepLKey       = "5e4e98a7-332f-9230-d405-3f6f28f7ebaf:fx"
)

type comment struct {
	Id          string `json:"id"`
	TextFR      string `json:"textfr"`
	TextEn      string `json:"texten"`
	PublishedAt string `json:"publishedat"`
	AuthorID    string `json:"authorid"`
	TargetId    string `json:"targetid"`
}

type deepLResponse struct {
	Translations []translatedComment `json:"translations"`
}

type translatedComment struct {
	Detected_source_language string `json:"detected_source_language"`
	Text                     string `json:"text"`
}

type faultymessage struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}

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
		panic(err)
	}
	defer rows.Close()

	var comments []comment
	comments, err = rowsToSlice(rows)
	if err != nil {
		panic(err)
	}
	if len(comments) == 0 {
		c.String(http.StatusNotFound, "Comment not found")
	}

	c.IndentedJSON(http.StatusOK, comments)
}

// function who extract data from *sql.Rows to []comment
// []comment can be easely manipulate
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

// function trigger by POST /target/:targetId/comments
// get parameters from request to save a new comment in database
// If one or more parameters miss, do we must return a 400 Bad request ?
func postComment(c *gin.Context) {
	db := DBConnection()
	defer db.Close()
	var newComment comment
	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	//Translate
	if newComment.TextFR != "" {
		newComment.TextEn = Translate("FR", "EN", newComment.TextFR)
	} else if newComment.TextEn != "" {
		newComment.TextFR = Translate("EN", "FR", newComment.TextEn)
	} else {
		return
	}

	DbInsertNewComment(db, newComment)

	//Post sur FaultyBackend
	//extraction du message et de l'auteur et conversion en JSON
	var faulty faultymessage
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
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func Translate(source_lang string, target_lang string, comment string) string {
	urlstr := "https://api-free.deepl.com/v2/translate?auth_key=" + deepLKey + "&text=" + url.QueryEscape(comment) + "&target_lang=" + target_lang + "&source_lang=" + source_lang

	resp, err := http.Post(urlstr, "application/x-www-form-urlencoded", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	str, err := ioutil.ReadAll(resp.Body)

	var DeepLResponse deepLResponse
	json.Unmarshal(str, &DeepLResponse)

	return DeepLResponse.Translations[0].Text
}

func DbInsertNewComment(db *sql.DB, newComment comment) *sql.Rows {
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
