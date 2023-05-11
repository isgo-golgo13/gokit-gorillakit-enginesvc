SERVICE_NAME=service

.PHONY : clean 
.DEFAULT_GOAL : all

all:
	
compile: 
	go build github-actions-ci-pipeline/cmd/service

run:
	./service

clean: 
	rm -f ${SERVICE_NAME}

test:
	go test -v