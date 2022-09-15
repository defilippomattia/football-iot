package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/hamba/avro/v2"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://footballiot:footballiot@localhost:27017/?maxPoolSize=20&w=majority"

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
		coll := client.Database("footballiot").Collection("raw-sensor-readings")
		_, err = coll.InsertOne(context.TODO(), out)
		if err != nil {
			panic(err)
		}
	})

	runtime.Goexit()

}
