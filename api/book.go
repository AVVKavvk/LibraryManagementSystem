package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AVVKavvk/LMS/model"
	"github.com/AVVKavvk/LMS/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateBook(ctx echo.Context) error {
	var book model.Book

	if err := ctx.Bind(&book); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}
	if book.IsAllFieldEmpty() {
		return utils.Error(ctx, http.StatusBadRequest, "All field requried")
	}
	if !book.IsVaildCount() {
		return utils.Error(ctx, http.StatusBadRequest, "Book count should be greater than 0")
	}

	filter := bson.M{"bookId": book.BookId}
	err := Book.FindOne(ctx.Request().Context(), filter).Decode(&book)

	if err == nil {
		return utils.Error(ctx, http.StatusBadRequest, "Book already exists")
	} else if err != mongo.ErrNoDocuments {
		return utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to check existing book: %s", err))
	}

	currentTime := primitive.NewDateTimeFromTime(time.Now())
	book.UpdatedAt = currentTime
	book.CreatedAt = currentTime

	_, err = Book.InsertOne(ctx.Request().Context(), book)
	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book not added due to %s", err))
	}

	return utils.Success(ctx, http.StatusOK, "Book added successfully", book.BookId)
}

func UpdateBookCount(ctx echo.Context) error {
	// fmt.Println("viiii")
	type BookCount struct {
		Count int `json:"count"`
	}

	var book model.Book
	var count BookCount
	bookId := ctx.Param("bookId")

	if err := ctx.Bind(&count); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("request failed dues to %s", err.Error()))
	}
	if bookId == "" {
		return utils.Error(ctx, http.StatusBadRequest, "BookId requried")
	}
	if count.Count<= 0 {
		return utils.Error(ctx, http.StatusBadRequest, "Book count should be 1 or more")
	}

	fmt.Println(bookId, count)

	filter := bson.M{"bookId": bookId}

	err := Book.FindOne(ctx.Request().Context(), filter).Decode(&book)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book not found %s", err.Error()))
	}

	totalBook := book.Count + count.Count

	if totalBook < 0 {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Books can't be less than zero: %s", totalBook))
	}

	update := bson.M{"$set": bson.M{"count": totalBook}}

	err = Book.FindOneAndUpdate(ctx.Request().Context(), filter, update).Decode(&book)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book count not updated due to %s", err.Error()))
	}

	return utils.Success(ctx, http.StatusOK, "Book count updated ", echo.Map{"count": totalBook})
}

func GetBookByID(ctx echo.Context) error {
	var book model.Book

	bookId := ctx.Param("bookId")

	if bookId == "" {
		return utils.Error(ctx, http.StatusBadRequest, "BookId is required")
	}

	filter := bson.M{"bookId": bookId}

	err := Book.FindOne(ctx.Request().Context(), filter).Decode(&book)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "Book not found")
	}

	return utils.Success(ctx, http.StatusOK, "Book find", book)
}

func GetAllBooks(ctx echo.Context) error {
	var books []model.Book

	filter := bson.M{}

	cursor, err := Book.Find(ctx.Request().Context(), filter)
	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	defer cursor.Close(ctx.Request().Context())

	for cursor.Next(ctx.Request().Context()) {
		var book model.Book
		if err := cursor.Decode(&book); err != nil {
			return utils.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(ctx, http.StatusOK, "success", books)
}

func GetAllMISAssociateWithBook(ctx echo.Context) error {
	var MISNumbers []string
	var book model.Book
	bookId := ctx.Param("bookId")

	if bookId == "" {
		return utils.Error(ctx, http.StatusBadRequest, "BookId is required")
	}

	filter := bson.M{"bookId": bookId}

	err := Book.FindOne(ctx.Request().Context(), filter).Decode(&book)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book not fount with ID: %s", bookId))
	}

	if len(book.Students) == 0 {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book not assign yet"))
	}

	for _, mis := range book.Students {
		MISNumbers = append(MISNumbers, mis)
	}

	return utils.Success(ctx, http.StatusOK, "Success", MISNumbers)
}

func AssignBookToStudent(ctx echo.Context) error {
	var book model.Book
	var student model.Student

	bookId := ctx.QueryParam("bookId")
	mis := ctx.QueryParam("mis")
	day := ctx.QueryParam("day")

	if bookId == "" || mis == "" || day==""{
		return utils.Error(ctx, http.StatusBadRequest, "BookId ,day and MIS required")
	}

	filter := bson.M{"bookId": bookId}

	err := Book.FindOne(ctx.Request().Context(), filter).Decode(&book)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book not found with error: %s", err.Error()))
	}

	filter = bson.M{"mis": mis}

	err = Student.FindOne(ctx.Request().Context(), filter).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student not found with mis %s due to error %s", mis, err.Error()))
	}

	if book.Count <= 0 {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("No book available with bookId : %s", bookId))
	}

	for _, student := range book.Students {
		if mis == student {
			return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student with mis: %s already have book with ID: %s", mis, bookId))
		}
	}
	newStudnetData := book.Students

	newStudnetData = append(newStudnetData, mis)

	newCount := book.Count
	newCount -= 1

	filter = bson.M{"bookId": bookId}
	update := bson.M{"$set": bson.M{"students": newStudnetData, "count": newCount}}

	err = Book.FindOneAndUpdate(ctx.Request().Context(), filter, update).Decode(&book)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book not assigned to student with MIS %s", mis))
	}

	filter = bson.M{"mis": mis}

	newBooks := student.Books

	newBooks = append(newBooks, bookId)

	update = bson.M{"$set": bson.M{"books": newBooks}}

	err = Student.FindOneAndUpdate(ctx.Request().Context(), filter, update).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintln("Book not assigned to student due to error %s", err.Error()))
	}
	err = AddEntryInIssued(&book,&student,day)
	
	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintln("Book not assigned to student due to error %s", err.Error()))
	}
	return utils.Success(ctx, http.StatusOK, "Book assigned successfully", fmt.Sprintf("Remaining books with book ID %s is %d", bookId, newCount))
}

