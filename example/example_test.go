package example

import (
	"testing"
	"time"
)

const longTime = 2 * time.Second

func TestPassShortTime(t *testing.T) {
}

func TestPassLongTime(t *testing.T) {
	time.Sleep(longTime)
}

func TestFailShortTime(t *testing.T) {
	t.Fail()
}

func TestFailLongTime(t *testing.T) {
	time.Sleep(longTime)
	t.Fail()
}

func TestFailShortTimeParallel(t *testing.T) {
	t.Parallel()
	t.Fail()
}

func TestFailLongTimeParallel(t *testing.T) {
	t.Parallel()
	time.Sleep(longTime)
	t.Fail()
}
