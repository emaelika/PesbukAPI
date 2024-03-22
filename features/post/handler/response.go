package handler

type PostResponse struct {
	ID        uint      `json:"id"`
	Fullname  string    `json:"fullname"`
	Avatar    string    `json:"avatar"`
	Picture   string    `json:"picture"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
	Comments  []comment `json:"comments"`
}

type comment struct {
	ID        uint   `json:"id"`
	Komentar  string `json:"content"`
	Fullname  string `json:"fullname"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}

type PostsResponse struct {
	ID           string `json:"id"`
	Picture      string `json:"picture"`
	Content      string `json:"content"`
	CreatedAt    string `json:"created_at"`
	CommentCount int    `json:"comment_count"`
}
