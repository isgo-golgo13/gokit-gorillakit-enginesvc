# stage 1
FROM golang:1.20-alpine as stage

WORKDIR /gokit-enginesvc
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
# copy the source from the current directory to the Working Directory inside the container
#COPY . .
# The 'k8s' Kubernetes directory is NOT copied into the target image as it out of app code lifecycle
COPY client/ client/
COPY servicekit/ servicekit/
COPY cmd/service cmd/service/
COPY storagekit/ storagekit/

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
# the following line works if uncommented
#RUN go build k8s-reference-app/gokit-enginesvc/cmd/service
# the following line works
RUN go build -o service cmd/service/main.go


# stage 2
FROM alpine:latest
RUN apk add --no-cache git tzdata curl 
COPY --from=stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/
# copy ONLY the app ('service' exe) NOT the source from '/gokit-enginesvc'
COPY --from=stage /gokit-enginesvc/service .
# healthcheck
HEALTHCHECK CMD curl --fail http://localhost:8080/ || exit 1
CMD ["./service"]