package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/db"
)

const CCRegisterEndpoint = "/clusters/register/"

//Register cluster to CCP
func RegToClusterController() {
	server := os.Getenv("CLUSTER_CONTROLLER_URL")
	port := os.Getenv("CLUSTER_CONTROLLER_PORT")

	url := server + ":" + port + CCRegisterEndpoint
	fmt.Println("URL:>", url)
	var metadata db.ClusterControllerMetadata
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
