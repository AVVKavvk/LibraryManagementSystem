package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	BookId    string             `json:"bookId,omitempty" bson:"bookId,omitempty"`
	Course    string             `json:"course,omitempty" bson:"course,omitempty"`
	Sem       string             `json:"sem,omitempty" bson:"sem,omitempty"`
	Count     int                `json:"count,omitempty" bson:"count,omitempty"`
	Penalty  float32            `json:"penalty,omitempty" bson:"penalty,omitempty"`
	Students  []string           `json:"students,omitempty" bson:"students,omitempty"`
	CreatedAt primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt primitive.DateTime `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (book *Book) IsAllFieldEmpty() bool {
	if book.BookId == "" || book.Count == 0 || book.Name == "" || book.Course == "" || book.Sem == "" {
		return true
	} else {
		return false
	}
}

func (book *Book) IsVaildCount() bool {
	if book.Count <= 0 {
		return false
	} else {
		return true
	}
}
