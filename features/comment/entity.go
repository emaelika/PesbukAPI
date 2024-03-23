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
	ShowAllComments() echo.HandlerFunc
}

type CommentModel interface {
	AddComment(userid uint, postID uint, contentBaru string) (Comment, error)
	UpdateComment(userid uint, commentID uint, data Comment) (Comment, error)
	DeleteComment(userid uint, commentID uint) error
	GetCommentByOwner(userid uint) ([]Comment, error)
	GetAllComments() ([]Comment, error)
}

type CommentService interface {
	AddComment(userid *jwt.Token, postID uint, contentBaru string) (Comment, error)
	UpdateComment(userid *jwt.Token, commentID uint, data Comment) (Comment, error)
	DeleteComment(userid *jwt.Token, commentID uint) error
	GetCommentByOwner(userid *jwt.Token) ([]Comment, error)
	GetAllComments() ([]Comment, error)
}

type Comment struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Avatar    string `json:"avatar"`
	Fullname  string `json:"fullname"`
}
