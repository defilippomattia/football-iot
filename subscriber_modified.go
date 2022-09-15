package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"runtime"

	"github.com/hamba/avro/v2"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://footballiot:footballiot@localhost:27018/?maxPoolSize=20&w=majority"

func main() {

	type SensorReading struct {
		TIMESTAMP     string  `avro:"timestamp"`
		TAGID         int32   `avro:"tag_id"`
		XPOS          float64 `avro:"x_pos"`
		YPOS          float64 `avro:"y_pos"`
		HEADING       float64 `avro:"heading"`
		DIRECTION     float64 `avro:"direction"`
		ENERGY        float64 `avro:"energy"`
		SPEED         float64 `avro:"speed"`
		TOTALDISTANCE float64 `avro:"total_distance"`
	}

	type SensorReadingDoc struct {
		TAGID          int32     `bson:"tag_id"`
		XPOSS          []float64 `bson:"x_positions,omitempty"`
		YPOSS          []float64 `bson:"y_positions,omitempty"`
		HEADINGS       []float64 `bson:"headings,omitempty"`
		DIRECTIONS     []float64 `bson:"directions,omitempty"`
		ENERGIES       []float64 `bson:"energies,omitempty"`
		SPEEDS         []float64 `bson:"speeds,omitempty"`
		TOTALDISTANCES []float64 `bson:"total_distances,omitempty"`
	}
	// Id   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	schema, err := avro.Parse(`{
		"type": "record",
		"name": "reading",
		"namespace": "org.hamba.avro",
		"fields" : [
			{"name": "timestamp", "type": "string"},
			{"name": "tag_id", "type": "int"},
			{"name": "x_pos", "type": "double"},
			{"name": "y_pos", "type": "double"},
			{"name": "heading", "type": "double"},
			{"name": "direction", "type": "double"},
			{"name": "energy", "type": "double"},
			{"name": "speed", "type": "double"},
			{"name": "total_distance", "type": "double"}
		]
	}`)

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}
	nc.Subscribe("football-iot.players.*", func(m *nats.Msg) {
		out := SensorReading{}
		err = avro.Unmarshal(schema, m.Data, &out)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out)

		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}

		//
		var result bson.M
		coll := client.Database("footballiot").Collection("players")

		err = coll.FindOne(context.TODO(), bson.D{{"tag_id", out.TAGID}}).Decode(&result)
		if err != nil {
			//document doesn't exists with tagid...
			out := SensorReadingDoc{
				TAGID: out.TAGID,
			}
			_, err = coll.InsertOne(context.TODO(), out)
			if err != nil {
				panic(err)
			}
		} else {
			//document exists, append to array
			fmt.Println("doc exists")
			result["directions"] = out.DIRECTION
			fmt.Println(result["directions"])

			filter := bson.D{{"tag_id", out.TAGID}}
			update := bson.D{
				{
					"$push", bson.D{
						{"x_positions", math.Round(out.XPOS*100) / 100},
						{"y_positions", math.Round(out.YPOS*100) / 100},
						{"headings", math.Round(out.HEADING*100) / 100},
						{"directions", math.Round(out.DIRECTION*100) / 100},
						{"energies", math.Round(out.ENERGY*100) / 100},
						{"speeds", math.Round(out.SPEED*100) / 100},
						{"total_distances", math.Round(out.TOTALDISTANCE*100) / 100},
					},
				},
			}
			result, err := coll.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				panic(err)
			}
			print(result)

		}

	})

	runtime.Goexit()

}
