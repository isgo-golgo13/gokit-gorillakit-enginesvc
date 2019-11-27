SERVICE_NAME=enginesvc

.PHONY : clean 
.DEFAULT_GOAL : all

all:
	
compile: 
	go build github.com/isgo-golgo13/enginesvc/cmd/enginesvc

clean: 
	rm -f ${SERVICE_NAME}

test:
	go test -v