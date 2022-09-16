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
		TIMESTAMP     string  `avro:"timestamp" db:"TIME_STAMP"`
		TAGID         int32   `avro:"tag_id" db:"TAGID"`
		XPOS          float64 `avro:"x_pos" db:"XPOS"`
		YPOS          float64 `avro:"y_pos" db:"YPOS"`
		HEADING       float64 `avro:"heading" db:"HEADING"`
		DIRECTION     float64 `avro:"direction" db:"DIRECTION"`
		ENERGY        float64 `avro:"energy" db:"ENERGY"`
		SPEED         float64 `avro:"speed" db:"SPEED"`
		TOTALDISTANCE float64 `avro:"total_distance" db:"TOTALDISTANCE"`
	}

	/*
		type Page struct {
		    PageId string                 `bson:"pageId" json:"pageId"`
		    Meta   map[string]interface{} `bson:"meta" json:"meta"`
		}
	*/

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

	/*
		import (
			_ "github.com/lib/pq"
			"github.com/jmoiron/sqlx"
			"log"
		)

		type ApplyLeave1 struct {
			LeaveId           int       `db:"leaveid"`
			EmpId             string    `db:"empid"`
			SupervisorEmpId   string    `db:"supervisorid"`
		}

		db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
		if err != nil {
			log.Fatalln(err)
		}

		query := `INSERT INTO TABLENAME(leaveid, empid, supervisorid)
				VALUES(:leaveid, :empid, :supervisorid)`

		var leave1 ApplyLeave1
		_, err := db.NamedExec(query, leave1)
		if err != nil {
			log.Fatalln(err)
		}
	*/

	runtime.Goexit()

}
