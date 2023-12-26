package main

import (
	"database/sql"
	"log"

	"github.com/broemp/red_card/api"
	db "github.com/broemp/red_card/db/sqlc"
	"github.com/broemp/red_card/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	dbSource := "postgresql://" + config.DBUser + ":" + config.DBPassword + "@" + config.DBAdress + ":5432/" + config.DBDatabase + "?sslmode=disable"
	conn, err := sql.Open(config.DBDriver, dbSource)
	if err != nil {
		log.Fatal("could not establish database connection: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = server.Start(config.WebPort)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
