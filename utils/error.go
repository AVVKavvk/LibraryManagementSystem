package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type success struct {
	StatusCode int         `json:"statusCode"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

func NewSuccess(msg string, data interface{}) success {
	return success{
		StatusCode: http.StatusOK,
		Msg:        msg,
		Data:       data,
	}
}

func Success(ctx echo.Context, statusCode int, msg string, data interface{}) error {
	response := NewSuccess(msg, data)
	return ctx.JSON(statusCode, response)
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func NewError(msg string, statusCode int) ErrorResponse {
	return ErrorResponse{
		StatusCode: statusCode,
		Msg:        msg,
	}
}

func Error(ctx echo.Context, statusCode int, msg string) error {
	response := NewError(msg, statusCode)
	return ctx.JSON(statusCode, response)
}
