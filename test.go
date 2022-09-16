package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func readPlayerCsv(csvPath string) {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(csvFile)
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(csvLines)

}

func main() {
	playerCsvs := []string{
		"./resources/sensor-data/player_1.csv",
		"./resources/sensor-data/player_2.csv",
		"./resources/sensor-data/player_3.csv",
		"./resources/sensor-data/player_5.csv",
		"./resources/sensor-data/player_6.csv",
		"./resources/sensor-data/player_7.csv",
		"./resources/sensor-data/player_8.csv",
		"./resources/sensor-data/player_9.csv",
		"./resources/sensor-data/player_10.csv",
		"./resources/sensor-data/player_11.csv",
		"./resources/sensor-data/player_12.csv",
		"./resources/sensor-data/player_13.csv",
		"./resources/sensor-data/player_14.csv",
		"./resources/sensor-data/player_15.csv",
	}

	for _, file := range playerCsvs {

		readPlayerCsv(file)

	}
	//readPlayerCsv("./resources/sensor-data/player_x.csv")
}
