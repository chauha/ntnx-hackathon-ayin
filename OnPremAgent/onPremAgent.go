package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const cCRegEndpoint = "/clusters/register/" //Endpoint where this agent registers
const defaultCCServer = "http://localhost"
const defaultCCPort = "9090"
const getMachineUUID = "sudo cat /sys/class/dmi/id/product_uuid"
const getMasterCMD = "kubectl get no | grep -i \"Master\" | wc -l"
const getWorkerCMD = "kubectl get no | grep -i -v \"Master\" | tail -n +2 | wc -l"

type ClusterControllerMetadata struct {
	ID            string `json:"id"`         // Id of the Cluster
	Name          string `json:"name"`       // Name of the Cluster
	Workers       int    `json:"no_workers"` // Number of worker nodes
	Masters       int    `json:"no_masters"` // Number of Master nodes
	NetworkPlugin string `json:"network_plugin,omitempty"`
}

func main() {

	log.SetPrefix(fmt.Sprintf("AYIN OnPremAgent [%v] ", os.Getpid()))
	go regToClusterController()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/status/{resource}", getDeploymentStatusGeneric).Methods("GET")
	router.HandleFunc("/status/{resource}/{object}", getDeploymentStatus).Methods("GET")
	router.HandleFunc("/create", createResource).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createResource(w http.ResponseWriter, r *http.Request) {
	tempFileName := "/tmp/" + genFileName()
	f, err := os.OpenFile(tempFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	io.Copy(f, r.Body)
	result := executeSysCommand("kubectl create -f " + tempFileName)
	fmt.Println("Result ", result)
	fmt.Fprintf(w, result)
}

func getDeploymentStatusGeneric(w http.ResponseWriter, r *http.Request) {
	resourceID := mux.Vars(r)["resource"]
	result := getStatusGeneric(resourceID)
	fmt.Fprintf(w, result)
}

func getDeploymentStatus(w http.ResponseWriter, r *http.Request) {

	objectID := mux.Vars(r)["object"]
	result := getStatus("deployment", objectID)
	fmt.Fprintf(w, result)
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func getStatusGeneric(resource string) string {

	cmd := exec.Command("kubectl", "get", resource, "-o", "json")
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

func getStatus(resource string, object string) string {
	cmd := exec.Command("kubectl", "get", resource, object, "-o", "json")
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

func regToClusterController() {
	server := os.Getenv("CLUSTER_CONTROLLER_URL")
	if server == "" {
		server = defaultCCServer
	}
	port := os.Getenv("CLUSTER_CONTROLLER_PORT")
	if port == "" {
		port = defaultCCPort
	}

	url := server + ":" + port + cCRegEndpoint

	var metadata ClusterControllerMetadata
	var err error
	metadata.ID = executeSysCommand(getMachineUUID)
	metadata.Name = "Demo"
	metadata.Masters, err = strconv.Atoi(executeSysCommand(getMasterCMD))
	if err != nil {
		panic(err)
	}
	metadata.Workers, err = strconv.Atoi(executeSysCommand(getWorkerCMD))
	if err != nil {
		panic(err)
	}
	fmt.Println("metatdata  ", metadata)
	j, _ := json.Marshal(metadata)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Can't connect to Server")
	} else {
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
}

func executeSysCommand(cmd string) string {
	command := exec.Command("bash", "-c", cmd)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stdErr
	err := command.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stdErr.String())
		return stdErr.String()

	}
	stringOut := strings.TrimSuffix(stdOut.String(), "\n")
	return stringOut

}

func genFileName() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
