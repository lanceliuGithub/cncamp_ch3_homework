export tag=1.0

build:
	echo "Building HTTP Server Binary"
	mkdir -pv bin/linux/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64

release: build
	echo "Building HTTP Server Container Image"
	docker build -t lanceliu2022/myhttpserver:${tag} .

push: release
	echo "Pushing Local Container Image to Docker Hub"
	docker push lanceliu2022/myhttpserver:${tag}
