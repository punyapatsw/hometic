package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/punyapatsw/hometic/logger"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Pair struct {
	DeviceID int `json:"DeviceID"`
	UserID   int `json:"UserID"`
}

type CustomResponseWriter interface {
	JSON(statusCode int, data interface{})
}

type JSONResponseWriter struct {
	http.ResponseWriter
}

type CustomHandlerFunc func(w CustomResponseWriter, r *http.Request)

func (handler CustomHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(&JSONResponseWriter{w}, r)
}

func (w *JSONResponseWriter) JSON(statusCode int, data interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func PairDeviceHandler(device Device) CustomHandlerFunc {
	return func(w CustomResponseWriter, r *http.Request) {
		// i := r.Context().Value("logger")
		logger.L(r.Context()).Info("pair-device")
		// log := i.(*zap.Logger)
		// log.Info("pair-device")
		var p Pair
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.JSON(http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()
		fmt.Printf("pair: %#v\n", p)

		err = device.Pair(p)
		if err != nil {
			w.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// w.JSON(http.StatusOK, []byte(`{"status":"active"}`))
		w.JSON(http.StatusOK, map[string]interface{}{"status": "active"})

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

func run() error {
	fmt.Println("hello hometic")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	r.Use(logger.Middleware)
	r.Handle("/pair-device", PairDeviceHandler(NewCreatePairDevice(db))).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr :", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Println("starting")
	return server.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatal("can't start application", err)
	}

}
