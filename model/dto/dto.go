package dto

import "time"

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type ApiFileRequest struct {
	FileName string   `json:"file_name" validation:"required"`
	Ext      string   `json:"ext"`
	Tags     []string `json:"tags" validation:"required"`
	Bytes    []byte   `json:"data" validation:"required"` // base64
}

type ApiFileResponse struct {
	FileName string   `json:"file_name" validation:"required"`
	Ext      string   `json:"ext"`
	Tags     []string `json:"tags" validation:"required"`
	Bytes    []byte   `json:"data" validation:"required"` // base64
}

type ApiloginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ApiUserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	bytes []byte) ApiFileRequest {
	return ApiFileRequest{
		FileName: fName,
		Tags:     tags,
		Bytes:    bytes,
	}
}

func NewApiFileResponse(fName string, ext string, tags []string,
	bytes []byte) ApiFileResponse {
	return ApiFileResponse{
		FileName: fName,
		Tags:     tags,
		Bytes:    bytes,
	}
}
