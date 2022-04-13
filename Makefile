.PHONY: build

TAG=$(shell git describe --abbrev=0 --tags)
DATE=$(shell go run ./scripts/date.go)

build:
		@go mod tidy -compat=1.8 && \
		go build -ldflags "-X main.version=$(TAG) -X main.buildDate=$(DATE)" -o doko

install: doko
		@mv doko /usr/local/bin

brc: # build doko container
		@docker build -t dokocli/doko . && \
		docker push dokocli/doko

bfrc: # build full doko container
		@docker build -t dokocli/doko-full --file ./docker/doko-full/Dockerfile . && \
		docker push dokocli/doko-full

brcwc: # build doko container with cache
		@docker pull dokocli/doko:latest && \
		docker build -t dokocli/doko --cache-from dokocli/doko:latest . && \
		docker push dokocli/doko

bfrcwc: # build full doko container with cache
		@docker pull dokocli/doko-full:latest && \
		docker build -t dokocli/doko-full --cache-from dokocli/doko-full:latest . && \
		docker push dokocli/doko-full
