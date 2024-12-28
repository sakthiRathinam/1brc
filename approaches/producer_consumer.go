package approaches

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

var workerConcurrency = 0

func ProducerConsumerApproach(fileLoc string, chunkSize int) string {
	resultMap := make(map[string]measurements)
	chunkStream := make(chan []string)
	resultStream := make(chan map[string]measurements)
	countStream := make(chan int)
	waitGroup := sync.WaitGroup{}
	if workerConcurrency == 0 {
		workerConcurrency = runtime.NumCPU() - 1
	}
	fmt.Println(workerConcurrency, "worker concurrencly")
	// spawn workers
	for i := 0; i < workerConcurrency; i++ {
		waitGroup.Add(1)
		go chunkProcessWorker(&chunkStream, &resultStream, &countStream, &waitGroup)
	}
	// read file and put the chunk in chunkstream
	go chunkProducer(fileLoc, chunkSize, &chunkStream)
	// merge the results
	fmt.Println(resultMap)
	totalRowsProcessed := 0
	go func() {
		waitGroup.Wait()
		close(countStream)
	}()
	for count := range countStream {
		totalRowsProcessed += count
	}
	fmt.Println("this many rows got processed", totalRowsProcessed)
	return ""
	// get the finaloutput
}

func chunkProcessWorker(chunkChannel *chan []string, resultChannel *chan map[string]measurements, countStream *chan int, waitGroup *sync.WaitGroup) {
	for chunkString := range *chunkChannel {
		*countStream <- len(chunkString)
	}
	defer waitGroup.Done()
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
