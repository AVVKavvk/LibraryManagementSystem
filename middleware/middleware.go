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

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			adminEmail := c.Request().Header.Get("X-Admin-Email")
			if adminEmail == "" {
				return utils.Error(c, http.StatusUnauthorized, "You are not authorized")
			}

			var admin model.Admin
			filter := bson.M{"email": adminEmail}

			err := api.Admin.FindOne(c.Request().Context(), filter).Decode(&admin)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return utils.Error(c, http.StatusForbidden, "You are not an admin")
				}
				return utils.Error(c, http.StatusInternalServerError, "Error querying database")
			}

			return next(c)
		}
	}
}
