package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func main() {
	code := runGoTest()
	os.Exit(code)
}

// Run go test with args
func runGoTest() int {
	// Pass all args
	args := []string{"test"}
	args = append(args, os.Args[1:]...)
	cmd := exec.Command("go", args...)

	// Output pipe
	outReader, outWriter := io.Pipe()
	defer outReader.Close()
	defer outWriter.Close()

	// Error pipe
	errReader, errWriter := io.Pipe()
	defer errReader.Close()
	defer errWriter.Close()

	// Redirect cmd pipes to our pipes
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter

	// See https://stackoverflow.com/questions/8875038/redirect-stdout-pipe-of-child-process-in-go
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start: %s", err)
		return 1
	}
	defer func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("Failed to wait: %s", err)
		}
	}()

	go colorOutputReader(outReader)
	go colorErrorReader(errReader)

	return 0
}

func colorOutputReader(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "--- PASS") ||
			strings.HasPrefix(line, "PASS") ||
			strings.HasPrefix(line, "ok") {
			color.Green("%s\n", line)
			continue
		}

		if strings.HasPrefix(line, "--- SKIP") {
			color.Yellow("%s\n", line)
			continue
		}

		if strings.HasPrefix(line, "--- FAIL") ||
			strings.HasPrefix(line, "FAIL") {
			color.Red("%s\n", line)
			continue
		}

		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %s", err)
	}
}

func colorErrorReader(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		color.Red("%s\n", line)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %s", err)
	}
}
