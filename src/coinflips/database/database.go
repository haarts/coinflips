package database

import (
	"github.com/speps/go-hashids"
	pq "github.com/bmizerany/pq"
	"database/sql"
	"fmt"
	"time"
	"coinflips/settings"
)

// TODO: move to settings file
var (
	databaseName = "coinflips"
	databaseUser = "harm"
	salt = settings.ReadSalt()
)

type Coinflip struct {
	Head	string
	Tail	string
	Result	sql.NullString
	Id		int
}

type Participant struct {
	Email		string
	Seen		pq.NullTime
	CoinflipId	int
	Id			int
}

func OpenDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", databaseUser, databaseName))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func decodeKey(key string) int {
	h := hashids.New()
	h.Salt = salt
	return h.Decrypt(key)[0]
}

func encodeKey(key int) string {
	h := hashids.New()
	h.MinLength = 10
	h.Salt = salt
	return h.Encrypt([]int{key})
}

func (coinflip *Coinflip) EncodedKey() string {
	return encodeKey(coinflip.Id)
}

func (coinflip *Coinflip) SetResult(result string) {
	coinflip.Result = sql.NullString{String: result}
}

func (coinflip *Coinflip) FindParticipantByEmail(email string) (Participant, error) {
	db, _ := OpenDatabase()
	defer db.Close()

	var participant Participant
	err := db.QueryRow("SELECT id, email, seen, coinflip_id FROM participants WHERE email = $1 AND coinflip_id = $2", email, coinflip.Id).Scan(&participant.Id, &participant.Email, &participant.Seen, &participant.CoinflipId)
	if err != nil {
		return Participant{}, err
	}
	return participant, nil
}

func (coinflip *Coinflip) NumberOfUnregisteredParticipants() (int, error) {
	db, _ := OpenDatabase()
	defer db.Close()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM participants WHERE seen IS NULL AND coinflip_id = $1", coinflip.Id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (coinflip *Coinflip) FindParticipants() []Participant {
	db, _ := OpenDatabase()
	defer db.Close()

	var participants []Participant
	rows, _ := db.Query("SELECT id, email, seen FROM participants WHERE coinflip_id = $1", coinflip.Id)
	for rows.Next() {
		var participant Participant
		rows.Scan(&participant.Id, &participant.Email, &participant.Seen)
		participants = append(participants, participant)
	}
	return participants
}


func (participant *Participant) Register() error {
	participant.Seen = pq.NullTime{Time: time.Now()}
	return participant.Update()
}

func (participant *Participant) Update() error {
	db, _ := OpenDatabase()
	defer db.Close()
	
	_, err := db.Exec("UPDATE participants SET email = $1, seen = $2 WHERE id = $3", participant.Email, participant.Seen.Time, participant.Id)
	return err
}

func (coinflip *Coinflip) CreateParticipant(participant *Participant) error {
	db, _ := OpenDatabase()
	defer db.Close()

	var id int
	err := db.QueryRow("INSERT INTO participants (email, coinflip_id) VALUES($1, $2) RETURNING id", participant.Email, coinflip.Id).Scan(&id)
	if err != nil {
		return err
	}
	participant.CoinflipId = coinflip.Id
	participant.Id = id
	return nil
}

func (coinflip *Coinflip) Create() error {
	db, _ := OpenDatabase()
	defer db.Close()

	var id int
	err := db.QueryRow("INSERT INTO coinflips (head, tail) VALUES($1, $2) RETURNING id", coinflip.Head, coinflip.Tail).Scan(&id)
	if err != nil {
		return err
	}
	coinflip.Id = id
	return nil
}

func (coinflip *Coinflip) Update() error {
	db, _ := OpenDatabase()
	defer db.Close()
	
	_, err := db.Exec("UPDATE coinflips SET head = $1, tail = $2, result = $3 WHERE id = $4", coinflip.Head, coinflip.Tail, coinflip.Result.String, coinflip.Id)
	return err
}

func FindCoinflip(key string) (*Coinflip, error) {
	keyId := decodeKey(key)
	db, _ := OpenDatabase()
	defer db.Close()

	row := db.QueryRow("SELECT id, head, tail, result FROM coinflips WHERE id = $1", keyId)
	
	coinflip := new(Coinflip)
	if err := row.Scan(&coinflip.Id, &coinflip.Head, &coinflip.Tail, &coinflip.Result); err != nil {
	  return nil, err
	}
	return coinflip, nil
}

func TotalNumberOfParticipants() int {
	db, _ := OpenDatabase()
	defer db.Close()
	var count int
	db.QueryRow("SELECT COUNT(*) FROM participants").Scan(&count)
	return count
}
