package handler

type CommentRequest struct {
	Komentar string `json:"komentar" form:"komentar"`
}