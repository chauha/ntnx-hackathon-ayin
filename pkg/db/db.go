package db

import "github.com/nutanix/ntnx-hackathon-ayin/v1/pkg/clusterManager"

type ClusterStorage interface {
	InsertOrUpdateCluster(c *clusterManager.ClusterControllerMetadata) error
	ListClusters() ([]clusterManager.ClusterControllerMetadata, error)
	Get(id string) *clusterManager.ClusterControllerMetadata
}
