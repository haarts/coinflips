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
	databaseUser = "postgres"
)

/*type Coinflip struct {*/
	/*Head	string*/
	/*Tail	string*/
	/*Result	string*/
	/*Id		int*/
/*}*/

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

func encodeKey(key int) string  {
	h := hashids.New()
	h.MinLength = 10
	h.Salt = "dit is zout, heel zout" //TODO: put this in a conf file
	return h.Encrypt([]int{key})
}
