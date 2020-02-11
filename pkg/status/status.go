package status

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

func GetDeploymentStatusGeneric(w http.ResponseWriter, r *http.Request) {
	resourceID := mux.Vars(r)["resource"]
	result := GetStatusGeneric(resourceID)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}

func GetDeploymentStatus(w http.ResponseWriter, r *http.Request) {
	objectID := mux.Vars(r)["object"]
	result := GetStatus("deployment", objectID)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}

func HomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func GetStatusGeneric(resource string) string {
	cmd := exec.Command("kubectl", "get", resource, "-o", "json")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		return stderr.String()
	}
	return out.String()
}

func GetStatus(resource string, object string) string {
	cmd := exec.Command("kubectl", "get", resource, object, "-o", "json")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		return stderr.String()
	}
	return out.String()
}
