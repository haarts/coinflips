package database

import (
	"github.com/speps/go-hashids"
	_ "github.com/bmizerany/pq"
	"database/sql"
	"fmt"
	"time"
)

// TODO: move to settings file
var (
	databaseName = "coinflips"
	databaseUser = "harm"
)

type Coinflip struct {
	Head	string
	Tail	string
	Result	string
	Id		int
}

type Participant struct {
	Email		string
	Seen		time.Time
	CoinflipId	int
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
	h.Salt = "dit is zout, heel zout" //TODO: put this in a conf file
	return h.Decrypt(key)[0]
}

func encodeKey(key int) string {
	h := hashids.New()
	h.MinLength = 10
	h.Salt = "dit is zout, heel zout" //TODO: put this in a conf file
	return h.Encrypt([]int{key})
}

func (coinflip *Coinflip) EncodeKey() string {
	return encodeKey(coinflip.Id)
}

func (coinflip *Coinflip) FindParticipantByEmail(email string) (Participant, error) {
	return Participant{}, nil
}

func (coinflip *Coinflip) NumberOfUnregisteredParticipants() (int, error) {
	return 0, nil
}

func (coinflip *Coinflip) FindParticipants() []Participant {
	return []Participant{Participant{}}
}

func (participant *Participant) Update() error {
	return nil
}

func (coinflip *Coinflip) CreateParticipant(participant *Participant) error {
	return nil
}

func (coinflip *Coinflip) Create() error {
	return nil
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
