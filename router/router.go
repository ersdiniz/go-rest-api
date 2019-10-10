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
	group := router.Group("/v1/companies")
	{
		group.GET("", route.Get)
		group.GET("/:id", route.GetUnique)
		group.POST("", route.Insert)
		group.DELETE("/:id", route.Delete)
	}
}

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