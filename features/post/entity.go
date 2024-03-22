package post

import (
	"PesbukAPI/features/user"
	"time"

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
	ID 			uint   `json:"id"`
	Picture     string `json:"picture"`
	Content 	string `json:"content"`
	CreatedAt	time.Time `json:"created_at"`
}

type PostWithUser struct {
	ID      uint           `json:"id"`
	UserID  uint           `json:"userid"`
	Picture string         `json:"picture"`
	Content string         `json:"content"`
	UserInfo user.User   `json:"userinfo"` 
}