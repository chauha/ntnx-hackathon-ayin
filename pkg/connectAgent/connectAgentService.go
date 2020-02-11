package connectAgent

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/status"
)

//Specify Port
type ConnectAgent struct {
	WebserverAddress string
}

// Runs connect Agent
func (ca *ConnectAgent) RunService() error {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", status.HomeLink)
	router.HandleFunc("/status/{resource}", status.GetDeploymentStatusGeneric).Methods("GET")
	router.HandleFunc("/status/{resource}/{object}", status.GetDeploymentStatus).Methods("GET")
	log.Fatal(http.ListenAndServe(ca.WebserverAddress, router))

	return nil

}
