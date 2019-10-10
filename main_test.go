package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go-rest-api/app"
	"go-rest-api/database"
	"go-rest-api/env"
	"go-rest-api/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const API = "/api/v1/companies"

func Init() {
	//os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("FILE_SOURCE", "q1_catalog_test.csv")
}

func TestInsertEmpty(t *testing.T) {
	Init()

	w, persisted := insertByApi("source/empty.csv")
	defer finishProcess(persisted)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestFindAll(t *testing.T) {
	Init()

	_, persisted := insertByApi("source/q2_clientData_test.csv")
	defer finishProcess(persisted)

	request, _ := http.NewRequest("GET", API, nil)

	company := requestCompany(t, request)

	name := company.Name
	site := company.Website

	assert.Equal(t, strings.ToUpper(name), name)
	assert.Equal(t, strings.ToLower(site), site)
}

func TestFindById(t *testing.T) {
	Init()

	entity := model.Company {
		AddressZip: "24580",
		Name: "Company Limited",
		Website: "WwW.CompanyLimited.com",
	}

	objectId := insertInDatabase(entity)
	defer finishProcessUnique(objectId)

	//Not found
	request, _ := http.NewRequest("GET", API + "/" + bson.NewObjectId().Hex(), nil)
	response := httptest.NewRecorder()
	router := app.GetEngine()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNoContent, response.Code)

	//Return one
	request, _ = http.NewRequest("GET", API + "/" + objectId.Hex(), nil)
	response = httptest.NewRecorder()
	router = app.GetEngine()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	var company model.Company
	err := json.Unmarshal([]byte(response.Body.String()), &company)
	if err != nil {
		panic(err)
	}

	name := company.Name
	zip := company.AddressZip
	site := company.Website

	assert.Equal(t, entity.Name, name)
	assert.Equal(t, entity.AddressZip, zip)
	assert.Equal(t, entity.Website, site)
}

func TestFindByName(t *testing.T) {
	Init()

	_, persisted := insertByApi("source/q2_clientData_test.csv")
	defer finishProcess(persisted)

	request, _ := http.NewRequest("GET", API, nil)
	q := request.URL.Query()
	q.Add("name", "sale")
	request.URL.RawQuery = q.Encode()

	company := requestCompany(t, request)

	name := company.Name
	site := company.Website

	assert.Equal(t, strings.ToUpper(name), name)
	assert.Equal(t, strings.ToLower(site), site)
}

func TestFindByZip(t *testing.T) {
	Init()

	_, persisted := insertByApi("source/q2_clientData_test.csv")
	defer finishProcess(persisted)

	request, _ := http.NewRequest("GET", API, nil)
	q := request.URL.Query()
	q.Add("zip", "78229")
	request.URL.RawQuery = q.Encode()

	company := requestCompany(t, request)

	name := company.Name
	zip := company.AddressZip
	site := company.Website

	assert.Equal(t, strings.ToUpper(name), name)
	assert.Equal(t, "78229", zip)
	assert.Equal(t, strings.ToLower(site), site)
}

func TestFindByNameAndZip(t *testing.T) {
	Init()

	_, persisted := insertByApi("source/q2_clientData_test.csv")
	defer finishProcess(persisted)

	request, _ := http.NewRequest("GET", API, nil)
	q := request.URL.Query()
	q.Add("name", "sale")
	q.Add("zip", "78229")
	request.URL.RawQuery = q.Encode()

	company := requestCompany(t, request)

	name := company.Name
	zip := company.AddressZip
	site := company.Website

	assert.Equal(t, strings.ToUpper(name), name)
	assert.Equal(t, "78229", zip)
	assert.Equal(t, strings.ToLower(site), site)
}

func TestDelete(t *testing.T) {
	Init()

	entity := model.Company {
		AddressZip: "24580",
		Name: "Company Limited",
	}

	objectId := insertInDatabase(entity)
	defer finishProcessUnique(objectId)

	//Not found
	request, _ := http.NewRequest("DELETE", API + "/" + bson.NewObjectId().Hex(), nil)
	response := httptest.NewRecorder()
	router := app.GetEngine()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)

	//Delete one
	request, _ = http.NewRequest("DELETE", API + "/" + objectId.Hex(), nil)
	response = httptest.NewRecorder()
	router = app.GetEngine()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}

func insertByApi(fileToInsert string) (*httptest.ResponseRecorder,  []bson.ObjectId) {
	Init()

	persisted := model.LoadSourceFile()

	archive, err := os.Open(fileToInsert)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	requestPart := &bytes.Buffer{}
	writer := multipart.NewWriter(requestPart)
	part, err := writer.CreateFormFile("file", fileToInsert)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, archive)
	err = writer.Close()
	if err != nil {
		log.Println(err)
	}

	router := app.GetEngine()

	req, _ := http.NewRequest("POST", API, requestPart)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w, persisted
}

func insertInDatabase(entity model.Company) bson.ObjectId {
	Init()

	session, err := mgo.Dial(env.DatabaseHost() + ":" + env.DatabasePort())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	entity.Id = bson.NewObjectId()
	session.DB(database.DATABASE).C(model.COLLECTION).Insert(entity)
	if err != nil {
		log.Fatal(err)
	}
	return entity.Id
}

func requestCompany(t *testing.T, request *http.Request) model.Company {
	response := httptest.NewRecorder()
	router := app.GetEngine()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	var company []model.Company
	err := json.Unmarshal([]byte(response.Body.String()), &company)
	if err != nil {
		panic(err)
	}
	return company[0]
}

func finishProcessUnique(persisted bson.ObjectId) {
	finish(persisted, model.NewService())
}

func finishProcess(persisted []bson.ObjectId) {
	service := model.NewService()
	for _, id := range persisted {
		finish(id, service)
	}
}

func finish(persisted bson.ObjectId, service model.Service) {
	service.Delete(persisted.Hex())
}
