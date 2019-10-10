package exception

import (
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"log"
	"net/http"
)

func ReturnErrorStacktrace(c *gin.Context, err error) {
	LogStacktrace(err)
	c.JSON(http.StatusBadRequest, c.Error(err))
}

func ReturnError(c *gin.Context, message string) {
	log.Println(message)
	c.JSON(http.StatusBadRequest, message)
}

func LogStacktrace(err error) {
	err = errors.Wrap(err, 0)
	log.Println(err.(*errors.Error).ErrorStack())
}