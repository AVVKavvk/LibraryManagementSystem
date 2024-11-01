package middleware

import (
	"net/http"

	"github.com/AVVKavvk/LMS/api"
	"github.com/AVVKavvk/LMS/model"
	"github.com/AVVKavvk/LMS/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsAuthorized() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			adminEmail := c.Request().Header.Get("X-Admin-Email")
			if adminEmail == "" {
				return utils.Error(c, http.StatusUnauthorized, "You are not authorized")
			}

			filter := bson.M{"email": adminEmail}

			var admin model.Admin
			err := api.Admin.FindOne(c.Request().Context(), filter).Decode(&admin)
			if err == nil {
				return next(c)
			} else if err != mongo.ErrNoDocuments {
				return utils.Error(c, http.StatusInternalServerError, "Error querying database for admin")
			}

			var student model.Student
			err = api.Student.FindOne(c.Request().Context(), filter).Decode(&student)
			if err == nil {
				return next(c)
			} else if err == mongo.ErrNoDocuments {
				return utils.Error(c, http.StatusForbidden, "You are not logged in.")
			}

			return utils.Error(c, http.StatusInternalServerError, "Error querying database for student")
		}
	}
}
