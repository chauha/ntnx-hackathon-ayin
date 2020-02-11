package db

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type FileClusterStorage struct {
	filePath string
	data     map[string]ClusterControllerMetadata
}

func NewFileClusterStorage(filePath string) (*FileClusterStorage, error) {
	fs := &FileClusterStorage{}
	if filePath != "" {
		fs.filePath = filePath
	} else {
		fs.filePath = "FileClusterStorage.json"
	}
	fs.data = make(map[string]ClusterControllerMetadata)
	err := fs.loadFile()
	return fs, err
}

func (fs *FileClusterStorage) InsertOrUpdateCluster(c *ClusterControllerMetadata) error {
	log.Printf("Storing cluster %+v", c)
	fs.data[c.ID] = *c
	return fs.updateFile()
}

func (fs *FileClusterStorage) ListClusters() ([]ClusterControllerMetadata, error) {
	values := make([]ClusterControllerMetadata, 0, len(fs.data))
	for _, v := range fs.data {
		values = append(values, v)
	}
	return values, nil
}

func (fs *FileClusterStorage) Get(id string) *ClusterControllerMetadata {
	c := fs.data[id]
	return &c
}

func (fs *FileClusterStorage) loadFile() error {
	_, err := os.Stat(fs.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	jsonBytes, err := ioutil.ReadFile(fs.filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, &fs.data)
	log.Printf("Read from storage %+v", fs.data)
	return err
}

func (fs *FileClusterStorage) updateFile() error {
	jsonString, err := json.Marshal(fs.data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fs.filePath, []byte(jsonString), 0644)
	return err
}
