package clusterManager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type ClusterControllerMetadata struct {
	ID            string `json:"id"`   // Id of the Cluster
	Name          string `json:"name"` // Name of the Cluster
	Health        string `json:"health,omitempty"`
	Workers       int    `json:"no_workers"` // Number of worker nodes
	Masters       int    `json:"no_masters"` // Number of Master nodes
	NetworkPlugin string `json:"network_plugin,omitempty"`
}

type Ping struct {
	ID      string `json:"id"` // Id of the Cluster
	Health  string `json:"health"`
	Workers int    `json:"no_workers"` // Number of worker nodes
	Masters int    `json:"no_masters"` // Number of Master nodes
}

type ClusterManager struct {
	ClusterControllerBaseURL string
	PingIntervalInSeconds    time.Duration
	lastPingTime             time.Time
}

func (cm *ClusterManager) RunService() error {
	cm.regToClusterController()

	pingStruct := &Ping{}
	for {
		updatePingArgs(pingStruct)
		cm.pingCloudAgent(pingStruct)
		cm.lastPingTime = time.Now()
		time.Sleep(cm.PingIntervalInSeconds * time.Second)

	}
}

func updatePingArgs(pingStruct *Ping) {
	id := ExecuteSysCommand(getMachineUUID)
	pingStruct.ID = id
	pingStruct.Health = "UP"
	pingStruct.Masters, _ = strconv.Atoi(ExecuteSysCommand(getMasterCMD))
	pingStruct.Workers, _ = strconv.Atoi(ExecuteSysCommand(getWorkerCMD))
}

func (cm *ClusterManager) pingCloudAgent(pingArgs *Ping) error {
	url := cm.ClusterControllerBaseURL + "/ping"
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

const cCRegEndpoint = "/clusters/register/" //Endpoint where this agent registers
const getMachineUUID = "sudo cat /sys/class/dmi/id/product_uuid"
const getMasterCMD = "kubectl get no | grep -i \"Master\" | wc -l"
const getWorkerCMD = "kubectl get no | grep -i -v \"Master\" | tail -n +2 | wc -l"

//Register cluster to CCP
func (cm *ClusterManager) regToClusterController() {
	url := cm.ClusterControllerBaseURL + cCRegEndpoint
	var metadata ClusterControllerMetadata
	var err error
	metadata.ID = ExecuteSysCommand(getMachineUUID)
	metadata.Name = "Demo"
	metadata.Masters, err = strconv.Atoi(ExecuteSysCommand(getMasterCMD))
	if err != nil {
		panic(err)
	}
	metadata.Workers, err = strconv.Atoi(ExecuteSysCommand(getWorkerCMD))
	if err != nil {
		panic(err)
	}
	j, _ := json.Marshal(metadata)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Can't connect to Server")
	} else {
		defer resp.Body.Close()
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response Body:", string(body))
	}
}

//Execute a command
func ExecuteSysCommand(cmd string) string {
	command := exec.Command("bash", "-c", cmd)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stdErr
	err := command.Run()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + stdErr.String())
		return stdErr.String()

	}
	return strings.TrimSuffix(stdOut.String(), "\n")
}

// This guy generates a 10 character long random string
func GenFileName() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
