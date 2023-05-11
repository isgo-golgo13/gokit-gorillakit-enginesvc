SERVICE_NAME=service

.PHONY : clean 
.DEFAULT_GOAL : all

all:
	
compile: 
	go build gokit-gorillakit-enginesvc/cmd/service

run:
	./service

clean: 
	rm -f ${SERVICE_NAME}

test:
	go test -v