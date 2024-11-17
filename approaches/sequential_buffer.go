package approach1singlethreaded

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type measurements struct {
	min        float64
	max        float64
	mean       float64
	totalSum   float64
	totalCount float64
}

type ComputedResult struct {
	stationName string
	max         float64
	min         float64
	mean        float64
}

func SequentialScanner(fileLoc string) string {
	stations := make(map[string]measurements)
	final_output, err := read_file_in_buffer_return_calc_results(fileLoc, &stations)
	if err != nil {
		fmt.Println("Error while reading the file", err)
		return ""
	}
	fmt.Println(final_output)
	return final_output
}

func read_file_in_buffer_return_calc_results(filepath string, stations *map[string]measurements) (string, error) {
	fileObj, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error while opening the file", err)
	}
	defer fileObj.Close()
	scanner := bufio.NewScanner(fileObj)
	for scanner.Scan() {
		line := scanner.Text()
		process_line_and_update_station(line, stations)
		if err != nil {
			return "", err
		}
	}
	return get_final_output(stations), nil

}

func get_final_output(stations *map[string]measurements) string {
	stationArr := make([]ComputedResult, 0)
	for key, station := range *stations {
		stationArr = append(stationArr, ComputedResult{stationName: key, max: math.Round(station.max*10) / 10.0, min: math.Round(station.min*10) / 10.0, mean: math.Round(station.mean*10) / 10.0})
	}
	sort.Slice(stationArr, func(i, j int) bool {
		return stationArr[i].stationName < stationArr[j].stationName
	})

	final_output := "{"
	for indx, station := range stationArr {
		if indx == len(stationArr)-1 {
			final_output += fmt.Sprintf("%s=%.1f/%.1f/%.1f}", station.stationName, station.min, station.mean, station.max)
			break
		}
		fmt.Println(station)
		final_output += fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", station.stationName, station.min, station.mean, station.max)
	}
	return final_output
}

func process_line_and_update_station(line string, STATIONS *map[string]measurements) {
	lineVals := strings.Split(line, ";")
	if len(lineVals) < 2 {
		return
	}
	stationName := lineVals[0]
	temp := lineVals[1]
	station, ok := (*STATIONS)[stationName]
	tempVal, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		fmt.Println("error here", err)
		return
	}
	if !ok {
		(*STATIONS)[stationName] = measurements{min: tempVal, max: tempVal, mean: tempVal, totalSum: tempVal, totalCount: 1}
		return
	}
	if station.max < tempVal {
		station.max = tempVal
	}
	if station.min > tempVal {
		station.min = tempVal
	}

	station.totalCount += 1
	station.totalSum += tempVal

	station.mean = station.totalSum / station.totalCount

	(*STATIONS)[stationName] = station

}
