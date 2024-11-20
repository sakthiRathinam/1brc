package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	approaches "github.com/sakthiRathinam/1brc/approaches"
	"github.com/stretchr/testify/assert"
)

func TestSequentialScannerLogic(t *testing.T) {
	testCases := []struct {
		fileLoc            string
		expectedOutputPath string
	}{
		{
			fileLoc:            "test_cases/measurements-1.txt",
			expectedOutputPath: "test_cases/measurements-1.out",
		},
		{
			fileLoc:            "test_cases/measurements-2.txt",
			expectedOutputPath: "test_cases/measurements-2.out",
		},

		{
			fileLoc:            "test_cases/measurements-3.txt",
			expectedOutputPath: "test_cases/measurements-3.out",
		},
	}

	for _, tc := range testCases {
		output := approaches.LineByLineApproach(tc.fileLoc)
		expectedOutput := ReadFile(tc.expectedOutputPath)
		assert.Equal(t, expectedOutput, output, "Expected and actual output should be same")
	}

}

func ReadFile(filepath string) string {
	fileObj, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error while opening the file", err)
		return ""
	}
	defer fileObj.Close()
	scanner := bufio.NewScanner(fileObj)
	var output string
	for scanner.Scan() {
		output += scanner.Text()
	}
	return output
}
