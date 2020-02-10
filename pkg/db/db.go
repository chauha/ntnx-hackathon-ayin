package db

type ClusterControllerMetadata struct {
	ID            string `json:"id"`         // Id of the Cluster
	Name          string `json:"name"`       // Name of the Cluster
	Workers       int    `json:"no_workers"` // Number of worker nodes
	Masters       int    `json:"no_masters"` // Number of Master nodes
	NetworkPlugin string `json:"network_plugin,omitempty"`
}

type ClusterStorage interface {
	InsertOrUpdateCluster(c *ClusterControllerMetadata) error
	Get(id string) *ClusterControllerMetadata
}
