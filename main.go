package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

type server struct {
	db redis.Conn
}

func makeResponse(w http.ResponseWriter, state int, msg string) {
	w.WriteHeader(state)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Printf("Error: Creating error response with state %d: %s", state, err)
	}
}

func (s *server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := redis.Strings(s.db.Do("SMEMBERS", "users"))
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

	pref, err := redis.String(s.db.Do("GET", "U:"+vars["name"]))
	if err != nil {
		makeResponse(w, http.StatusNotFound, "user not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, pref)
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
	_, err = s.db.Do("SADD", "users", vars["name"])
	if err != nil {
		makeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = s.db.Do("SET", "U:"+vars["name"], string(content))
	if err != nil {
		makeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) GetPlan(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	port := os.Getenv("PORT")
	redisURL := os.Getenv("REDIS_URL")

	c, err := redis.DialURL(redisURL)
	if err != nil {
		log.Fatalf("connecting to redis on %q: %v", redisURL, err)
	}
	defer c.Close()

	s := &server{db: c}

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

	log.Fatal(http.ListenAndServe(":"+port, r))
}
