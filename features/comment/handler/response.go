package handler

type CommentResponse struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
}