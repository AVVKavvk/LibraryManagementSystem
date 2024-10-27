package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Issued struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Course      string             `json:"course,omitempty" bson:"course,omitempty"`
	Sem         string             `json:"sem,omitempty" bson:"sem,omitempty"`
	BookId      string             `json:"bookId,omitempty" bson:"bookId,omitempty"`
	StudentId   string             `json:"studentId,omitempty" bson:"studentId,omitempty"`
	IssueDate   string             `json:"issueDate,omitempty" bson:"issueDate,omitempty"`
	ReturnDate  string             `json:"returnDate,omitempty" bson:"returnDate,omitempty"`
	Penality    float32            `json:"penality,omitempty" bson:"penality,omitempty"`
	BookName    string             `json:"bookName,omitempty" bson:"bookName,omitempty"`
	StudentName string             `json:"studentName,omitempty" bson:"studentName,omitempty"`
}
