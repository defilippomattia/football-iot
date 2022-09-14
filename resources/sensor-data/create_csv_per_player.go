package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	csvFile, err := os.Open("./interpolated_zxy.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines {

		fileName := "./player_" + line[1] + ".csv"

		_, err := os.Stat(fileName)

		if os.IsNotExist(err) {
			newf, err := os.Create(fileName)
			if err != nil {
				fmt.Println(err)
			}
			newf.Close()
		}

		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		w := csv.NewWriter(f)
		w.Write(line)
		w.Flush()
		f.Close()

	}

}
