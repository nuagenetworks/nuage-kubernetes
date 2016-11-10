all: build 

test: build
	go test -v ./...

build: dep fmt lint
	go build -v
dep:
	go get ./...
fmt: 
	go fmt ./...

lint:
	cd api; go install; cd ..
	cd ovsdb; go install; cd ..
	go install
	gometalinter --disable=dupl --disable=gocyclo --deadline 300s ./... 	
