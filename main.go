package main

import (
	"go-rest-api/app"
	"go-rest-api/model"
	"log"
)

func main() {
	log.Println(":: Preparando banco de dados ::")

	model.CleanDatabase()

	log.Println(":: Carregando dados iniciais ::")

	model.LoadSourceFile()

	log.Println(":: Inicializando aplicação ::")

	engine := app.GetEngine()

	log.Println(":: Aplicação pronta!!! ::")

	err := engine.Run(":8082")
	if err != nil {
		panic(err)
	}
}
