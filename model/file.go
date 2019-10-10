package model

import (
	"encoding/csv"
	"go-rest-api/env"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strings"
)

const (
	PATH = "./source/"
)

func LoadSourceFile() []bson.ObjectId {
	archive, err := os.Open(PATH + env.FileSource())
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	lines := FileToList(csv.NewReader(archive))
	if len(lines) == 0 {
		log.Println("The source file is empty!")
		return []bson.ObjectId{}
	}

	errors := Companies{}
	companies := Companies{}
	for i, line := range lines {
		if i != 0 {
			zip := line[1]
			entity := Company{Name: strings.ToUpper(line[0]), AddressZip: zip}

			if len(zip) == ADDRESS_ZIP_LENGTH {
				companies = append(companies, entity)
			} else {
				errors = append(errors, entity)
			}
		}
	}
	persisted := Insert(companies)

	logErrors(errors)

	return persisted
}

func FileToList(archive *csv.Reader) [][]string {
	reader := archive
	reader.Comma = env.FileSourceSeparator()

	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	return lines
}

func logErrors(errors Companies) {
	if len(errors) > 0 {
		log.Println(" ** Has not imported items because address zip is invalid: **")
		for _, entity := range errors {
			log.Println("   -- Name: " + entity.Name + ". Zip: " + entity.AddressZip)
		}
	}
}
