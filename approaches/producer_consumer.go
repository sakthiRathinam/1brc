package approaches

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var workerConcurrency = 0

func ProducerConsumerApproach(fileLoc string, chunkSize int) string {
	resultMap := make(map[string]measurements)
	chunkStream := make(chan []string)
	resultStream := make(chan map[string]measurements)
	waitGroup := sync.WaitGroup{}
	if workerConcurrency == 0 {
		workerConcurrency = runtime.NumCPU() - 1
	}
	fmt.Println(workerConcurrency, "worker concurrencly")
	// spawn workers
	for i := 0; i < workerConcurrency; i++ {
		waitGroup.Add(1)
		go chunkProcessWorker(&chunkStream, &resultStream, &waitGroup)
	}
	// read file and put the chunk in chunkstream
	go chunkProducer(fileLoc, chunkSize, &chunkStream)
	// merge the results
	fmt.Println(resultMap)
	totalRowsProcessed := 0
	go func() {
		waitGroup.Wait()
		close(resultStream)
	}()

	for tempMap := range resultStream {
		for stationName, measurement := range tempMap {
			station, ok := resultMap[stationName]
			if !ok {
				resultMap[stationName] = measurement
			}

			if station.max < measurement.max {
				station.max = measurement.max
			}
			if station.min > measurement.min {
				station.min = measurement.min
			}

			station.totalCount += measurement.totalCount

			station.totalSum += measurement.totalSum

			station.mean = station.totalSum / station.totalCount

			resultMap[stationName] = station

		}
	}
	fmt.Println("this many rows got processed", totalRowsProcessed)
	finalOutput := finalFormatting(&resultMap)
	fmt.Println(finalOutput)
	return finalOutput
}

func chunkProcessWorker(chunkChannel *chan []string, resultChannel *chan map[string]measurements, waitGroup *sync.WaitGroup) {
	for chunkString := range *chunkChannel {
		processedTempMap := processChunk(chunkString)
		*resultChannel <- processedTempMap
	}
	defer waitGroup.Done()
}
func processChunk(fileLines []string) map[string]measurements {
	tempMap := make(map[string]measurements)
	for _, line := range fileLines {
		lineVals := strings.Split(line, ";")
		if len(lineVals) != 2 {
			continue
		}
		stationName := lineVals[0]
		temperature := lineVals[1]

		station, ok := tempMap[stationName]

		tempVal, err := strconv.ParseFloat(temperature, 64)
		if err != nil {
			continue
		}

		if !ok {
			tempMap[stationName] = measurements{min: tempVal, max: tempVal, mean: tempVal, totalSum: tempVal, totalCount: 1}
			continue
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

	}
	return tempMap
}

func chunkProducer(filepath string, chunkSize int, chunkChannel *chan []string) {
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
				fmt.Println(partialLine)
				break
			}
			fmt.Println("Error while reading the file" + err.Error())
		}

		buffer = append(partialLine, buffer[:n]...)

		chunkString := string(buffer)

		splittedLines := strings.Split(chunkString, "\n")
		partialLine = []byte(splittedLines[len(splittedLines)-1])
		splittedLines = splittedLines[:len(splittedLines)-1]
		*chunkChannel <- splittedLines
	}
	close(*chunkChannel)

	fmt.Println("total processed lines", count)
}
