package main

import (
	// "database/sql"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	// _ "github.com/lib/pq"
)

type Pair struct {
	DeviceID int `json:"DeviceID"`
	UserID   int `json:"UserID"`
}

func PairDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var p Pair
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()
	fmt.Printf("pair : %#v\n", p)
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to database error", err)
	}
	defer db.Close()
	stmt, err := db.Prepare(`INSERT INTO pairs (deviceid, userid)
							values ($1, $2)`)
	if err != nil {
		log.Fatal(err)
		return
	}
	if _, err = stmt.Exec(p.DeviceID, p.UserID); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("insert table success.")

	w.Write([]byte(`{"status":"active"}`))
}

func main() {
	fmt.Println("hello hometic")
	r := mux.NewRouter()
	r.HandleFunc("/pair-device", PairDeviceHandler).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr :", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Println("starting")
	log.Fatal(server.ListenAndServe())
}
