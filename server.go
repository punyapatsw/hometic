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
	"go.uber.org/zap"
)

type Pair struct {
	DeviceID int `json:"DeviceID"`
	UserID   int `json:"UserID"`
}

func PairDeviceHandler(device Device) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := zap.NewExample()
		l = l.With(zap.Namespace("hometic"), zap.String("error", "error"))
		l.Info("pair-device")

		var p Pair
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		defer r.Body.Close()
		fmt.Printf("pair: %#v\n", p)

		err = device.Pair(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.Write([]byte(`{"status":"active"}`))
	}
}

type Device interface {
	Pair(p Pair) error
}

type CreatePairDeviceFunc func(p Pair) error

func (fn CreatePairDeviceFunc) Pair(p Pair) error {
	return fn(p)
}

func NewCreatePairDevice(db *sql.DB) CreatePairDeviceFunc {
	return func(p Pair) error {
		_, err := db.Exec("INSERT INTO pairs VALUES ($1,$2);", p.DeviceID, p.UserID)
		return err
	}
}

func main() {
	fmt.Println("hello hometic")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Handle("/pair-device", PairDeviceHandler(NewCreatePairDevice(db))).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	// addr := fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT"))
	fmt.Println("addr :", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Println("starting")
	log.Fatal(server.ListenAndServe())
}
