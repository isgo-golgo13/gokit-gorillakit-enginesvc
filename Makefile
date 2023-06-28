SERVICE_NAME=service

.PHONY : clean 
.DEFAULT_GOAL : all

all:
	
compile: 
	GO111MODULE=on go build gokit-gorillakit-enginesvc/cmd/service

run:
	./service

clean: 
	rm -f ${SERVICE_NAME}

test:
	go test -v