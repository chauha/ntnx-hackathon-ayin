package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

const statusEndpoint = "/status"
const CCRegisterEndpoint = "/clusters/register/" //Endpoint where this agent registers

type ClusterControllerMetadata struct {
	ID            string `json:"id"`         // Id of the Cluster
	Name          string `json:"name"`       // Name of the Cluster
	Workers       int    `json:"no_workers"` // Number of worker nodes
	Masters       int    `json:"no_masters"` // Number of Master nodes
	NetworkPlugin string `json:"network_plugin,omitempty"`
}

func main() {

	log.SetPrefix(fmt.Sprintf("AYIN OnPremAgent [%v] ", os.Getpid()))
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/status/{resource}", getDeploymentStatusGeneric).Methods("GET")
	router.HandleFunc("/status/{resource}/{object}", getDeploymentStatus).Methods("GET")
	// router.HandleFunc("/status/{id}", updateEvent).Methods("PATCH")
	// router.HandleFunc("/status/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
	go regToClusterController()

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
	port := os.Getenv("CLUSTER_CONTROLLER_PORT")

	url := server + ":" + port + CCRegisterEndpoint
	fmt.Println("URL:>", url)
	var metadata ClusterControllerMetadata
	metadata.ID = executeSysCommand("sudo", []string{"cat", "/sys/class/dmi/id/product_uuid"})
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

func executeSysCommand(command string, args []string) string {
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
