package main

import (
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Port     int  `short:"p" description:"The port on which the server runs on" default:"8081"`
	Insecure bool `short:"i" description:"Start the API using http. Not recommended"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting server")
	startServer(opts.Port, opts.Insecure)
}
