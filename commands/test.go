package main

import (
	"log"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	DatabaseLocation string
	DownloadLocation string
	APIKey string
}

func main () {
	var c Configuration
	envconfig.Process("gbvideo", &c)
	log.Printf("%#v", c)
}
