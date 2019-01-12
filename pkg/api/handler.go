package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jacobkaufmann/arxivlib-users/datastore"
)

var store = datastore.NewDatastore(datastore.DB)

// Handler handles incoming API requests
func Handler() *gin.Engine {
	m := gin.Default()

	m.GET("/api/users/:id", serveUser)
	m.GET("/api/users", serveUsers)
	m.POST("/api/users", createUser)

	return m
}
