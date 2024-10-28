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
	admin.GET("/all", api.GetAllAdmin)
	admin.PUT("", api.UpdateAdminNamePassword)
	admin.PUT("/password", api.UpdateAdminPassword)
	admin.PUT("/forget", api.AdminForgetPassword)
	admin.DELETE("", api.DeleteAdmin)

}
