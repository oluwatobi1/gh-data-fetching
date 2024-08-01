package utils

import "github.com/gin-gonic/gin"

type ResponseType struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func InfoResponse(ctx *gin.Context, message string, data interface{}, statusCode int) {
	ctx.JSON(statusCode, ResponseType{
		Code:    0,
		Data:    data,
		Message: message,
	})
}
