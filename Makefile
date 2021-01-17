SERVICE_NAME=service

.PHONY : clean 
.DEFAULT_GOAL : all

all:
	
build: 
	go build github.com/isgo-golgo13/go-gokit-gorilla-restsvc/cmd/service

run:
	./service

clean: 
	rm -f ${SERVICE_NAME}

docker-build:
	sh kill-docker.sh
	sh kr8-docker.sh 

docker-run:
	sh kr8-docker-run.sh

docker-clean:
	sh kill-docker.sh

test:
	go test -v