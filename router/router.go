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

	admin.GET("/student/:mis", api.GetStudentWithoutPassword)
	admin.POST("/assign/:mis/:bookId", api.AddBookToStudent)

	student := e.Group("/student")
	student.POST("", api.CreateStudent)
	student.PUT("/name-phone", api.UpdateStudentNamePhone)
	student.PUT("/password", api.UpdateStudentPassword)
	student.GET("/password", api.GetStudentWithPassword)
	student.GET("/dues-penality/:mis", api.GetStudentPenalityDues)
	student.GET("/books/:mis", api.GetBooksAssociateWithStudent)

}
