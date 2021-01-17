# stage 1
FROM golang:1.15-alpine as stage

WORKDIR /go-gokit-gorilla-restsvc
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
# copy the source from the current directory to the Working Directory inside the container
#COPY . .
# The 'k8s' Kubernetes directory is NOT copied into the target image as it out of app code lifecycle
COPY client/ client/
COPY servicekit/ servicekit/
COPY cmd/service cmd/service/

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
# the following line works if uncommented
#RUN go build github.com/isgo-golgo13/go-gokit-gorilla-restsvc/cmd/service
# the following line works
RUN go build -o service cmd/service/main.go


# stage 2
FROM alpine:latest
RUN apk add --no-cache git tzdata curl 
COPY --from=stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/
# copy ONLY the app ('service' exe) NOT the source from '/go-gokit-gorilla-restsvc'
COPY --from=stage /go-gokit-gorilla-restsvc/service .
# healthcheck
HEALTHCHECK CMD curl --fail http://localhost:8080/ || exit 1
CMD ["./service"]