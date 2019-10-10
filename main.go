package main

import (
	"go-rest-api/app"
	"go-rest-api/env"
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

	err := engine.Run(":" + env.SystemPort())
	if err != nil {
		panic(err)
	}
}
