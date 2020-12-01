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
	args := []string{"test"}
	args = append(args, os.Args[1:]...)

	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()

	outReader, outWriter := io.Pipe()
	defer outWriter.Close()

	errReader, errWriter := io.Pipe()
	defer errWriter.Close()

	cmd.Stdout = outWriter
	cmd.Stderr = errWriter

	go func() {
		defer outReader.Close()

		scanner := bufio.NewScanner(outReader)

		for scanner.Scan() {
			line := scanner.Text()
			trimmedLine := strings.TrimSpace(line)

			if strings.HasPrefix(trimmedLine, "--- PASS") ||
				strings.HasPrefix(line, "PASS") ||
				strings.HasPrefix(line, "ok") {
				color.Green("%s\n", line)
				continue
			}

			if strings.HasPrefix(trimmedLine, "--- FAIL") ||
				strings.HasPrefix(line, "FAIL") {
				color.Red("%s\n", line)
				continue
			}

			fmt.Println(line)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("scanner error: %s", err)
		}
	}()

	go func() {
		defer errReader.Close()

		scanner := bufio.NewScanner(errReader)

		for scanner.Scan() {
			line := scanner.Text()
			color.Red("%s\n", line)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("scanner error: %s", err)
		}
	}()

	if err := cmd.Run(); err != nil {
		return 1
	}

	return 0
}
