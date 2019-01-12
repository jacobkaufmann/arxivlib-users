package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	arxivlib "github.com/jacobkaufmann/arxivlib-users"
)

func serveUser(c *gin.Context) {
	user := c.Param("username")
	opt := &arxivlib.UserListOptions{Username: user}

	users, err := store.Users.List(opt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	if err = writeJSON(c.Writer, users[0]); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
	}
}

func serveUsers(c *gin.Context) {
	opt := &arxivlib.UserListOptions{}
	if err := c.ShouldBindQuery(opt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := store.Users.List(opt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	if users == nil {
		users = []*arxivlib.User{}
	}

	if err = writeJSON(c.Writer, users); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
	}
}

func authenticateUser(c *gin.Context) {
	username := c.PostForm("username")
	passwd := c.PostForm("passwd")

	user, err := store.Users.Authenticate(username, passwd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, "status: not found")
		return
	}

	if err = writeJSON(c.Writer, user); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
	}
}

func createUser(c *gin.Context) {
	var user arxivlib.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := store.Users.Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
	}
}
