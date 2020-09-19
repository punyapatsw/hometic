package main

import (
	// "database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// "os"
	// _ "github.com/lib/pq"
)

func PairDeviceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"active"}`))
}

func main() {
	fmt.Println("hello hometic")
	r := mux.NewRouter()
	r.HandleFunc("/pair-device", PairDeviceHandler).Methods(http.MethodPost)

	server := http.Server{
		Addr:    "127.0.0.1:2009",
		Handler: r,
	}
	log.Println("starting")
	log.Fatal(server.ListenAndServe())
}
