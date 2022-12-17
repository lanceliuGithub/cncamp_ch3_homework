export tag=1.0

build:
	echo "Building HTTP Server Binary"
	mkdir -pv bin/linux/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64
