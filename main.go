package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
)

const cmdName = "go-test-color"

var (
	colorGreen  = color.New(color.FgGreen)
	colorYellow = color.New(color.FgYellow)
	colorRed    = color.New(color.FgRed)
)

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

	// Drain stdout and stderr concurrently to avoid deadlock when the child
	// fills one pipe buffer while the other is not being read.
	output := &syncWriter{
		writer: os.Stdout,
	}

	eg := &errgroup.Group{}

	eg.Go(func() error {
		return colorOutputReader(outReader, output)
	})

	eg.Go(func() error {
		return colorErrorReader(errReader, output)
	})

	if err := eg.Wait(); err != nil {
		log.Printf("%s failed to read output: %s", cmdName, err)
		return 1
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("%s failed: %s", cmdName, err)
		return 1
	}

	return 0
}

func colorOutputReader(reader io.Reader, writer io.Writer) error {
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("bufio: reader failed to read %w", err)
		}

		line = strings.TrimSpace(line)

		if line == "" && errors.Is(err, io.EOF) {
			return nil
		}

		if strings.HasSuffix(line, "[no test files]") {
			if errors.Is(err, io.EOF) {
				return nil
			}

			continue
		}

		if strings.HasPrefix(line, "--- PASS") ||
			strings.HasPrefix(line, "PASS") ||
			strings.HasPrefix(line, "ok") {
			colorGreen.Fprintf(writer, "%s\n", line)
			if errors.Is(err, io.EOF) {
				return nil
			}

			continue
		}

		if strings.HasPrefix(line, "--- SKIP") {
			colorYellow.Fprintf(writer, "%s\n", line)
			if errors.Is(err, io.EOF) {
				return nil
			}

			continue
		}

		if strings.HasPrefix(line, "--- FAIL") ||
			strings.HasPrefix(line, "FAIL") {
			colorRed.Fprintf(writer, "%s\n", line)
			if errors.Is(err, io.EOF) {
				return nil
			}

			continue
		}

		fmt.Fprintln(writer, line)

		if errors.Is(err, io.EOF) {
			return nil
		}
	}
}

func colorErrorReader(reader io.Reader, writer io.Writer) error {
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("bufio: reader failed to read %w", err)
		}

		line = strings.TrimSpace(line)

		if line == "" && errors.Is(err, io.EOF) {
			return nil
		}

		if strings.HasPrefix(line, "# ") {
			fmt.Fprintln(writer, line)
			if errors.Is(err, io.EOF) {
				return nil
			}

			continue
		}

		// https://github.com/golang/go/issues/61229
		if strings.HasPrefix(line, "ld: warning: ") {
			colorYellow.Fprintf(writer, "%s\n", line)
			if errors.Is(err, io.EOF) {
				return nil
			}

			continue
		}

		colorRed.Fprintf(writer, "%s\n", line)

		if errors.Is(err, io.EOF) {
			return nil
		}
	}
}

// syncWriter support multiple goroutines writing to the same writer
type syncWriter struct {
	writer io.Writer
	mu     sync.Mutex
}

func (w *syncWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.writer.Write(p)
}
