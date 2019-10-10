package database

import (
	"gopkg.in/mgo.v2"
	"go-rest-api/env"
)

const (
	DATABASE = "dinizDB"
)

func Connect() *mgo.Session {
	session, err := mgo.Dial(env.DatabaseHost() + ":" + env.DatabasePort())
	if err != nil {
		panic(err)
	}
	return session
}

func Clean(collection string) {
	session := Connect()
	defer session.Close()

	session.DB(DATABASE).C(collection).RemoveAll(nil)
}
