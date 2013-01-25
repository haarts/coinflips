package database

import (
	"testing"
	"time"
	"database/sql"
)

func TestKeyEncoding(t *testing.T) {
	key := 123
	if decodeKey(encodeKey(key)) != key {
		t.Error("Encoding and then decoding doesn't yield the same value")
	}
}

func TestUnmarchalParticipant(t *testing.T) {
	db, _ := OpenDatabase()
	defer cleanAndCloseDatabase(db)

	time, _ := time.Parse("2006-01-02 15:04", "2011-01-19 22:15")
	email := "harm@awesome.com"
	_, err := db.Exec("INSERT INTO participants (email, seen, coinflip_id) VALUES($1, $2, 0)", email, time)
	if err != nil {
	  t.Fatal(err)
	}
	rows, _ := db.Query("SELECT email, seen FROM participants LIMIT 1")
	for rows.Next() {
		var participant Participant
		rows.Scan(&participant.Email, &participant.Seen)
		if participant.Email != email || participant.Seen != time {
			t.Fatalf("Unexpected attributes in: %v, expected %v and %v", participant, email, time)
		}
	}
}

func TestCreateCoinflip(t *testing.T) {
	db, _ := OpenDatabase()
	defer cleanAndCloseDatabase(db)

	coinflip := Coinflip{ Head: "head", Tail: "tail" }
	coinflip.Create()
	
	row := db.QueryRow("SELECT id, head, tail FROM coinflips WHERE id = $1", coinflip.Id)
	var storedCoinflip Coinflip
	row.Scan(&storedCoinflip.Id, &storedCoinflip.Head, &storedCoinflip.Tail)

	if storedCoinflip != coinflip {
		t.Fatalf("Expected %v to be equal to %v", coinflip, storedCoinflip)
	}
}

func TestCreateParticipant(t *testing.T) {
	db, _ := OpenDatabase()
	defer cleanAndCloseDatabase(db)

	coinflip := Coinflip{ Head: "head", Tail: "tail" }
	coinflip.Create()

	participant := Participant{ Email: "harm" }
	coinflip.CreateParticipant(&participant)

	var storedParticipant Participant
	db.QueryRow("SELECT id, email, coinflip_id FROM participants WHERE coinflip_id = $1", coinflip.Id).Scan(&storedParticipant.Id, &storedParticipant.Email, &storedParticipant.CoinflipId)
	if participant.Id == 0 || storedParticipant != participant {
		t.Fatalf("Expected %v to equal %v", storedParticipant, participant)
	}
}

/*func TestFindParticipantByEmail(t *testing.T) {*/
	/*db, _ := OpenDatabase()*/
	/*defer cleanAndCloseDatabase(db)*/

	/*email := "harm@awesome.com"*/
	/*_, err := db.Exec("INSERT INTO participants (email) VALUES($1)", email)*/
	/*_, err := db.Exec("INSERT INTO participants (email) VALUES($1)", "some@other.com")*/
/*}*/

func cleanAndCloseDatabase(db *sql.DB) {
	db.Exec("DELETE FROM participants")
	db.Exec("DELETE FROM coinflips")
	db.Close()
}
