package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLastMidnight(t *testing.T) {
	now := time.Now()
	lastMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	assert.Equal(t, lastMidnight, LastMidnight())
}

func TestNextMidnight(t *testing.T) {
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	assert.Equal(t, nextMidnight, NextMidnight())
}
