package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"

	c "app/config"
	d "app/dao"
	m "app/models"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var config = c.Config{}
var dao = d.FlightsDAO{}
var es *elasticsearch.Client

type ESLogHook struct {
}

func (hook *ESLogHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook *ESLogHook) Fire(entry *log.Entry) error {
	// Execute your custom logic here

	entryStr, err := entry.String()
	if err != nil {
		fmt.Println("Error converting log entry to string")
	}

	req := esapi.IndexRequest{
		Index:      "log-prjctr",
		DocumentID: uuid.New().String(),
		Body:       strings.NewReader(entryStr),
		Refresh:    "true",
	}

	res, reqErr := req.Do(context.Background(), es)
	if reqErr != nil {
		fmt.Printf("Error performing the request: %s\n", reqErr)
	}
	defer res.Body.Close()

	if res.IsError() {
		fmt.Printf("Error indexing document: %s\n", res.Status())
	}

	return nil
}

// GET list of flights
func AllFlights(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting all flights")
	flights, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	retResponse(w, http.StatusOK, flights)
}

// GET a flight by its ID
func FindFlightEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Info("Looking for flight")
	params := mux.Vars(r)
	flight, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Flight ID")
		return
	}
	retResponse(w, http.StatusOK, flight)
}

// POST a new flight
func CreateFlight(w http.ResponseWriter, r *http.Request) {
	log.Info("Create Flight")
	defer r.Body.Close()
	var flight m.Flight
	if err := json.NewDecoder(r.Body).Decode(&flight); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	flight.ID = bson.NewObjectId()
	if err := dao.Insert(flight); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	retResponse(w, http.StatusCreated, flight)
}

// PUT update an existing flight
func UpdateFlight(w http.ResponseWriter, r *http.Request) {
	log.Info("Updating flight")
	defer r.Body.Close()
	var flight m.Flight
	if err := json.NewDecoder(r.Body).Decode(&flight); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(flight); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	retResponse(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing flight
func DeleteFlight(w http.ResponseWriter, r *http.Request) {
	log.Info("Delete flight")
	defer r.Body.Close()
	var flight m.Flight
	if err := json.NewDecoder(r.Body).Decode(&flight); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(flight); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	retResponse(w, http.StatusOK, map[string]string{"result": "success"})
}

func Health(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting Health")
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	retResponse(w, http.StatusOK, map[string]string{"server": name, "result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	log.Error("Error executing request", msg)
	retResponse(w, code, map[string]string{"error": msg})
}

func retResponse(w http.ResponseWriter, code int, payload interface{}) {
	log.Info("Responding", code)
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	fmt.Println("Init")
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()

	var err error
	log.SetFormatter(&log.JSONFormatter{})

	es, err = elasticsearch.NewClient(
		elasticsearch.Config{Addresses: []string{"http://127.0.0.1:9200"}})
	if err != nil {
		fmt.Println(err, es)
		os.Exit(1)
	}

	_, perr := es.Ping()
	if perr != nil {
		fmt.Println("Cannot ping elasticsearch")
		os.Exit(1)
	}

	log.Println(elasticsearch.Version)
	log.Println(es.Info())

	hook := &ESLogHook{}
	log.AddHook(hook)
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", Health)
	r.HandleFunc("/flights", AllFlights).Methods("GET")
	r.HandleFunc("/flights", CreateFlight).Methods("POST")
	r.HandleFunc("/flights", UpdateFlight).Methods("PUT")
	r.HandleFunc("/flights", DeleteFlight).Methods("DELETE")
	r.HandleFunc("/flights/{id}", FindFlightEndpoint).Methods("GET")

	log.Info("Starting Service")

	if err := http.ListenAndServe("127.0.0.1:3030", r); err != nil {
		log.Fatal(err)
	}
}
