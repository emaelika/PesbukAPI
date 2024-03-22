package comment

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CommentController interface {
	Add() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	ShowMyComments() echo.HandlerFunc
}

type CommentModel interface {
	AddComment(userid uint, komentarBaru string) (Comment, error)
	UpdateComment(userid uint, commentID uint, data Comment) (Comment, error)
	DeleteComment(userid uint, commentID uint) error
	GetCommentByOwner(userid uint) ([]Comment, error)
}

type CommentService interface {
	AddComment(userid *jwt.Token, komentarBaru string) (Comment, error)
	UpdateComment(userid *jwt.Token, commentID uint, data Comment) (Comment, error)
	DeleteComment(userid *jwt.Token, commentID uint) error
	GetCommentByOwner(userid *jwt.Token) ([]Comment, error)
}

type Comment struct {
	Komentar  string
	ID        uint
	PostId    uint
	Fullname  string
	Avatar    string
	CreatedAt string
}
