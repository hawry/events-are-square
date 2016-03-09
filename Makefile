all: 64 32

64:
	echo **************** BUILDING 64-BIT EXECUTABLE ****************
	go build -v -o events
	echo *************************** DONE ***************************

32:
	echo **************** BUILDING 32-BIT EXECUTABLE ****************
	export GOARCH=386; \
	go build -v -o events32
	echo *************************** DONE ***************************
