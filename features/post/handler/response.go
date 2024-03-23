package handler

import "PesbukAPI/features/post"

type PostResponse struct {
	ID        uint   `json:"id"`
	Picture   string `json:"picture"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Avatar    string `json:"avatar"`
	Fullname  string `json:"fullname"`
}

type GetResponse struct {
	ID        uint           `json:"id"`
	Picture   string         `json:"picture"`
	Content   string         `json:"content"`
	CreatedAt string         `json:"created_at"`
	Avatar    string         `json:"avatar"`
	Fullname  string         `json:"fullname"`
	Comments  []post.Comment `json:"comments"`
}

type GetAllResponse struct {
	ID           uint   `json:"id"`
	Picture      string `json:"picture"`
	Content      string `json:"content"`
	CreatedAt    string `json:"created_at"`
	Avatar       string `json:"avatar"`
	Fullname     string `json:"fullname"`
	CommentCount int    `json:"comment_count"`
}
