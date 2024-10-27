package router

import (
	"github.com/AVVKavvk/LMS/api"
	"github.com/AVVKavvk/LMS/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {

	admin := e.Group("/admin", middleware.IsAdmin())
	admin.POST("", api.CreateAdmin)
	admin.GET("", api.GetAdmin)

}
