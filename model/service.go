package model

import (
	"go-rest-api/database"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Service struct {}

func NewService() Service {
	return Service{}
}

func Insert(companies Companies) []bson.ObjectId {
	session := database.Connect()
	defer session.Close()

	collection := session.DB(database.DATABASE).C(COLLECTION)

	var persisted = []bson.ObjectId{}
	for _, company := range companies {
		company.Id = bson.NewObjectId()
		collection.Insert(company)
		persisted = append(persisted, company.Id)
	}
	return persisted
}

func InsertOrUpdate(lines [][]string) {
	for i, line := range lines {
		if i != 0 {
			name := strings.ToUpper(line[0])
			zip := line[1]

			existing := findByFullNameAndZip(name, zip)
			if existing.Name != "" {
				existing.Website = strings.ToLower(line[2])
				update(existing.Id, existing)
			}
		}
	}
}

func (service *Service) FindById(id string) Company {
	session := database.Connect()
	defer session.Close()

	company := Company{}

	session.DB(database.DATABASE).C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&company)

	return company
}

func (service *Service) FindByName(name string) Companies {
	session := database.Connect()
	defer session.Close()

	companies := Companies{}
	collection := session.DB(database.DATABASE).C(COLLECTION)

	search := bson.RegEx{name + ".*", "i"}

	collection.Find(bson.M{ "name": search }).All(&companies)

	return companies
}

func (service *Service) FindByZip(zip string) Companies {
	session := database.Connect()
	defer session.Close()

	companies := Companies{}
	collection := session.DB(database.DATABASE).C(COLLECTION)

	collection.Find(bson.M{ "addresszip" : zip }).All(&companies)

	return companies
}

func (service *Service) FindByNameAndZip(name string, zip string) Companies {
	session := database.Connect()
	defer session.Close()

	companies := Companies{}
	collection := session.DB(database.DATABASE).C(COLLECTION)

	search := bson.RegEx{name + ".*", "i"}

	collection.Find(bson.M{ "name": search, "addresszip" : zip }).All(&companies)

	return companies
}

func (service *Service) FindAll() Companies {
	session := database.Connect()
	defer session.Close()

	companies := Companies{}

	session.DB(database.DATABASE).C(COLLECTION).Find(bson.M{}).All(&companies)

	return companies
}

func (service *Service) Delete(id string) bool {
	session := database.Connect()
	defer session.Close()

	company := Company{}
	collection := session.DB(database.DATABASE).C(COLLECTION)

	collection.FindId(bson.ObjectIdHex(id)).One(&company)

	if company.Name != "" {
		collection.RemoveId(company.Id)
		return true;
	}
	return false;
}

func CleanDatabase() {
	database.Clean(COLLECTION)
}

func update(id bson.ObjectId, company Company) {
	session := database.Connect()
	defer session.Close()

	session.DB(database.DATABASE).C(COLLECTION).UpdateId(id, company)
}

func findByFullNameAndZip(name string, zip string) Company {
	session := database.Connect()
	defer session.Close()

	company := Company{}

	session.DB(database.DATABASE).C(COLLECTION).Find(bson.M{"name": name, "addresszip" : zip}).One(&company)

	return company
}
