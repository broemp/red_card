package main

import (
	"database/sql"
	"fmt"
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
	fmt.Println(config.DB_Driver)
	conn, err := sql.Open(config.DB_Driver, config.DB_Source)
	if err != nil {
		log.Fatal("could not establish database connection: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = server.Start(config.WEB_Port)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
