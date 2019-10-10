package model

import "gopkg.in/mgo.v2/bson"

const COLLECTION = "companies"
const ADDRESS_ZIP_LENGTH = 5

type Company struct {
	Id bson.ObjectId `bson:"_id"`
	Name string `json:"name"`
	AddressZip string `json:"zip"`
	Website string `json:"website"`
}

type Companies []Company