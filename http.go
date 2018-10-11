package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type progressBarWebServer struct {
	router *mux.Router
}

func NewWebServer(db ProgressDB) progressBarWebServer {
	router := mux.NewRouter()
	router.Handle("/{id}", GetProgressBarHandler(db)).Methods("GET")
	router.Handle("/{id}", UpdateProgressBarHandler(db)).Methods("PUT")
	router.Handle("/", CreateProgressBarHandler(db)).Methods("POST")
	return progressBarWebServer{router}
}

func (ws *progressBarWebServer) Serve(address string) {
	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      ws.router,
	}
	srv.ListenAndServe()
}

func GetProgressBarHandler(db ProgressDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		progress := db.Get(vars["id"])
		if "application/text" == r.Header.Get("Accept") {
			fmt.Fprintf(w, "%d", progress)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"id": vars["id"], "progress": progress})
		}
	})
}

func CreateProgressBarHandler(db ProgressDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewV4().String()
		token := newToken()
		db.Create(id, token, 0) // TODO: Starting progress

		if "application/text" == r.Header.Get("Accept") {
			fmt.Fprintf(w, "ID: %s\nToken: %s", id, token)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "token": token, "progress": 0})
		}
	})
}

func UpdateProgressBarHandler(db ProgressDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		r.ParseForm()
		token := r.Form.Get("token")
		progress := r.Form.Get("progress")
		iProgress, err := strconv.Atoi(progress)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !db.Update(vars["id"], token, uint32(iProgress)) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid Token")
			return
		}

		if "application/text" == r.Header.Get("Accept") {
			fmt.Fprintf(w, "%d", iProgress)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"id": vars["id"], "progress": iProgress})
		}
	})
}
