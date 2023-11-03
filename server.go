package main

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"net/http"
	"sync"
)

var (
	hostIP string
	cache  sync.Map
)

func main() {

	hostIP = "localhost"
	//hostIP = "db-container"

	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     hostIP,
		Port:     5432,
		Database: "L0_db",
		User:     "some_user",
		Password: "zxczxc",
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	loadCacheFromDB(conn, &cache)
	//fmt.Println(">>> CACHE LOADED: ")
	//fmt.Println(cache.Load("993somemoreuid666"))
	//fmt.Println(cache.Load("b563feb7b2b84b6test"))
	http.HandleFunc("/", cacheHandler)
	http.ListenAndServe(":8000", nil)

	SubscribeToNATS(conn)

	select {}
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for cache data")
	cacheData := make(map[string]interface{})
	cache.Range(func(key, value interface{}) bool {
		cacheData[key.(string)] = value
		return true
	})

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(cacheData)
	if err != nil {
		fmt.Println("Error encoding cache data:", err)
	}
}