func DeleteBookFromStudent(ctx echo.Context) error {
	var book model.Book
	var student model.Student

	bookId := ctx.QueryParam("bookId")
	mis := ctx.QueryParam("mis")

	if bookId == "" || mis == "" {
		return utils.Error(ctx, http.StatusBadRequest, "BookId and MIS required")
	}

	bookFilter := bson.M{"bookId": bookId}

	err := Book.FindOne(ctx.Request().Context(), bookFilter).Decode(&book)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book not found with error: %s", err.Error()))
	}

	studentFilter := bson.M{"mis": mis}

	err = Student.FindOne(ctx.Request().Context(), studentFilter).Decode(&student)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student not found with mis %s due to error %s", mis, err.Error()))
	}
	var newBooks []string

	if len(student.Books) == 0 {
		return utils.Error(ctx, http.StatusBadRequest, "Student have no books")
	}
	var cnt = 0
	for _, book := range student.Books {
		if bookId == book {
			cnt += 1
		} else {
			newBooks = append(newBooks, book)
		}
	}

	if cnt == 0 {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student have no book with ID: %s", bookId))
	}

	var newMISNumber []string
	cnt = 0

	for _, studentMIS := range book.Students {
		if studentMIS == mis {
			cnt += 1
		} else {
			newMISNumber = append(newMISNumber, studentMIS)
		}
	}
	if cnt == 0 {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book with ID %s, not assigned to student with MIS: %s", bookId, mis))
	}

	newBookCount := book.Count
	newBookCount += 1

	studentUpdate := bson.M{"$set": bson.M{"books": newBooks}}

	err = Student.FindOneAndUpdate(ctx.Request().Context(), studentFilter, studentUpdate).Decode(&student) //

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Book with ID %s, not deleted form student with MIS", bookId, mis))
	}

	bookUpdate := bson.M{"$set": bson.M{"students": newMISNumber, "count": newBookCount}}

	err = Book.FindOneAndUpdate(ctx.Request().Context(), bookFilter, bookUpdate).Decode(&book)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student with MIS %s, not deleted form book with ID", mis, bookId))
	}
	
	err = RemoveEntryInIssued(bookId,mis)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Student with MIS %s, not deleted form book with ID", mis, bookId))
	}
	return utils.Success(ctx, http.StatusOK, "Book deleted successfully from student", fmt.Sprintf("Remaining books with book ID %s is %d", bookId, newBookCount))
}


func GetBooksByCourse(ctx echo.Context) error {
	var books []model.Book
	course := ctx.Param("course")

	if course == "" {
		return utils.Error(ctx, http.StatusBadRequest, "Course is required")
	}

	filter := bson.M{"course": course}

	cursor, err := Book.Find(ctx.Request().Context(), filter)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Unable to fetch books with course %s", course))
	}

	for cursor.Next(ctx.Request().Context()) {
		var book model.Book
		if err := cursor.Decode(&book); err != nil {
			return utils.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		books = append(books, book)
	}
	if len(books) == 0 {
		return utils.Error(ctx, http.StatusInternalServerError, "No books find for this course")
	}
	return utils.Success(ctx, http.StatusOK, "Successfully", books)
}

func GetBooksBySemWithCourse(ctx echo.Context) error {
	var books []model.Book
	course := ctx.Param("course")
	sem := ctx.Param("sem")

	if course == "" {
		return utils.Error(ctx, http.StatusBadRequest, "Course is required")
	}

	filter := bson.M{"course": course}

	cursor, err := Book.Find(ctx.Request().Context(), filter)

	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Unable to fetch books with course %s", course))
	}

	for cursor.Next(ctx.Request().Context()) {
		var book model.Book
		if err := cursor.Decode(&book); err != nil {
			return utils.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		if book.Sem == sem {
			books = append(books, book)
		}
	}
	if len(books) == 0 {
		return utils.Error(ctx, http.StatusInternalServerError, "No books find for this course and sem")
	}
	return utils.Success(ctx, http.StatusOK, "Successfully", books)
}
