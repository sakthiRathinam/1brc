package main

import (
	"fmt"
	"math/rand"
	"os"
)

var STATIONS = []string{"arabic", "brazil", "canada", "denmark", "egypt", "france", "germany", "hungary", "india", "japan", "korea", "london", "mexico", "nigeria", "oman", "paris", "qatar", "russia", "spain", "texas", "ukraine", "vietnam", "washington", "xray", "yemen", "zambia"}

func GenerateFakeMeasurements(numberOfMeasurements int, fileLoc string) {
	fileObj, _ := os.Create(fileLoc)
	for i := 0; i < numberOfMeasurements; i++ {
		// generate random station name
		stationName := STATIONS[rand.Intn(len(STATIONS))]
		// generate random temperature
		randTemp := rand.Float64() * 100
		// write to file
		station_str := fmt.Sprintf("%s;%.2f\n", stationName, randTemp)
		fileObj.WriteString(station_str)
	}
	defer fileObj.Close()
}
