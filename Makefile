GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/ci-src/generator/bin
GOARCH=amd64

MAKEFLAGS+="--silent"

download:
	go mod download

build compile: clean download
	@mkdir -p ${GOBIN}
	env GOOS=linux GOARCH=${GOARCH} go build -o ${GOBIN}/generator_${GOARCH}_linux

run:
	go run main.go

clean:
	go clean
	rm -rf ${GOBIN}/*

docker: compile
	cd ci-src/ && \
	docker compose rm --stop --force && \
	docker compose build && \
	docker compose up -d
