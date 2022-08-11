package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertNewCommentBDD(t *testing.T) {
	var newComment comment
	newComment.Id = "1"
	newComment.TextFR = "Je suis un commentaire!"
	newComment.TextEn = "I am a new comment!"
	newComment.PublishedAt = "23423412"
	newComment.AuthorID = "12345"
	newComment.TargetId = "Photo_12331423"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' as occured when opening stub database connection", err)
	}
	defer db.Close()

	//SqlQuery := fmt.Sprintf("INSERT INTO comments VALUES ('%s', '%s', '%s', '%s', '%s', '%s')", newComment.Id, newComment.TextFR, newComment.TextEn, newComment.PublishedAt, newComment.AuthorID, newComment.TargetId)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO comments").WithArgs(newComment.Id, newComment.TextFR, newComment.TextEn, newComment.PublishedAt, newComment.AuthorID, newComment.TargetId).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	DbInsertNewComment(db, newComment)
	/*
		if err := DbInsertNewComment(newComment); err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}*/
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
