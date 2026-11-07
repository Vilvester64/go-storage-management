package utils

import "github.com/gin-gonic/gin"

// APIResponse adalah format response JSON standar sesuai docs/api_reference.md.
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse mengirimkan response JSON berstatus "success".
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirimkan response JSON berstatus "error".
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Status:  "error",
		Message: message,
	})
}
