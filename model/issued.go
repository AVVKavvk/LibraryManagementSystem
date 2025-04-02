package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Issued struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Course      string             `json:"course,omitempty" bson:"course,omitempty"`
	Sem         string             `json:"sem,omitempty" bson:"sem,omitempty"`
	BookId      string             `json:"bookId,omitempty" bson:"bookId,omitempty"`
	StudentId   string             `json:"studentId,omitempty" bson:"studentId,omitempty"`
	IssueDate   primitive.DateTime `json:"issueDate,omitempty" bson:"issueDate,omitempty"`
	ReturnDate  primitive.DateTime `json:"returnDate,omitempty" bson:"returnDate,omitempty"`
	Penalty    float32            `json:"penalty,omitempty" bson:"penalty,omitempty"`
	BookName    string             `json:"bookName,omitempty" bson:"bookName,omitempty"`
	StudentName string             `json:"studentName,omitempty" bson:"studentName,omitempty"`
}
