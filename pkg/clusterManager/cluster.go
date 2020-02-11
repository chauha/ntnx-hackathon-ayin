package clusterManager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
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
	id := ExecuteSysCommand("sudo", []string{"cat", "/sys/class/dmi/id/product_uuid"})
	pingStruct.ID = id
	pingStruct.Health = "UP"
	pingStruct.Masters = 1
	pingStruct.Workers = 3
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

const CCRegisterEndpoint = "/clusters/register/"

//Register cluster to CCP
func (cm *ClusterManager) regToClusterController() {
	url := cm.ClusterControllerBaseURL + CCRegisterEndpoint
	fmt.Println("URL:>", url)
	var metadata ClusterControllerMetadata
	metadata.ID = ExecuteSysCommand("sudo", []string{"cat", "/sys/class/dmi/id/product_uuid"})
	metadata.Name = "Demo"
	metadata.Masters = 1 //TODO strconv.Atoi(executeSysCommand("kubectl", []string{"get", "no"}))
	metadata.Workers = 2 //TODO strconv.Atoi(executeSysCommand("/bin/sh", []string{"kubectl", "get", "no", "|", "grep", "--ignore-case", "\"Master\"", "|", "tail", "-n", "+2", "|", "wc", "-l"}))

	j, _ := json.Marshal(metadata)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

//Execute a command
func ExecuteSysCommand(command string, args []string) string {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return stderr.String()
	}
	return out.String()
}
