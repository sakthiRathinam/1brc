package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	approaches "github.com/sakthiRathinam/1brc/approaches"
)

var cpu_profile = flag.String("cpuprofile", "", "write cpu profile to file")
var mem_profile = flag.String("memprofile", "", "write memory profile to this file")
var file_location = flag.String("fileloc", "data/measurements.txt", "get the file for processing")

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
	start_time := time.Now()
	approaches.SequentialScanner(*file_location)
	elapsed_time := time.Since(start_time)
	fmt.Println("Time taken to execute the program", elapsed_time)

}
