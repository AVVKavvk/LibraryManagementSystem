package main

import (
	"net/http"
	"os"

	"github.com/LMS/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	router := echo.New()

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_custom} [INFO] [<NONE>, ${path}]: {"http_method": "${method}", "http_ip": "${remote_ip}", "http_status": ${status}, "http_latency": ${latency}, "http_agent": "${user_agent}", "http_error": "${error}", "http_path": "${path}", "http_bytes_in": ${bytes_in}, "http_bytes_out": ${bytes_out}}` + "\n",
		CustomTimeFormat: "20060102150405.000",
		Output:           os.Stdout,
	}))
	router.Use(middleware.RequestID())

	router.GET("/", func(ctx echo.Context) error {

		// return utils.Success(ctx, http.StatusOK, "Hello World", nil)
		return utils.Error(ctx, http.StatusBadRequest, "Hello World")
	})

	router.Start(":8080")
}
