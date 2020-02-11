package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/clusterManager"
	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/connectAgent"
	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/register"
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
	go register.RegToClusterController()

	log.Print("Starting On premise nutanix K8s Connect agent")
	cm := clusterManager.ClusterManager{
		PingIntervalInSeconds: 100,
	}
	go cm.RunService()

	ca := connectAgent.ConnectAgent{
		WebserverAddress: ":8080",
	}

	errC := ca.RunService()
	if errC != nil {
		log.Printf("Connect Agent failed to start with error: %+v", errC)
		os.Exit(1)
	}
}
