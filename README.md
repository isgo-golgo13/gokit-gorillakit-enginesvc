# Go, Go-kit, Gorilla (Go) Web Toolkit REST Service using Docker

It uses the [gorilla mux package](https://github.com/gorilla/mux) for routing.

### Build the app  (not using Docker)

The module definition for the project is in go.mod as : `github.com/isgo-golgo13/go-gokit-gorilla-restsvc` and to build this at the root of the project (github.com/isgo-golgo13/go-gokit-gorilla-restsvc/ from the git clone) issue:
`go build github.com/isgo-golgo13/go-gokit-gorilla-restsvc/cmd/service` or `go build -o service cmd/service/main.go` (the latter is used in the Dockerfile) and the exe will deposit in the root of the project or just issue the provided Makdefile as follows:

`make compile`

### Run the app (not using Docker)
`go run github.com/isgo-golgo13/go-gokit-gorilla-restsvc/cmd/service`

### Clean the built binary artifact
Issue the following: `make clean`


### Build the Docker Container Image

```bash

sh kr8-docker.sh 

```
This shell script executes the following:

```
#! /bin/sh

set -ex

docker image build -t go-gokit-gorilla-restsvc:1.0 .
```

### Run the Docker Container
```bash

sh kr8-docker-run.sh
```


### Destroy the Docker Image and Container
```

sh kill-docker.sh
```


Test the HEALTHCHECK endpoint of the Docker container (order prior or after hitting service endpoints isn't affecting). The goal of the container is to serve traffic on port 8080 and our healthchk should ensure that is occurring. The default options of the HEALTHCHECK flag are interval 30s, timeout 30s, start-period 0s, and retries 3. If different options are reqiured, look to the Docker HEALTHCHECK option page at docker.io.
```
docker inspect --format='{{json .State.Health}}' enginesvc-healthchk
```

If run right at the start after the container starts you see the Status as `starting`
```
{"Status":"starting","FailingStreak":0,"Log":[]}
```

And after the health check runs (after the default interval of 30s):
```
{"Status":"healthy","FailingStreak":0,"Log":[{"Start":"2017-07-21T06:10:51.809087707Z","End":"2017-07-21T06:10:51.868940223Z","ExitCode":0,"Output":"output is specific to the curl POST "}]}
```
It is possible that container can accept and process requests properly and yet you will see the Status as `'unhealthy'` and this likely issue from a root image in the Docker container.



### Create/Register an Engine:

```bash
$ curl -d '{"id":"00001","factory_id":"utc_pw_10-0001", "engine_config" : "Radial", "engine_capacity": 660.10, "fuel_capacity": 400.00, "fuel_range": 240.60}' -H "Content-Type: application/json" -X POST http://localhost:8080/engines/
{}
```

### Retrieve an Engine
 
```bash
$ curl localhost:8080/engines/00001
{"engine":{"id":"00001","factory_id":"utc_pw_10-0001", "engine_config" : "Radial", "engine_capacity": 660.10, "fuel_capacity": 400.00, "fuel_range": 240.60}}
```

### To Exec into Docker Container (In Linux Shell Terminal of Container)

```
Syntax:

docker container exec -it <container-name> <shell> 

Actual Use:

docker container exec -it go-gokit-gorilla-restsvc /bin/sh
```

The run `ls -al` to see the directory structure of the application as laid out in the container.


### To Push Docker Container Image to  DockerHub

```
1. docker login 
2. docker image tag the container image 
3. docker image push to DockerHub container image repo

In Detail:

Step #1. 
docker login

Results in:

Username: isgogolgo13
Password: <dockerhub password>


Step #2 
docker image tag the container image 

Do 'docker image ls' to get the Docker container image (including container image tag)
to tag as pre-step required to push the tagged container image to DockerHub repo.

To do this:

Syntax:

docker tag <container-image-id (see 'IMAGE ID')> <dockerhub-username>/<container-image-name>:<tag>

Actual Use:
docker tag ############ isgogolgo13/go-gokit-gorilla-restsvc:1.0


Step #3 
docker image push to DockerHub container image repo 

Syntax:

docker push <dockerhub-username/<container-image-tag>

Actual Use:

docker push isgogolgo13/go-gokit-gorilla-restsvc:1.0

```


### Running the K8s (K8s YAML Approach, K8s Helm Approach)

To run the `sgogolgo13/go-gokit-gorilla-restsvc` on K8s cluster we will use Rancher `k3d` using Traefik 2.2 LoadBalancer/Ingress Controller as of 1-25-2020
k3d included Traefik uses v1 and to use v2 the step to include v2 is to provide arguments to the `k3d cluster create` command.

**1.** Create the K3D cluster (default 1 server node, 1 agent node) - Disable default Traefik v1 Load Balancer/Ingress
```
k3d cluster create go-gokit-gorilla-restsvc-cluster \
--api-port 127.0.0.1:6443 \
-p 80:80@loadbalancer \
-p 443:443@loadbalancer \
--k3s-server-arg "--no-deploy=traefik"
```

The following flags and flag arguments are as follows:

| Flag                    |                        | Result                                                                      | 
|:-----------------------:|:----------------------:|:--------------------------------------------------------------------------- | 
| --api-port              |  127.0.0.1:6443        | k3d cluster will serve at IP: 127.0.0.1 (localhost) on port 6443            | 
| -p                      | 80:80@loadbalancer     | associate localhost port 80 to port 80 to k3d virtual load balancer         | 
| -p                      | 443:443@loadbalancer   | associate localhost port 443 to port 443 to k3d virtual load balancer       | 
| --k3s-server-arg        | "--no-deploy=traefik"  | disable Traefik load balancer/ingress v1 to override with v2 (using Helm) * |



**2.** Install Traefik Load Balancer/Ingress Controller v2

```
helm repo list  # Check if traefik helm repo is downloaded/cached on the host (if NOT the helm repo add will do it in the next step)
helm repo add traefik https://containous.github.io/traefik-helm-chart  # Do only if traefik helm  repo not downloaded/cached on the host)

helm install traefik traefik/traefik
```