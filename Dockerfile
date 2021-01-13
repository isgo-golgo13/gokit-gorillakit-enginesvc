# stage 1
FROM golang:1.15-alpine as stage
COPY . /cmd
WORKDIR /cmd
ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go build github.com/isgo-golgo13/go-gokit-gorilla-restsvc/cmd/service 

# stage 2
FROM alpine:latest
# healthcheck req to run curl
RUN apk --no-cache add curl
WORKDIR /root/
COPY --from=stage /cmd .
EXPOSE 8080
# healthcheck
HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
CMD ["./service"]
