package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"runtime/pprof"
)

/*
	When given a list of JSON objects like {"value": "..."},
	return how many contained the word "buffalo".

	Must read the records from os.Stdin and write the result to os.Stdout.
*/

var profile = flag.Bool("profile", false, "write profiles")

func main() {
	flag.Parse()
	// CPU Profiling (if enabled)
	if *profile {
		log.Println("profiling")
		f, err := os.Create("cpu.profile")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	// Count buffalos
	n := Count(os.Stdin)
	fmt.Println(n)

	// Memory Profiling (if enabled)
	if *profile {
		f, err := os.Create("mem.profile")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			panic(err)
		}
	}
}

type Record struct {
	Value string `json:"value"`
}

func Count(r io.Reader) int {
	n := 0
	b := bufio.NewReader(r)
	for {
		s, err := b.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		var record Record
		if err := json.Unmarshal([]byte(s), &record); err != nil {
			panic(err)
		}

		rx := regexp.MustCompile(`.*(\s+)buffalo(\s+).*`)
		if rx.Match([]byte(record.Value)) {
			n++
		}
	}
	return n
}
