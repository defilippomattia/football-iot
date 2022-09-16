package main

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"

	"github.com/hamba/avro/v2"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

const uri = "mongodb://footballiot:footballiot@localhost:27017/?maxPoolSize=20&w=majority"

func main() {
	const (
		host     = "localhost"
		port     = 5532
		user     = "footballiot"
		password = "footballiot"
		dbname   = "footballiot"
	)

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
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		fmt.Println(psqlInfo)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			panic(err)
		}

		sqlStatement := `
			INSERT INTO public."RAW_SENSOR_DATA"
			(TIME_STAMP, TAGID, XPOS, YPOS, HEADING, DIRECTION, ENERGY, SPEED, TOTALDISTANCE)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`
		_, err = db.Exec(sqlStatement,
			out.TIMESTAMP,
			out.TAGID,
			out.XPOS,
			out.YPOS,
			out.HEADING,
			out.DIRECTION,
			out.ENERGY,
			out.SPEED,
			out.TOTALDISTANCE)
		if err != nil {
			panic(err)
		}
		//mongo raw db...
		// client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		// if err != nil {
		// 	panic(err)
		// }
		// coll := client.Database("footballiot").Collection("raw-sensor-readings")
		// _, err = coll.InsertOne(context.TODO(), out)
		// if err != nil {
		// 	panic(err)
		// }
	})

	runtime.Goexit()

}
