# Nanobit backend engineer assignment

## Running
The easiest way to run the components is using docker-compose

```bash
docker-compose build
docker-compose up
```
# Kubernetes deployment
```
kubectl apply -f k8s/redis-deployment.yaml
kubectl apply -f k8s/redis-service.yaml
kubectl apply -f k8s/web-deployment.yaml
kubectl apply -f k8s/web-service.yaml
kubectl apply -f k8s/worker-deployment.yaml
```

# Nomad deployment
If nomad is installed on a remote machine set:
```bash
export NOMAD_ADDR="NOMAD REMOTE ADDRESS"
```

Validate nomad jobs by running:
```bash
nomad validate /deplyoments/nomad/redis/redis.nomad
```

Run the job with the command:
```bash
nomad run /deployments/nomad/redis/redis.nomad
```

## Usage
To connect to the websocket server use the following command:
```bash
wsta -I ws://127.0.0.1:8080
```

There are two types of messages that can be sent to the websocket server.

To list all of the users and their favourite numbers send:
```json
{"action": "listUsers"}
```

To set a user's favourite number send:
```json
{"action": "setNumber", "msg": {"username": "test-user", "number": 123}}
```
