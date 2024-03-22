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

func ResponseArrayFormat(code int, message any, data any, pagination Pagination) map[string]any {
	var result = map[string]any{}
	result["code"] = code
	result["message"] = message
	if data != nil {
		result["data"] = data
	}
	result["pagination"] = pagination

	return result
}

type Pagination struct {
	page       int
	pageSize   int
	totalPages int
}
