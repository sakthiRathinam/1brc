package approaches

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type StationsStore struct {
	Stations *map[string]measurements
	lock     *sync.RWMutex
}

func (storeObj *StationsStore) updateStations(line string) {
	lineVals := strings.Split(line, ";")
	stationName := lineVals[0]
	temp := lineVals[1]
	storeObj.lock.Lock()
	defer storeObj.lock.Unlock()

	station, ok := (*storeObj.Stations)[stationName]
	tempVal, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		fmt.Println("error while parsing the temp val")
	}
	storeObj.lock.Lock()
	defer storeObj.lock.Unlock()
	if !ok {
		(*storeObj.Stations)[stationName] = measurements{min: tempVal, max: tempVal, mean: tempVal, totalSum: tempVal, totalCount: 1}
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

	(*storeObj.Stations)[stationName] = station
}

func createStationStore() StationsStore {
	stations := make(map[string]measurements)
	return StationsStore{Stations: &stations, lock: &sync.RWMutex{}}
}

func ThreadedBuffer(fileLoc string, chunkSize int) string {
	stationStore := createStationStore()
	final_output, err := read_file_in_chunk_and_assign_to_goroutine(fileLoc, chunkSize, &stationStore)
	if err != nil {
		fmt.Println("Error while reading the file", err)
		return ""
	}
	fmt.Println(final_output)
	return final_output
}

func read_file_in_chunk_and_assign_to_goroutine(filepath string, chunkSize int, stationStoreObj *StationsStore) (string, error) {
	fmt.Println("Chunk size", chunkSize)
	fileObj, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error while opening the file", err)
	}
	defer fileObj.Close()

	reader := bufio.NewReaderSize(fileObj, chunkSize)
	buffer := make([]byte, chunkSize)
	count := 0
	partialLine := []byte{}

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				stationStoreObj.updateStations(string(partialLine))
				break
			}
			fmt.Println("Error while reading the file" + err.Error())
		}

		buffer = append(partialLine, buffer[:n]...)
		chunkString := string(buffer)

		splittedLines := strings.Split(chunkString, "\n")
		partialLine = []byte(splittedLines[len(splittedLines)-1])
		for index, chunkLine := range strings.Split(chunkString, "\n") {
			if index == len(splittedLines)-1 {
				break
			}
			go stationStoreObj.updateStations(chunkLine)
		}

	}

	fmt.Println("total processed lines", count)
	return get_final_output(stationStoreObj.Stations), nil
}
