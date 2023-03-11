package dto

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type ApiFileRequest struct {
	FileName string   `json:"file_name"`
	Ext      string   `json:"ext"`
	Tags     []string `json:"tags"`
	Data     []byte   `json:"data"` // base64
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

func NewApiFileRequest(fName string, ext string, tags []string,
	data []byte) ApiFileRequest {
	return ApiFileRequest{
		FileName: fName,
		Tags:     tags,
		Data:     data,
	}
}
