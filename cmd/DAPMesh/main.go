package main

import (
	"flag"
	"log"

	"github.com/philmacfly/DAPMesh/pkg/config"
	"github.com/philmacfly/DAPMesh/pkg/gossiper"
)

var confpath string

func init() {
	flag.StringVar(&confpath, "c", "./config.toml", "Path to the configuration toml file")
	flag.Parse()
}

func main() {
	c, err := config.LoadConfig(confpath)
	if err != nil {
		log.Fatalln("Error loadin config:" + err.Error())
	}
	_, _ = gossiper.StartGossiper(c.Gossiper)

}
