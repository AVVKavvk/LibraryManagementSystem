package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/AVVKavvk/LMS/model"
	"github.com/AVVKavvk/LMS/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAdmin(ctx echo.Context) error {
	var admin model.Admin
	if err := ctx.Bind(&admin); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if admin.Email == "" || admin.Password == "" || admin.Name == "" || admin.Phone == "" {
		return utils.Error(ctx, http.StatusBadRequest, "Email, Password, Name, Phone are required")
	}

	hash := sha256.New()
	hash.Write([]byte(admin.Password))
	admin.Password = hex.EncodeToString(hash.Sum(nil))

	res, err := Admin.InsertOne(ctx.Request().Context(), admin)

	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx, http.StatusCreated, "Admin created successfully", res)
}

func GetAdmin(ctx echo.Context) error {
	var admin model.Admin
	if err := ctx.Bind(&admin); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if admin.Email == "" || admin.Password == "" {
		return utils.Error(ctx, http.StatusBadRequest, "Email, Password are required")
	}

	hash := sha256.New()
	hash.Write([]byte(admin.Password))
	haxPassword := hex.EncodeToString(hash.Sum(nil))

	filter := bson.M{"email": admin.Email, "password": haxPassword}

	err := Admin.FindOne(ctx.Request().Context(), filter).Decode(&admin)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.Error(ctx, http.StatusUnauthorized, "Email not found or wrong password")
		}
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())

	}

	return utils.Success(ctx, http.StatusOK, "Admin found", admin)
}
