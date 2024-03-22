package handler

type PostRequest struct {
	Picture string `json:"picture" form:"picture"`
	Content string `json:"content" form:"content"`
}