package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	approaches "github.com/sakthiRathinam/1brc/approaches"
)

var cpu_profile = flag.String("cpuprofile", "", "write cpu profile to file")
var mem_profile = flag.String("memprofile", "", "write memory profile to this file")
var file_location = flag.String("fileloc", "data/measurements.txt", "get the file for processing")
var chunk_size = flag.String("chunksize", "60000", "read it given chunk size")
var generate_fake_measurements = flag.String("generateFake", "0", "Number if you want to generate fake measuremetns")
var approach = flag.String("approach", "producerConsumer", "Approach to be used for processing")

var approachesMap = map[string]func(string, int) string{
	"threadedBuffer":   approaches.ThreadedBuffer,
	"sequentialBuffer": approaches.SequentialBuffer,
	"lineByline":       approaches.LineByLineApproach,
	"producerConsumer": approaches.ProducerConsumerApproach,
}

func main() {
	flag.Parse()

	if *cpu_profile != "" {
		cpuProfileFile, err := os.Create("./profiles/" + *cpu_profile)
		if err != nil {
			panic("Error while creating the cpu profile file" + err.Error())
		}
		defer cpuProfileFile.Close()
		if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
			panic("Error while starting the cpu profiler" + err.Error())
		}

		defer pprof.StopCPUProfile()

	}

	chunk_size, err := strconv.Atoi(*chunk_size)
	if chunk_size != 60000 {
		chunk_size = chunk_size * 1024 * 1024
	}
	if err != nil {
		panic("Chunk size should be an integer without decimal points" + err.Error())
	}

	no_of_fake_measurements, err := strconv.Atoi(*generate_fake_measurements)

	if no_of_fake_measurements != 0 && err == nil {
		GenerateFakeMeasurements(no_of_fake_measurements, "data/fake_measurements.txt")
		*file_location = "data/fake_measurements.txt"
	}

	selected_apporach, ok := approachesMap[*approach]
	if !ok {
		panic("Approach not found")
	}

	start_time := time.Now()
	selected_apporach(*file_location, chunk_size)
	elapsed_time := time.Since(start_time)
	fmt.Println("Time taken to execute the program", elapsed_time)

	if *mem_profile != "" {
		memProfileFile, err := os.Create("./profiles/" + *mem_profile)
		if err != nil {
			panic("Error while creating the memory profile file" + err.Error())
		}
		defer memProfileFile.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
			panic("Error while starting the cpu profiler" + err.Error())
		}
	}
}
