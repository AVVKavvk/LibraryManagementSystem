package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/AVVKavvk/LMS/model"
	"github.com/AVVKavvk/LMS/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	currentTime := primitive.NewDateTimeFromTime(time.Now())
	admin.CreatedAt = currentTime
	admin.UpdatedAt = currentTime

	filter := bson.M{"email": admin.Email}
	err := Admin.FindOne(ctx.Request().Context(), filter).Decode(&admin)

	if err == nil {
		return utils.Error(ctx, http.StatusConflict, "Admin already exists with this email")
	}

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
	admin.Password=""
	return utils.Success(ctx, http.StatusOK, "Admin login successfully", admin)
}

func UpdateAdminNamePassword(ctx echo.Context) error {
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

	currentTime := primitive.NewDateTimeFromTime(time.Now())
	admin.UpdatedAt = currentTime

	filter := bson.M{"email": admin.Email, "password": admin.Password}
	update := bson.M{"$set": bson.M{"name": admin.Name, "phone": admin.Phone}}

	name := admin.Name
	phone := admin.Phone

	err := Admin.FindOneAndUpdate(ctx.Request().Context(), filter, update).Decode(&admin)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.Error(ctx, http.StatusNotFound, "Admin not found check email or password")
		}
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx, http.StatusOK, "Admin updated successfully", map[string]interface{}{"name": name, "phone": phone})

}

func DeleteAdmin(ctx echo.Context) error {
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

	err := Admin.FindOneAndDelete(ctx.Request().Context(), filter).Decode(&admin)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.Error(ctx, http.StatusNotFound, "Admin not found check email or password")
		}
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx, http.StatusOK, "Admin deleted successfully", admin.Name)
}

func GetAllAdmin(ctx echo.Context) error {
	var admins []model.Admin
	cursor, err := Admin.Find(ctx.Request().Context(), bson.M{})
	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	defer cursor.Close(ctx.Request().Context())

	for cursor.Next(ctx.Request().Context()) {
		var admin model.Admin
		cursor.Decode(&admin)
		admins = append(admins, admin)
	}
	return utils.Success(ctx, http.StatusOK, "All Admins", admins)

}

type updateAdmin struct {
	Admin       model.Admin `json:"admin"`
	NewPassword string      `json:"newPassword"`
}

func UpdateAdminPassword(ctx echo.Context) error {
	var req updateAdmin

	if err := ctx.Bind(&req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if req.Admin.Email == "" || req.Admin.Password == "" || req.NewPassword == "" {
		return utils.Error(ctx, http.StatusBadRequest, "Email, Password, NewPassword are required")
	}

	hash := sha256.New()
	hash.Write([]byte(req.Admin.Password))
	req.Admin.Password = hex.EncodeToString(hash.Sum(nil))

	hash = sha256.New()
	hash.Write([]byte(req.NewPassword))
	newPassword := hex.EncodeToString(hash.Sum(nil))

	currentTime := primitive.NewDateTimeFromTime(time.Now())
	req.Admin.UpdatedAt = currentTime

	filter := bson.M{"email": req.Admin.Email, "password": req.Admin.Password}

	update := bson.M{"$set": bson.M{"password": newPassword, "updatedAt": currentTime}}

	var admin model.Admin

	err := Admin.FindOneAndUpdate(ctx.Request().Context(), filter, update).Decode(&admin)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.Error(ctx, http.StatusNotFound, "Admin not found check email or password")
		}
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx, http.StatusOK, "Admin password updated successfully", admin.Email)
}

func AdminForgetPassword(ctx echo.Context) error {

	// implement later with OTP

	return nil
}

func GetAdminByID(ctx echo.Context) error  {
	var admin model.Admin
	_id:=ctx.Param("id")

	if _id==""{
		return utils.Error(ctx, http.StatusInsufficientStorage, "Please login again")
	}
	filter:=bson.M{"_id":_id}

	err:= Admin.FindOne(ctx.Request().Context(),filter).Decode(&admin)

	if err!=nil{
		if _id==""{
			return utils.Error(ctx, http.StatusInsufficientStorage, "Please login again")
		}
	}
	admin.Password=""

	return utils.Success(ctx, http.StatusOK,"Profile",admin)
}