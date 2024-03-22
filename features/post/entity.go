package post

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PostController interface {
	Add() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	ShowAllPosts() echo.HandlerFunc
	ShowPostByID() echo.HandlerFunc
}

type PostModel interface {
	AddPost(userid uint, pictureBaru string, contentBaru string) (Post, error)
	UpdatePost(userid uint, postID uint, data Post) (Post, error)
	DeletePost(postID uint) error
	GetAllPosts() ([]Post, error)
	GetPostByID(postID uint) (*Post, error)
}

type PostService interface {
	AddPost(userid *jwt.Token, pictureBaru string, contentBaru string) (Post, error)
	UpdatePost(userid *jwt.Token, postID uint, data Post) (Post, error)
	DeletePost(userid *jwt.Token, postID uint) error
	GetAllPosts() ([]Post, error)
	GetPostByID(postID uint) (*Post, error)
}

type Post struct {
	ID           uint      `json:"id"`
	Fullname     string    `json:"fullname"`
	Avatar       string    `json:"avatar"`
	Picture      string    `json:"picture"`
	Content      string    `json:"content"`
	CreatedAt    string    `json:"created_at"`
	Comments     []Comment `json:"comments"`
	CommentCount int
}

type Comment struct {
	ID        uint   `json:"id"`
	Komentar  string `json:"content"`
	Fullname  string `json:"fullname"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}
