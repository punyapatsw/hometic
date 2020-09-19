package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Pair struct {
	DeviceID int `json:"DeviceID"`
	UserID   int `json:"UserID"`
}

type PairDeviceHandler struct {
	createPairDevice CreatePairDevice
}

func (ph *PairDeviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var p Pair
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()
	fmt.Printf("pair: %#v\n", p)

	err = ph.createPairDevice(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Write([]byte(`{"status":"active"}`))
}

type CreatePairDevice = func(p Pair) error

func createPairDevice(p Pair) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO pairs VALUES ($1,$2);", p.DeviceID, p.UserID)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("hello hometic")
	r := mux.NewRouter()
	r.Handle("/pair-device", &PairDeviceHandler{createPairDevice: createPairDevice}).Methods(http.MethodPost)

	// addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	addr := fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT"))
	fmt.Println("addr :", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Println("starting")
	log.Fatal(server.ListenAndServe())
}
