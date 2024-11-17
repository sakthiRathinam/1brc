package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	approaches "github.com/sakthiRathinam/1brc/approaches"
)

func main() {
	start_time := time.Now()
	fileLoc := "data/measurements.txt"
	if len(os.Args) == 3 {
		fmt.Println("Usage: ./main <number_of_measurements> <output_file>")
		total_measurements, _ := strconv.Atoi(os.Args[1])
		GenerateFakeMeasurements(total_measurements, os.Args[2])
		fileLoc = os.Args[2]
	}

	approaches.SequentialScanner(fileLoc)
	elapsed_time := time.Since(start_time)
	fmt.Println("Time taken to execute the program", elapsed_time)

}
