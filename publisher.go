package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hamba/avro/v2"
	"github.com/nats-io/nats.go"
)

func readPlayerCsv(csvPath string) {

	nc, _ := nats.Connect(nats.DefaultURL)

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
	if err != nil {
		fmt.Println(err)
	}
	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines {
		timestampL := line[0]
		tagIdL, _ := strconv.Atoi(line[1])
		xPosL, _ := strconv.ParseFloat(line[2], 64)
		yPosL, _ := strconv.ParseFloat(line[3], 64)
		headingL, _ := strconv.ParseFloat(line[4], 64)
		directionL, _ := strconv.ParseFloat(line[5], 64)
		energryL, _ := strconv.ParseFloat(line[6], 64)
		speedL, _ := strconv.ParseFloat(line[7], 64)
		totalDistanceL, _ := strconv.ParseFloat(line[8], 64)

		in := SensorReading{
			TIMESTAMP:     timestampL,
			TAGID:         int32(tagIdL),
			XPOS:          xPosL,
			YPOS:          yPosL,
			HEADING:       headingL,
			DIRECTION:     directionL,
			ENERGY:        energryL,
			SPEED:         speedL,
			TOTALDISTANCE: totalDistanceL,
		}

		data, err := avro.Marshal(schema, in)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(data)
		subject := "football-iot.players." + line[1] // line[1] == player id

		nc.Publish(subject, data)
		time.Sleep(3 * time.Second)

		// out := SensorReading{}
		// err = avro.Unmarshal(schema, data, &out)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println(out)
	}

}

func main() {

	readPlayerCsv("./resources/sensor-data/player_x.csv")

}
