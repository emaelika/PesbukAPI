package post

import (
	"PesbukAPI/helper"

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

type PostService interface {
	AddPost(userid *jwt.Token, pictureBaru string, contentBaru string) (Post, error)
	UpdatePost(userid *jwt.Token, postID uint, data Post) (Post, error)
	DeletePost(userid *jwt.Token, postID uint) error
	GetAllPosts(paginasi helper.Pagination) ([]Post, int, error)
	GetPostByID(postID uint) (*Post, error)
}

type PostModel interface {
	AddPost(userid uint, pictureBaru string, contentBaru string) (Post, error)
	UpdatePost(userid uint, postID uint, data Post) (Post, error)
	DeletePost(postID uint) error
	GetAllPosts(paginasi helper.Pagination) ([]Post, int, error)
	GetPostByID(postID uint) (*Post, error)
}

type Post struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Picture      string    `json:"picture"`
	Content      string    `json:"content"`
	CreatedAt    string    `json:"created_at"`
	Avatar       string    `json:"avatar"`
	Fullname     string    `json:"fullname"`
	Comments     []Comment `json:"comments"`
	CommentCount int       `json:"comment_count"`
}

type Comment struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Avatar    string `json:"avatar"`
	Fullname  string `json:"fullname"`
}
