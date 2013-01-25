package database

import (
	"testing"
	"time"
	"database/sql"
	"fmt"
)

func TestKeyEncoding(t *testing.T) {
	key := 123
	if decodeKey(encodeKey(key)) != key {
		t.Error("Encoding and then decoding doesn't yield the same value")
	}
}

/*func TestFailedOpenDatabase(t *testing.T) {*/
	/*// TODO: How to induce an error???*/
	/*databaseName = "foo bar"*/
	/*databaseUser = "thisreallyisfake"*/
	/*db, err := OpenDatabase()*/
	/*if err == nil {*/
		/*t.Error("expected connecting to non existing db to return errors")*/
	/*}*/
	/*defer db.Close()*/
/*}*/

func TestUnmarchalParticipant(t *testing.T) {
	db, _ := OpenDatabase()
	defer cleanAndCloseDatabase(db)

	time, _ := time.Parse("2006-01-02 15:04", "2011-01-19 22:15")
	email := "harm@awesome.com"
	_, err := db.Exec("INSERT INTO participants (email, seen) VALUES($1, $2)", email, time)
	if err != nil {
	  t.Fatal(err)
	}
	rows, _ := db.Query("SELECT * FROM participants LIMIT 1")
	for rows.Next() {
		var participant Participant
		rows.Scan(&participant.Email, &participant.Seen)
		if participant.Email != email || participant.Seen != time {
			t.Fatalf("Unexpected attributes in: %v, expected %v and %v", participant, email, time)
		}
	}
}

func cleanAndCloseDatabase(db *sql.DB) {
	_, err := db.Exec("DELETE FROM participants")
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
}
