package create

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/clusterManager"
)

func CreateResource(w http.ResponseWriter, r *http.Request) {
	tempFileName := "/tmp/" + clusterManager.GenFileName()
	f, err := os.OpenFile(tempFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Println("Error : ", err)
	}
	io.Copy(f, r.Body)
	result := clusterManager.ExecuteSysCommand("kubectl create -f " + tempFileName + " -o json")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	log.Println("Result ", result)
	fmt.Fprintf(w, result)
}

func DeleteResource(w http.ResponseWriter, r *http.Request) {
	//TODO
}
