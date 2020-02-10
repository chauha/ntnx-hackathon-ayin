package curate

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/db"
	"github.com/pkg/errors"
)

type CurateClustersService struct {
	Db            db.ClusterStorage
	WebserverPort int
}

func (cs *CurateClustersService) RunService() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serviceHTTPHandle(w, r, cs)
	})
	server := &http.Server{}
	log.Printf("webserverPort=%d", cs.WebserverPort)
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: cs.WebserverPort})
	if err != nil {
		log.Printf("ListenTCP: %s", err)
		return errors.Wrap(err, "Failed to create a TCP listener.")
	}
	return server.Serve(listener)
}

func serviceHTTPHandle(w http.ResponseWriter, req *http.Request, cs *CurateClustersService) {
	log.Printf("serviceHTTPHandle")
	log.Printf("req.URL.Query %v", req.URL.Query())
	if req.URL.Path == "/clusters/register/" {
		var c db.ClusterControllerMetadata
		err := json.NewDecoder(req.Body).Decode(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ClusterControllerMetadata: %+v", c)
		cs.Db.InsertOrUpdateCluster(&c)
		w.Write([]byte("ok"))
	} else {
		http.Error(w, "Path not found", http.StatusNotFound)
	}
}
