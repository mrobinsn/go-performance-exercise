package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Build the attempt binary
	cmd := exec.Command("go", "build", "-o", "attempt", ".")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.WithError(err).
			WithField("output", out.String()).
			Fatal("failed to build your attempt")
	}

	// Find test files
	files, err := ioutil.ReadDir("data")
	if err != nil {
		log.WithError(err).Fatal("failed to list data files")
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".gz") {
			log.Fatal("you didn't decompress the test data yet, naughty")
		}
		if !file.IsDir() && strings.HasPrefix(file.Name(), "raw") {
			run(filepath.Join("data", file.Name()))
		}
	}
}

func run(file string) {
	fmt.Printf("=== Running against '%s' ===\n", file)
	// Read in the test data
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.WithError(err).Fatal("failed to read the data file")
	}

	// Run golden (simple)
	cmd := exec.Command("./golden-simple")
	cmd.Stdin = bytes.NewReader(b)
	var out bytes.Buffer
	cmd.Stdout = &out
	start := time.Now()
	if err := cmd.Run(); err != nil {
		log.WithError(err).Fatal("golden simple failed")
	}
	attemptDur := time.Since(start)

	n, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		log.WithError(err).Fatal("golden simple did not output a number")
	}
	fmt.Printf("Golden (simple) returned %d in %s (1x)\n", n, attemptDur)
	expectedAnswer := n
	bestDur := attemptDur

	// Run golden (complicated)
	cmd = exec.Command("./golden-complicated")
	cmd.Stdin = bytes.NewReader(b)
	out.Reset()
	cmd.Stdout = &out
	start = time.Now()
	if err := cmd.Run(); err != nil {
		log.WithError(err).Fatal("golden complicated failed")
	}
	attemptDur = time.Since(start)

	n, err = strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		log.WithError(err).Fatal("golden complicated did not output a number")
	}
	fmt.Printf("Golden (complicated) returned %d in %s (%.2fx)\n", n, attemptDur, float64(attemptDur)/float64(bestDur))

	// Run your attempt
	cmd = exec.Command("./attempt")
	cmd.Stdin = bytes.NewReader(b)
	out.Reset()
	cmd.Stdout = &out
	start = time.Now()
	if err := cmd.Run(); err != nil {
		log.WithError(err).Fatal("your attempt failed")
	}
	attemptDur = time.Since(start)

	n, err = strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		log.WithError(err).Fatal("your attempt did not output a number")
	}

	if n != expectedAnswer {
		log.WithField("expected", expectedAnswer).
			WithField("your_answer", n).
			Fatal("your attempt did not output the correct answer")
	}

	fmt.Printf("Your attempt returned %d in %s (%.2fx)\n", n, attemptDur, float64(attemptDur)/float64(bestDur))
	fmt.Print("======\n\n")
}
