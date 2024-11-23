package approaches

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func SequentialBuffer(fileLoc string, chunkSize int) string {
	stations := make(map[string]measurements)
	final_output, err := read_file_in_chunk_and_process(fileLoc, chunkSize, &stations)
	if err != nil {
		fmt.Println("Error while reading the file", err)
		return ""
	}
	fmt.Println(final_output)
	return final_output
}

func read_file_in_chunk_and_process(filepath string, chunkSize int, stations *map[string]measurements) (string, error) {
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
				process_line_and_update_station(string(partialLine), stations)
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
			process_line_and_update_station(chunkLine, stations)
		}

	}

	fmt.Println("total processed lines", count)
	return get_final_output(stations), nil
}
