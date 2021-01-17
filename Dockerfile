# # stage 1
# FROM golang:1.15-alpine as stage
# COPY . /cmd
# WORKDIR /cmd
# ENV GO111MODULE=on

# RUN CGO_ENABLED=0 GOOS=linux go build github.com/isgo-golgo13/go-gokit-gorilla-restsvc/cmd/service 

# # stage 2
# FROM alpine:latest
# # healthcheck req to run curl
# RUN apk --no-cache add curl
# WORKDIR /root/
# COPY --from=stage /cmd .
# EXPOSE 8080
# # healthcheck
# HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
#   CMD curl -f http://localhost:8080/health || exit 1
# CMD ["./service"]


# stage 1
FROM golang:1.15-alpine as stage

WORKDIR /go-gokit-gorilla-restsvc
COPY go.mod go.sum ./
RUN go mod download
# copy the source from the current directory to the Working Directory inside the container
#COPY . .
COPY ./client ./cmd ./servicekit ./



ENV GO111MODULE=auto
RUN CGO_ENABLED=0 GOOS=linux go build github.com/isgo-golgo13/go-gokit-gorilla-restsvc/cmd/service

# stage 2
FROM alpine:latest
RUN apk add --no-cache tzdata curl 
COPY --from=stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/
COPY --from=stage /go-gokit-gorilla-restsvc .
# healthcheck
HEALTHCHECK CMD curl --fail http://localhost:8080/ || exit 1
CMD ["./service"]