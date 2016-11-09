OUT=eas
LDFLAGS = -ldflags "-s -w -X main.buildVersion=`git describe --dirty`"

.PHONY: all
.SILENT:

all: 64

64:
	echo **************** BUILDING 64-BIT EXECUTABLE ****************
	go build -v $(LDFLAGS) -o $(OUT)
	echo *************************** DONE ***************************

32:
	echo **************** BUILDING 32-BIT EXECUTABLE ****************
	export GOARCH=386; \
	go build -v $(LDFLAGS) -o $(OUT)32
	echo *************************** DONE ***************************

run: 64
	./$(OUT)

run32: 32
	./$(OUT)32

install:
	ln -s `pwd`/$(OUT) /usr/local/bin/$(OUT)
