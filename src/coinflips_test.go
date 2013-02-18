package main

import (
	"testing"
)

func TestStripUrls(t *testing.T) {
	withUrl := "hi there go to spam.me"
	if stripUrls(withUrl) != "hi there go to " {
		t.Error("Failed to strip")
	}
	withoutUrl := "ken touch this"
	if stripUrls(withoutUrl) != "ken touch this" {
		t.Error("Stripped when not needed")
	}
}
