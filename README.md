# ntnx-hackathon-ayin
Nutanix hackathon Ayin 

# Register cluster

curl -X POST --data-binary '{"id": "f10868ed-5e17-4846-aca0-a15a0845dc5d", "name": "Test cluster", "no_workers": 1, "no_masters": 1}' http://localhost:9090/clusters/register/