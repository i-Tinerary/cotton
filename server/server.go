package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/i-tinerary/cotton/common"
	"github.com/i-tinerary/cotton/store"
)

func Serve(port string, storeURL *url.URL) error {
	storage, err := store.GetStore(storeURL)
	if err != nil {
		return fmt.Errorf("error getting store: %v", err)
	}

	s := &server{store: storage}

	r := mux.NewRouter()
	r.HandleFunc("/users", s.GetUsers).Methods("GET")
	r.HandleFunc("/users/{name}", s.GetUser).Methods("GET")
	r.HandleFunc("/users/{name}", s.SetUser).Methods("POST")
	// create a plan
	r.HandleFunc("/plans", s.GetPlan).Methods("GET")
	//
	r.HandleFunc("/places/{place_id}", nil).Methods("GET")
	// get all plans sorted chronologic
	r.HandleFunc("/plans/{name}", nil)
	// get a plan by id
	r.HandleFunc("/plans/{name}/{plan_id}", nil).Methods("GET")

	return http.ListenAndServe(":"+port, r)
}

type server struct {
	store store.Interface
}

func makeResponse(w http.ResponseWriter, state int, msg string) {
	w.WriteHeader(state)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Printf("Error: Creating error response with state %d: %s", state, err)
	}
}

func (s *server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.GetUsers()
	if err != nil {
		makeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		makeResponse(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (s *server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := s.store.GetUser(vars["name"])
	if err != nil {
		makeResponse(w, http.StatusNotFound, "user not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, user.Data)
}

func (s *server) SetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	content, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		makeResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := r.Body.Close(); err != nil {
		makeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.store.SetUser(common.User{
		Name: vars["name"],
		Data: string(content),
	})
	w.WriteHeader(http.StatusOK)
}

func (s *server) GetPlan(w http.ResponseWriter, r *http.Request) {
	return
}
