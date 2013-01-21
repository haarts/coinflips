package database

import (
	"testing"
)

func TestKeyEncoding(t *testing.T) {
	key := 123
	if decodeKey(encodeKey(key)) != key {
		t.Error("Encoding and then decoding doesn't yield the same value")
	}
}

func TestFailedOpenDatabase(t *testing.T) {
	// TODO: How to induce an error???
	databaseName = "foo bar"
	databaseUser = "thisreallyisfake"
	db, err := OpenDatabase()
	if err == nil {
		t.Error("expected connecting to non existing db to return errors")
	}
	defer db.Close()
}
