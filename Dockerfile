# stage 1
FROM golang:1.15-alpine as stage
COPY . /cmd
WORKDIR /cmd
ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go build github.com/isgo-golgo13/go-gokit-gorilla-restsvc/cmd/enginesvc 
RUN ls -al

# stage 2
FROM alpine:latest
# healthcheck req to run curl
RUN apk --no-cache add curl
WORKDIR /root/
COPY --from=stage /cmd .
# healthcheck
# HEALTHCHECK --interval=30s --timeout=3s \
#   CMD curl -f http://localhost:8080/health || exit 1
CMD ["./enginesvc"]
