package main

import (
	"log"
	"net/url"

	"github.com/fatih/color"
	"github.com/jacobkaufmann/arxivlib-users/api"
	"github.com/jacobkaufmann/arxivlib-users/datastore"
)

var (
	baseURL *url.URL
)

func main() {
	datastore.Connect()

	m := api.Handler()

	color.Magenta("Listening")
	log.Fatal(m.Run(":8080"))
}
