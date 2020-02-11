# ntnx-hackathon-ayin
Nutanix hackathon Ayin 

# Build & Run curate-clusters-service

```
go build -o ./build/curate-clusters-service ./cmd/curate-clusters-service
./build/curate-clusters-service
```
or
```
docker-compose build
docker-compose up curate-clusters-service
```

# Register cluster

```
curl -X POST --data-binary '{"id": "f10868ed-5e17-4846-aca0-a15a0845dc5d", "name": "Test cluster", "no_workers": 1, "no_masters": 1}' http://localhost:9090/clusters/register/
```

# List clusters

```
curl http://localhost:9090/clusters/
```

# Run on-prem-agent

```
go build -o ./build/on-prem-agent ./cmd/on-prem-agent
CLUSTER_CONTROLLER_BASE_URL=http://localhost:9090 ./build/on-prem-agent
```
