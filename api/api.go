package api

import (
	db "github.com/AVVKavvk/LMS/DB"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Admin   *mongo.Collection
	Student *mongo.Collection
	Book    *mongo.Collection
	Issued  *mongo.Collection
)

func init() {
	Admin = db.MongoClient.Database("LMS").Collection("Admin")
	Student = db.MongoClient.Database("LMS").Collection("Student")
	Book = db.MongoClient.Database("LMS").Collection("Book")
	Issued = db.MongoClient.Database("LMS").Collection("Issued")
}
