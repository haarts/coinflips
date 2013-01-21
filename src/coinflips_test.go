package main

import (
	"testing"
)

func TestEncodeKey(t *testing.T) {
	key := 123
	if decodeKey(encodeKey(123)) != 123 {
		t.Error("Encoding and then decoding doesn't yield the same value")
	}
}
