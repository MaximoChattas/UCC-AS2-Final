package db

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDb *mongo.Database
var client *mongo.Client

var HotelsCollection *mongo.Collection
var AmenitiesCollection *mongo.Collection

func Disconect_db() {

	client.Disconnect(context.TODO())
}

func Init_db() {

	clientOpts := options.Client().ApplyURI("mongodb://root:pass@mongodatabase:27017/?authSource=admin&authMechanism=SCRAM-SHA-256")
	cli, err := mongo.Connect(context.TODO(), clientOpts)
	client = cli

	if err != nil {
		log.Info("Connection Failed to Open")
		log.Fatal(err)
	} else {
		log.Info("Connection Established")
	}

	dbNames, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		log.Info("Failed to get databases available")
		log.Fatal(err)
	}

	mongoDb = client.Database("test")

	fmt.Println("Available datatabases:")
	fmt.Println(dbNames)

	HotelsCollection = mongoDb.Collection("hotels")
	AmenitiesCollection = mongoDb.Collection("amenities")

}
