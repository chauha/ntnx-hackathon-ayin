package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/curate"
	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/db"
)

func handleSignals(c chan os.Signal) {
	s := <-c
	log.Printf("Got signal %s", s)
	os.Exit(1)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	go handleSignals(c)

	log.Print("Starting Curate Clusters Service")

	fs, err := db.NewFileClusterStorage("")
	if err != nil {
		log.Printf("Cannot init storage: %+v", err)
		os.Exit(1)
	}
	cs := curate.CurateClustersService{
		Db:            fs,
		WebserverPort: 9090,
	}
	err = cs.RunService()
	if err != nil {
		log.Printf("Service ended with error: %+v", err)
		os.Exit(1)
	}
}
