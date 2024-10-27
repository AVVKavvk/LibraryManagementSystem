package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	BookId   string             `json:"bookId,omitempty" bson:"bookId,omitempty"`
	Course   string             `json:"course,omitempty" bson:"course,omitempty"`
	Sem      string             `json:"sem,omitempty" bson:"sem,omitempty"`
	Count    int                `json:"count,omitempty" bson:"count,omitempty"`
	Penality float32            `json:"penality,omitempty" bson:"penality,omitempty"`
	Student  []string           `json:"student,omitempty" bson:"student,omitempty"`
}
