package router

import (
	"github.com/AVVKavvk/LMS/api"
	"github.com/AVVKavvk/LMS/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {

	admin := e.Group("/admin", middleware.IsAdmin())
	admin.POST("", api.CreateAdmin)
	admin.GET("/all", api.GetAllAdmin)
	admin.PUT("", api.UpdateAdminNamePassword)
	admin.PUT("/password", api.UpdateAdminPassword)
	admin.PUT("/forget", api.AdminForgetPassword)
	admin.DELETE("", api.DeleteAdmin)

	admin.GET("/student/:mis", api.GetStudentWithoutPassword)
	admin.POST("/assign/:mis/:bookId", api.AddBookToStudent)
	
	admin.POST("/book", api.CreateBook)
	admin.PUT("/book/:bookId", api.UpdateBookCount)
	admin.PUT("/assign", api.AssignBookToStudent)
	admin.DELETE("/delete", api.DeleteBookFromStudent)
	

	student := e.Group("/student", middleware.IsAuthorized())
	student.POST("", api.CreateStudent)
	student.PUT("/name-phone", api.UpdateStudentNamePhone)
	student.PUT("/password", api.UpdateStudentPassword)
	student.POST("/password", api.GetStudentWithPassword)
	student.GET("/dues-penality/:mis", api.GetStudentPenalityDues)
	student.GET("/books/:mis", api.GetBooksAssociateWithStudent)
	
	book := e.Group("/book", middleware.IsAuthorized())
	book.GET("/:bookId", api.GetBookByID)
	book.GET("/student/:bookId", api.GetAllMISAssociateWithBook)
	book.GET("/course/:course", api.GetBooksByCourse)
	book.GET("/:course/:sem", api.GetBooksBySemWithCourse)
	

	all:=e.Group("/",middleware.IsAuthorized())
	all.POST("login/admin", api.GetAdmin)
	all.GET("admin-profile/:id", api.GetAdminByID)
	all.GET("books", api.GetAllBooks)

	// e.POST("/login/admin", api.GetAdmin)
	// e.GET("/admin-profile/:id",api.GetAdminByID, middleware.IsAdmin())
	// e.GET("/books", api.GetAllBooks,middleware.IsAuthorized())


	e.GET("/", func(ctx echo.Context) error {
		return ctx.File("./index.html")

	})
}
