package clusterManager

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/register"

	log "github.com/sirupsen/logrus"
)

type Ping struct {
	Id               string     `json:"id"`
	Health           string		`json:"health"`
	NoWorkersRunning int        `json:"noWorkersRunning"`
	NoMastersRunning int		`json:"noMastersRunning`
}

type ClusterManager struct {
	PingIntervalInSeconds time.Duration
	lastPingTime          time.Time
}

func (cm *ClusterManager) RunService() error {
	pingStruct := &Ping{}
	for {
		updatePingArgs(pingStruct)
		cm.pingCloudAgent(pingStruct)
		cm.lastPingTime = time.Now()
		time.Sleep(cm.PingIntervalInSeconds * time.Second)

	}

}

func updatePingArgs(pingStruct *Ping) {

	id := register.ExecuteSysCommand("sudo", []string{"cat", "/sys/class/dmi/id/product_uuid"})
	pingStruct.Id = id
	pingStruct.Health = "UP"
	pingStruct.NoMastersRunning = 1
	pingStruct.NoWorkersRunning = 3
}

func (cm *ClusterManager) pingCloudAgent(pingArgs *Ping) error {

	server := os.Getenv("CLUSTER_CONTROLLER_URL")
	port := os.Getenv("CLUSTER_CONTROLLER_PORT")
	url := server + ":" + port + "/ping"
	jBody, _ := json.Marshal(pingArgs)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println("Error while reaching cloud controller %s", err)
	}
	return nil
}
