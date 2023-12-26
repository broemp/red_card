package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/broemp/red_card/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("could not load config:", err)
	}

	dbSource := "postgresql://" + config.DBUser + ":" + config.DBPassword + "@" + config.DBAdress + ":5432/" + config.DBDatabase + "?sslmode=disable"
	testDB, err := sql.Open(config.DBDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
