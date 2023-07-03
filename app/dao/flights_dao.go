package dao

import (
	"log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	models "app/models"
)

type FlightsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "flights"
)

// Establish a connection to database
func (m *FlightsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of flights
func (m *FlightsDAO) FindAll() ([]models.Flight, error) {
	var flights []models.Flight
	err := db.C(COLLECTION).Find(bson.M{}).All(&flights)
	return flights, err
}

// Find a flight by its id
func (m *FlightsDAO) FindById(id string) (models.Flight, error) {
	var flight models.Flight
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&flight)
	return flight, err
}

// Insert a flight into database
func (m *FlightsDAO) Insert(flight models.Flight) error {
	err := db.C(COLLECTION).Insert(&flight)
	return err
}

// Delete an existing flight
func (m *FlightsDAO) Delete(flight models.Flight) error {
	err := db.C(COLLECTION).Remove(&flight)
	return err
}

// Update an existing flight
func (m *FlightsDAO) Update(flight models.Flight) error {
	err := db.C(COLLECTION).UpdateId(flight.ID, &flight)
	return err
}
