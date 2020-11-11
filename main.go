package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func main() {
	args := []string{"test"}
	args = append(args, os.Args[1:]...)

	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()

	reader, writer := io.Pipe()
	defer func() {
		if err := writer.Close(); err != nil {
			log.Printf("failed to close writer: %s", err)
		}

		if err := reader.Close(); err != nil {
			log.Printf("failed to close reader: %s", err)
		}
	}()

	cmd.Stdout = writer
	cmd.Stderr = writer

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(reader)
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

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
