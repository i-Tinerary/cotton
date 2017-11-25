package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

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
	r.HandleFunc("/places", s.SetPlace).Methods("POST")
	r.HandleFunc("/placse/{id}", s.GetPlace).Methods("GET")

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

func (s *server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
}

func (s *server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := s.store.GetUser(vars["name"])
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, user.Data)
}

func (s *server) SetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	content, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	err = s.store.SetUser(common.User{
		Name: vars["name"],
		Data: string(content),
	})
}

func (s *server) GetPlan(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *server) SetPlace(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	var place common.Place
	err = json.Unmarshal(content, &place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	err = s.store.SetPlace(place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
}

func (s *server) GetPlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	place, err := s.store.GetPlace(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(place); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
}
