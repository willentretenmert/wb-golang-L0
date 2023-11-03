package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"sync"
)

var cache sync.Map

func main() {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     "db-container",
		Port:     5432,
		Database: "L0_db",
		User:     "some_user",
		Password: "zxczxc",
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	SubscribeToNATS(conn)

	loadCacheFromDB(conn)

	fmt.Println(cache.Load("b563feb7b2b84b6test"))
	//http.ListenAndServe(":8000", nil)
	select {}
}
