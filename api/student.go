package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/AVVKavvk/LMS/model"
	"github.com/AVVKavvk/LMS/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateStudent(ctx echo.Context) error {
	var student model.Student

	if err := ctx.Bind(&student); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if student.Email == "" || student.Password == "" || student.Name == "" || student.Phone == "" || student.MIS == "" || student.Sem == "" || student.Course == "" {
		return utils.Error(ctx, http.StatusBadRequest, "All fields are required")
	}

	hash := sha256.New()
	hash.Write([]byte(student.Password))
	student.Password = hex.EncodeToString(hash.Sum(nil))

	currentTime := primitive.NewDateTimeFromTime(time.Now())
	student.CreatedAt = currentTime
	student.UpdatedAt = currentTime

	filter := bson.M{"email": student.Email}

	err := Student.FindOne(ctx.Request().Context(), filter).Decode(&student)

	if err == nil {
		return utils.Error(ctx, http.StatusConflict, "Student already exists with this email")
	}

	res, err := Student.InsertOne(ctx.Request().Context(), student)

	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx, http.StatusCreated, "Student created successfully", res)
}

func UpdateStudentNamePhone(ctx echo.Context) error {
	var student model.Student

	if err := ctx.Bind(&student); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if student.MIS == "" {
		return utils.Error(ctx, http.StatusBadRequest, "MIS is required")
	}
	name := student.Name
	phone := student.Phone

	hash := sha256.New()
	hash.Write([]byte(student.Password))
	haxPassword := hex.EncodeToString(hash.Sum(nil))

	filter := bson.M{"mis": student.MIS}

	err := Student.FindOne(ctx.Request().Context(), filter).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusNotFound, "Student not found")
	}

	if student.Password != haxPassword {
		return utils.Error(ctx, http.StatusUnauthorized, "Wrong password")
	}

	update := bson.M{"$set": bson.M{"name": name, "phone": phone, "updatedAt": primitive.NewDateTimeFromTime(time.Now())}}

	err = Student.FindOneAndUpdate(ctx.Request().Context(), filter, update).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx, http.StatusOK, "Student updated successfully", map[string]interface{}{"name": name, "phone": phone})

}

type updateStudent struct {
	Student     model.Student `json:"student"`
	NewPassword string        `json:"newPassword"`
}

func UpdateStudentPassword(ctx echo.Context) error {
	var updateStudent updateStudent

	if err := ctx.Bind(&updateStudent); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if updateStudent.Student.MIS == "" || updateStudent.Student.Password == "" || updateStudent.NewPassword == "" {
		return utils.Error(ctx, http.StatusBadRequest, "MIS, Password, NewPassword are required")
	}

	hash := sha256.New()
	hash.Write([]byte(updateStudent.Student.Password))
	haxPassword := hex.EncodeToString(hash.Sum(nil))

	filter := bson.M{"mis": updateStudent.Student.MIS, "password": haxPassword}

	err := Student.FindOne(ctx.Request().Context(), filter).Decode(&updateStudent.Student)

	if err != nil {
		return utils.Error(ctx, http.StatusNotFound, "Student not found")
	}

	hash = sha256.New()
	hash.Write([]byte(updateStudent.NewPassword))
	newPasswordHex := hex.EncodeToString(hash.Sum(nil))

	if updateStudent.Student.Password == newPasswordHex {
		return utils.Error(ctx, http.StatusBadRequest, "New password is same as old password")
	}

	update := bson.M{"$set": bson.M{"password": newPasswordHex, "updatedAt": primitive.NewDateTimeFromTime(time.Now())}}

	err = Student.FindOneAndUpdate(ctx.Request().Context(), filter, update).Decode(&updateStudent.Student)

	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx, http.StatusOK, "Password updated successfully", nil)

}

func GetStudentWithPassword(ctx echo.Context) error {
	var student model.Student

	if err := ctx.Bind(&student); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if student.MIS == "" || student.Password == "" {
		return utils.Error(ctx, http.StatusBadRequest, "MIS, Password are required")
	}

	hash := sha256.New()
	hash.Write([]byte(student.Password))
	student.Password = hex.EncodeToString(hash.Sum(nil))

	filter := bson.M{"mis": student.MIS, "password": student.Password}

	err := Student.FindOne(ctx.Request().Context(), filter).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student not found with MIS : %s , check your mis or password", student.MIS))
	}

	return utils.Success(ctx, http.StatusOK, "Stduent found", student)

}

func GetStudentWithoutPassword(ctx echo.Context) error {

	var student model.Student

	mis := ctx.Param("mis")

	if mis == "" {
		return utils.Error(ctx, http.StatusBadRequest, "MIS is required")
	}

	filter := bson.M{"mis": mis}

	err := Student.FindOne(ctx.Request().Context(), filter).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student not found with MIS : %s", mis))
	}

	return utils.Success(ctx, http.StatusOK, "Stdeunt found", student)

}

func GetStudentPenalityDues(ctx echo.Context) error {
	var student model.Student
	mis := ctx.Param("mis")

	if mis == "" {
		return utils.Error(ctx, http.StatusBadRequest, "MIS is required")
	}

	filter := bson.M{"mis": mis}

	err := Student.FindOne(ctx.Request().Context(), filter).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student not found with MIS : %s", mis))
	}

	return utils.Success(ctx, http.StatusOK, "success", map[string]interface{}{"penality": student.TotalPenality, "dues": student.Dues})
}

func GetBooksAssociateWithStudent(ctx echo.Context) error {
	var student model.Student
	mis := ctx.Param("mis")

	if mis == "" {
		return utils.Error(ctx, http.StatusBadRequest, "MIS is required")
	}

	filter := bson.M{"mis": mis}

	err := Student.FindOne(ctx.Request().Context(), filter).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student not found with MIS : %s", mis))
	}

	return utils.Success(ctx, http.StatusOK, "Stduent with books", student.Books)
}

func AddBookToStudent(ctx echo.Context) error {

	var student model.Student

	mis := ctx.Param("mis")
	bookId := ctx.Param("bookId")

	if mis == "" || bookId == "" {
		return utils.Error(ctx, http.StatusBadRequest, "MIS and BookId required")
	}

	filter := bson.M{"mis": mis}

	err := Student.FindOne(ctx.Request().Context(), filter).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student not found with MIS : %s", mis))
	}

	// Check if book is avaiable

	for _, book := range student.Books {
		if book == bookId {
			return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student already have this book with ID : %s", bookId))
		}
	}

	UpdateBooks := student.Books
	UpdateBooks = append(UpdateBooks, bookId)

	update := bson.M{"$set": bson.M{"books": UpdateBooks}}

	// remove book count from object

	err = Student.FindOneAndUpdate(ctx.Request().Context(), filter, update).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "Book not added to Student account")
	}
	student.Books = UpdateBooks
	return utils.Success(ctx, http.StatusOK, "Books assign successfully", student)
}
