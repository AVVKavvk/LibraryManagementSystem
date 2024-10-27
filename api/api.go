package api

import (
	db "github.com/AVVKavvk/LMS/DB"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Admin  *mongo.Collection
	User   *mongo.Collection
	Book   *mongo.Collection
	Issued *mongo.Collection
)

func init() {
	Admin = db.MongoClient.Database("LMS").Collection("Admin")
	User = db.MongoClient.Database("LMS").Collection("User")
	Book = db.MongoClient.Database("LMS").Collection("Book")
	Issued = db.MongoClient.Database("LMS").Collection("Issued")
}
