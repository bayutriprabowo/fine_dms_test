package dto

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewApiResponseSuccess(msg string, data any) ApiResponse {
	return ApiResponse{
		Status:  "Success",
		Message: msg,
		Data:    data,
	}
}

func NewApiResponseFailed(msg string) ApiResponse {
	return ApiResponse{
		Status:  "Failed",
		Message: msg,
	}
}
