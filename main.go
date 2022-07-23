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

func runGoTest() int {
	// Run go test

	args := []string{"test"}
	args = append(args, os.Args[1:]...)

	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()

	// Output pipe and error pipe

	outReader, outWriter := io.Pipe()
	defer outReader.Close()
	defer outWriter.Close()

	errReader, errWriter := io.Pipe()
	defer errReader.Close()
	defer errWriter.Close()

	cmd.Stdout = outWriter
	cmd.Stderr = errWriter

	go colorOutputReader(outReader)
	go colorErrorReader(errReader)

	if err := cmd.Run(); err != nil {
		return 1
	}

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
