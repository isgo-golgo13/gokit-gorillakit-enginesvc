# Go, Go-kit, Gorilla (Gorilla mux router package) REST Service using Docker

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

### Retrive an Engine
 
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