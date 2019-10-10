package router

import (
	"bufio"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"go-rest-api/exception"
	"go-rest-api/model"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
	"os"
)

type Route struct {
	service model.Service
}

func NewRoute() *Route {
	return &Route{
		service: model.NewService(),
	}
}

func (route *Route) BuildRoutes(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		group := v1.Group("/companies")
		{
			group.GET("", route.Get)
			group.GET("/:id", route.GetUnique)
			group.POST("", route.Insert)
			group.DELETE("/:id", route.Delete)
		}
	}
}

// @Summary Show a list of Companies
// @Description Get a list of Companies
// @Accept  json
// @Produce  json
// @Param name query string false "the name or part of name"
// @Param zip query string false "the address zip of company"
// @Success 200 {object} model.Companies
// @Failure 204 {string} string "No content"
// @Router /companies [get]
func (route *Route) Get(c *gin.Context) {
	name := c.Query("name")
	zip := c.Query("zip")

	response := model.Companies{}

	if name != "" && zip != "" {
		response = route.service.FindByNameAndZip(name, zip)
	} else if name != "" && zip == "" {
		response = route.service.FindByName(name)
	} else if name == "" && zip != "" {
		response = route.service.FindByZip(zip)
	} else {
		response = route.service.FindAll()
	}

	if len(response) > 0 {
		c.JSON(http.StatusOK, response)
	} else {
		c.Status(http.StatusNoContent)
	}
}

// @Summary Show a Company
// @Description Get a Company
// @Accept  json
// @Produce  json
// @Param id path string true "Identifier of company"
// @Success 200 {object} model.Company
// @Failure 204 {string} string "No content"
// @Router /companies/{id} [get]
func (route *Route) GetUnique(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		c.Status(http.StatusNoContent)
		return
	}

	response := route.service.FindById(id)

	if response.Name != "" {
		c.JSON(http.StatusOK, response)
	} else {
		c.Status(http.StatusNoContent)
	}
}

type CSV struct {}

// @Summary Insert Companies by file
// @Description Post Companies
// @Accept  mpfd
// @Produce  json
// @Param file formData string true "CSV file with companies"
// @Success 200 {object} model.Company
// @Failure 422 {string} string "Unprocessable"
// @Router /companies [post]
func (route *Route) Insert(c *gin.Context) {
	formFile, header, err := c.Request.FormFile("file")
	if err != nil {
		exception.ReturnErrorStacktrace(c, err)
		return
	}
	defer formFile.Close()

	f, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		exception.ReturnErrorStacktrace(c, err)
		return
	}
	defer f.Close()

	io.Copy(f, formFile)
	uploadedFile, _ := os.Open(header.Filename)

	lines := model.FileToList(csv.NewReader(bufio.NewReader(uploadedFile)))
	if len(lines) == 0 {
		c.Status(http.StatusUnprocessableEntity)
	} else {
		model.InsertOrUpdate(lines)
		c.Status(http.StatusCreated)
	}
}

// @Summary Remove a Company
// @Description Delete a Company
// @Accept  json
// @Produce  json
// @Param id path string true "Identifier of company"
// @Success 200 {object} model.Company
// @Failure 422 {string} string "Unprocessable company"
// @Router /companies/{id} [delete]
func (route *Route) Delete(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		c.Status(http.StatusNoContent)
		return
	}

	if route.service.Delete(id) {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusUnprocessableEntity)
	}
}