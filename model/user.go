package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	Email         string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone         string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Password      string             `json:"password,omitempty" bson:"password,omitempty"`
	MIS           string             `json:"mis,omitempty" bson:"mis,omitempty"`
	Sem           string             `json:"sem,omitempty" bson:"sem,omitempty"`
	Course        string             `json:"course,omitempty" bson:"course,omitempty"`
	TotalPenality float32            `json:"totalPenality,omitempty" bson:"totalPenality,omitempty"`
	Books         []string           `json:"books,omitempty" bson:"books,omitempty"`
}
