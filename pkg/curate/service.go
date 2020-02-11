package curate

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/clusterManager"
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
	log.Printf("req.URL.Path %s", req.URL.Path)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if strings.HasPrefix(req.URL.Path, "/clusters/register") {
		var c clusterManager.ClusterControllerMetadata
		err := json.NewDecoder(req.Body).Decode(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ClusterControllerMetadata: %+v", c)
		err = cs.Db.InsertOrUpdateCluster(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok"))
	} else if strings.HasPrefix(req.URL.Path, "/clusters") {
		c, err := cs.Db.ListClusters()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(c)
	} else if strings.HasPrefix(req.URL.Path, "/ping") {
		var p clusterManager.Ping
		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Ping: %+v", p)
		c := cs.Db.Get(p.ID)
		if c != nil {
			c.Health = p.Health
			c.Masters = p.Masters
			c.Workers = p.Workers
			err = cs.Db.InsertOrUpdateCluster(c)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Cluster not found", http.StatusNotFound)
		}
	} else {
		http.Error(w, "Path not found", http.StatusNotFound)
	}
}
