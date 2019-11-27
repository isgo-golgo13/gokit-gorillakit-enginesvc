# Docker/K8s Engine Microservice

This example demonstrates how to use Go kit to implement a REST-y HTTP service.
It uses the excellent [gorilla mux package](https://github.com/gorilla/mux) for routing.

Build the code (not using Docker)

The module definition for the project is in go.mod as : `github.com/isgo-golgo13/enginesvc` and to build this at the root of the project (enginesvc/ from the git clone) issue:
`go build github.com/isgo-golgo13/enginesvc/cmd/enginesvc` and the exe will deposit in the root of the project.

Run the code (not using Docker)
`go run github.com/marvincaspar/go-example/cmd/server`


Build the docker image

```bash

docker build -t enginesvc-healthchk:1.0 .

```

Run the docker container
```bash

docker run  --name enginesvc-healthchk -p 8080:8080 engine-healthchk:1.0
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



Create/Register an Engine:

```bash
$ curl -d '{"id":"1234","factory_id":"utc_pw_10-0001", "engine_config" : "Radial", "engine_capacity": 660.10, "fuel_capacity": 400.00, "fuel_range": 240.60}' -H "Content-Type: application/json" -X POST http://localhost:8080/engines/
{}
```

Get the engine you just created

```bash
$ curl localhost:8080/engines/1234
{"engine":{"id":"1234","factory_id":"utc_pw_10-0001", "engine_config" : "Radial", "engine_capacity": 660.10, "fuel_capacity": 400.00, "fuel_range": 240.60}}
```