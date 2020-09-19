package main

import (
	// "database/sql"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/mux"
	// _ "github.com/lib/pq"
)

func PairDeviceHandler(w http.ResponseWriter, r *http.Request) {
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
