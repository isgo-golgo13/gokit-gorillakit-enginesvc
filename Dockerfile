# stage 1
FROM golang:1.15.0-stretch as stage
# COPY . /cmd
WORKDIR /cmd
COPY go.mod .
COPY go.sum .
COPY vendor .
ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go build github.com/isgo-golgo13/go-gokit-gorilla-restsvc/cmd/enginesvc 

# stage 2
FROM alpine:latest
WORKDIR /root/
COPY --from=stage /cmd .
# healthcheck
HEALTHCHECK CMD curl --fail http://localhost:5000/ || exit 1
CMD ["./enginesvc"]
