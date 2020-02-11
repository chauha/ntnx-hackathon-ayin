package status

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

func GetDeploymentStatusGeneric(w http.ResponseWriter, r *http.Request) {
	resourceID := mux.Vars(r)["resource"]
	result := GetStatusGeneric(resourceID)
	fmt.Fprintf(w, result)
}

func GetDeploymentStatus(w http.ResponseWriter, r *http.Request) {

	objectID := mux.Vars(r)["object"]
	result := GetStatus("deployment", objectID)
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
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
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
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return stderr.String()
	}
	return out.String()
}
