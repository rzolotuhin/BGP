GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/ci-src/generator/bin
GOARCH=amd64 arm64
DOCKER_HUB_REPO=rzolotuhin

MAKEFLAGS+="--silent"

download:
	go mod download

build compile: clean download
	@mkdir -p ${GOBIN}
	for arc in ${GOARCH}; do \
		echo " - building for architecture $${arc} "; \
		env GOOS=linux GOARCH=$${arc} go build -o ${GOBIN}/generator_$${arc}_linux; \
	done

run:
	go run main.go

clean:
	go clean
	rm -rf ${GOBIN}/*

docker: compile
	cd ci-src/ && \
	docker compose rm --stop --force && \
	docker compose build --no-cache && \
	docker compose up -d

docker-multi-arc:
	docker buildx create --name=multi-arc --node=multi-arc --platform=linux/amd64,linux/arm64
	docker buildx build --no-cache --builder=multi-arc --platform=linux/amd64,linux/arm64 --push --tag ${DOCKER_HUB_REPO}/bird_ru_subnet_generator ci-src/generator/
	docker buildx build --no-cache --builder=multi-arc --platform=linux/amd64,linux/arm64 --push --tag ${DOCKER_HUB_REPO}/bird ci-src/bird/