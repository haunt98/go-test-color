package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

const cmdName = "go-test-color"

func main() {
	code := runGoTest()
	os.Exit(code)
}

// Run go test with args
func runGoTest() int {
	// Pass all args
	args := make([]string, 0, len(os.Args))
	args = append(args, "test")
	args = append(args, os.Args[1:]...)
	cmd := exec.CommandContext(context.Background(), "go", args...)

	// Read stdout and stderr
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("%s failed to get stdout pipe: %s", cmdName, err)
		return 1
	}

	errReader, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("%s failed to get stderr pipe: %s", cmdName, err)
		return 1
	}

	// See https://stackoverflow.com/questions/8875038/redirect-stdout-pipe-of-child-process-in-go
	if err := cmd.Start(); err != nil {
		log.Printf("%s failed to start: %s", cmdName, err)
		return 1
	}

	// Add color to both stdout and stderr
	if err := colorOutputReader(outReader); err != nil {
		log.Printf("%s failed to read stdout: %s", cmdName, err)
	}

	if err := colorErrorReader(errReader); err != nil {
		log.Printf("%s failed to read stderr: %s", cmdName, err)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("%s failed to wait: %s", cmdName, err)
		return 1
	}

	return 0
}

func colorOutputReader(reader io.Reader) error {
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil &&
			err != io.EOF {
			return fmt.Errorf("bufio: reader failed to read %w", err)
		}

		if line == "" &&
			err == io.EOF {
			return nil
		}

		line = strings.TrimSpace(line)

		if strings.HasSuffix(line, "[no test files]") {
			continue
		}

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

		if err == io.EOF {
			return nil
		}
	}
}

func colorErrorReader(reader io.Reader) error {
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil &&
			err != io.EOF {
			return fmt.Errorf("bufio: reader failed to read %w", err)
		}

		if line == "" &&
			err == io.EOF {
			return nil
		}

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "# ") {
			fmt.Println(line)
			continue
		}

		// https://github.com/golang/go/issues/61229
		if strings.HasPrefix(line, "ld: warning: ") {
			color.Yellow("%s\n", line)
			continue
		}

		color.Red("%s\n", line)

		if err == io.EOF {
			return nil
		}
	}
}
