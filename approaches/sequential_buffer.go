package approach1singlethreaded

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type measurements struct {
	min  int
	max  int
	mean int
}

type station struct {
	name         string
	measurements measurements
}

func SequentialBuffer() {
	fileLoc := "data/measurements.txt"
	read_file_in_buffer_return_calc_results(fileLoc)
}

func read_file_in_buffer_return_calc_results(filepath string) (measurements, error) {
	min := 0
	max := 0
	total_sum := 0
	count := 0
	fileObj, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error while opening the file", err)
	}
	defer fileObj.Close()
	bufferSize := 64 * 1024 * 1024 // 64 MB
	reader := bufio.NewReaderSize(fileObj, bufferSize)
	buffer := make([]byte, bufferSize)
	for {
		_, err := reader.Read(buffer)
		count += 1
		if err != nil {
			return measurements{}, nil
		}
		parse_buffer_get_calcs(buffer)
		if count == 3 {
			break
		}
	}
	measurements := measurements{min: min, max: max, mean: total_sum / count}
	return measurements, nil

}

func parse_buffer_get_calcs(buffer []byte) (int, int, int) {
	textFormat := string(buffer)
	lines := strings.Split(textFormat, "/n")
	fmt.Println(lines, "got the lines")
	return len(lines), 0, 0
}
