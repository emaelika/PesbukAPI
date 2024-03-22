package handler

type CommentResponse struct {
	ID        uint   `json:"id"`
	Komentar  string `json:"content"`
	Fullname  string `json:"fullname"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}
