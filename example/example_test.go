package example

import (
	"testing"
	"time"
)

const longTime = time.Second * 2

func TestPassShortTime(t *testing.T) {
}

func TestPassLongTime(t *testing.T) {
	time.Sleep(longTime)
}

func TestFailShorTime(t *testing.T) {
	t.Fail()
}

func TestFailLongTime(t *testing.T) {
	time.Sleep(longTime)
	t.Fail()
}

func TestFailShorTimeParallel(t *testing.T) {
	t.Parallel()
	t.Fail()
}

func TestFailLongTimeParallel(t *testing.T) {
	t.Parallel()
	time.Sleep(longTime)
	t.Fail()
}
