package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://footballiot:footballiot@localhost:27018/?maxPoolSize=20&w=majority"

func main() {

	type SensorReadingDoc struct {
		TAGID          int32     `bson:"tag_id"`
		XPOSS          []float64 `bson:"x_positions"`
		YPOSS          []float64 `bson:"y_positions"`
		HEADINGS       []float64 `bson:"headings"`
		DIRECTIONS     []float64 `bson:"directions"`
		ENERGIES       []float64 `bson:"energies"`
		SPEEDS         []float64 `bson:"speeds"`
		TOTALDISTANCES []float64 `bson:"total_distances"`
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	coll := client.Database("footballiot").Collection("players")
	out := SensorReadingDoc{}

	_, err = coll.InsertOne(context.TODO(), out)
	if err != nil {
		panic(err)
	}

}
