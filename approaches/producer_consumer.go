package approaches

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var workerConcurrency = 0

func ProducerConsumerApproach(fileLoc string, chunkSize int) string {
	resultMap := make(map[string]measurements)
	chunkStream := make(chan []string)
	resultStream := make(chan map[string]measurements)
	if workerConcurrency == 0 {
		workerConcurrency = runtime.NumCPU() - 1
	}
	// spawn workers
	for i := 0; i < workerConcurrency; i++ {
		go chunkProcessWorker(&chunkStream, &resultStream)
	}
	// read file and put the chunk in chunkstream
	chunkProducer(fileLoc, chunkSize, &chunkStream)
	// merge the results
	fmt.Println(resultMap)
	return ""
	// get the finaloutput
}

func chunkProcessWorker(chunkChannel *chan []string, resultChannel *chan map[string]measurements) {
	for chunkString := range *chunkChannel {
		fmt.Println(len(chunkString))
	}
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
				// add that missing line
				break
			}
			fmt.Println("Error while reading the file" + err.Error())
		}

		buffer = append(partialLine, buffer[:n]...)

		chunkString := string(buffer)

		splittedLines := strings.Split(chunkString, "\n")
		partialLine = []byte(splittedLines[len(splittedLines)-1])
		*chunkChannel <- splittedLines
	}

	fmt.Println("total processed lines", count)
}
