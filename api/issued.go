package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AVVKavvk/LMS/model"
	"github.com/AVVKavvk/LMS/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddEntryInIssued(book *model.Book, student *model.Student, day string) error {
	var issued model.Issued
	issued.Sem=book.Sem
	issued.Course=book.Course
	issued.BookId=book.BookId
	issued.BookName= book.Name
	issued.StudentId=student.MIS
	issued.StudentName=student.Name
	
	issued.IssueDate= primitive.NewDateTimeFromTime(time.Now())

	d, err := strconv.Atoi(day)
	if err != nil {
		return errors.New(fmt.Sprintf("Can't convert string to int for day %v",d))
	}
	issued.ReturnDate= primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, d))

	_, err = Issued.InsertOne(context.Background(),issued)

	if err!=nil{
		return errors.New(fmt.Sprintf("Book not added due to %s", err))
	}
	return nil
}

func RemoveEntryInIssued(bookId string, mis string) error {
	filter := bson.M{"studentId": mis, "bookId": bookId}

	result := Issued.FindOneAndDelete(context.Background(), filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("no matching record found to delete")
		}
		return fmt.Errorf("failed to delete book entry: %v", err)
	}
	return nil
}

func GetIssuedBooks(ctx echo.Context) error {

	cur, err := Issued.Find(context.Background(), bson.M{})
	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	var issued []model.Issued

	defer cur.Close(ctx.Request().Context())

	for cur.Next(context.Background()) {
		var item model.Issued
		err := cur.Decode(&item)
		if err != nil {
			return utils.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		issued = append(issued, item)
		}
	return utils.Success(ctx,200,"All Assigned Books",issued)
}

func FindAssignBookByMIS (ctx echo.Context) error {
	var issued []model.Issued

	mis:= ctx.Param("mis")

	if mis == "" || len(mis)<9 {
		return utils.Error(ctx, http.StatusBadRequest, "Invalid MIS")
	}

	filter := bson.M{"studentId":mis}

	cur, err := Issued.Find(ctx.Request().Context(),filter);

	if err !=nil{
		return utils.Error(ctx, http.StatusUnavailableForLegalReasons,fmt.Sprintf("No books find with MIS %s",mis))
	}
	defer cur.Close(ctx.Request().Context())

	for cur.Next(context.Background()){
		var item model.Issued
		err := cur.Decode(&item)
		if err != nil {
			return utils.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		issued = append(issued, item)
	}
	return utils.Success(ctx,200,"All Assigned Books",issued)
}

func FindAssignBookByBookId (ctx echo.Context) error{
	var issued []model.Issued
	bookId := ctx.Param("bookId")

	if bookId == "" || len(bookId)<=2{
		return utils.Error(ctx, http.StatusBadRequest, "Invalid BookId")
	}

	filter := bson.M{"bookId":bookId}

	cur, err := Issued.Find(ctx.Request().Context(), filter)

	if err != nil{
		return utils.Error(ctx, http.StatusUnavailableForLegalReasons, fmt.Sprintf("No books find with bookId %s",bookId))
	}
	defer cur.Close(ctx.Request().Context())
	for cur.Next(context.Background()){
		var item model.Issued
		err := cur.Decode(&item)
		if err != nil {
			return utils.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		issued = append(issued, item)
	}
	return utils.Success(ctx,200,"All Assigned Books",issued)
}

func GetDueBooks(ctx echo.Context) error {

	days:= ctx.Param("days")
	if days == "" || len(days)<1 {
    return utils.Error(ctx, http.StatusBadRequest, "Invalid days")
  }
	
	d, err := strconv.Atoi(days)
	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "Invalid days")
  }

	currentTime := time.Now()
	thresholdTime := currentTime.Add(time.Duration(d) * 24 * time.Hour)

	filter := bson.M{"returnDate": bson.M{"$lt": primitive.NewDateTimeFromTime(thresholdTime)}}
	
	cursor, err := Issued.Find(context.TODO(), filter)
	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	defer cursor.Close(context.TODO())

	var dueBooks []model.Issued
	if err = cursor.All(context.TODO(), &dueBooks); err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx,200,"Due Books",dueBooks)
}

func GetExpiryBooks(ctx echo.Context) error {
	currentTime := time.Now()
	filter := bson.M{"returnDate": bson.M{"$lt": primitive.NewDateTimeFromTime(currentTime)}}
	
	cursor, err := Issued.Find(context.TODO(), filter)
	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	defer cursor.Close(context.TODO())

	var dueBooks []model.Issued
	if err = cursor.All(context.TODO(), &dueBooks); err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx,200,"Due Books",dueBooks)
}