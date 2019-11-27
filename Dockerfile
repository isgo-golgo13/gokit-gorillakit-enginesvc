# stage 1
FROM golang:1.13.0-stretch as stage
COPY . /cmd
WORKDIR /cmd
ENV GO111MODULE=on
#RUN CGO_ENABLED=0 GOOS=linux go build -o cmd/enginesvc
RUN CGO_ENABLED=0 GOOS=linux go build github.com/isgo-golgo13/enginesvc/cmd/enginesvc 

# stage 2
FROM alpine:latest
WORKDIR /root/
COPY --from=stage /cmd .
# healthcheck
HEALTHCHECK CMD curl --fail http://localhost:5000/ || exit 1
CMD ["./enginesvc"]
