package helper

func ResponseFormat(code int, message any, data any) map[string]any {
	var result = map[string]any{}
	result["code"] = code
	result["message"] = message
	if data != nil {
		result["data"] = data
	}

	return result
}

type Pagination struct {
	Page       int `json:"page"`
	Pagesize   int `json:"pagesize"`
	TotalPages int `json:"total_pages"`
}

func ResponseArrayFormat(code int, message any, data any, paginasi Pagination) map[string]any {
	var result = map[string]any{}
	result["code"] = code
	result["message"] = message
	if data != nil {
		result["data"] = data
	}
	result["pagination"] = paginasi

	return result
}
