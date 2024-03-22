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

func ResponseFormatWithUser(code int, message any, data any, userInfo any) map[string]any {
	var result = map[string]any{}
	result["code"] = code
	result["message"] = message
	if data != nil {
		result["data"] = data
	}
	if userInfo != nil {
		result["user_info"] = userInfo
	}

	return result
}
