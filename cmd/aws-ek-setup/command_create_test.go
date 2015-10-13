package main

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestARNParser(t *testing.T) {
	arn := getAccountIdentifier("arn:aws:iam::170889777123:user/wolfeidau")
	assert.Equal(t, arn, "170889777123", "they should be equal")
}
