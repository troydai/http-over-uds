.PHONY: build
build:
	@ mkdir -p artifacts/bin
	@ go build -o artifacts/bin/server cmd/server/main.go 
	@ go build -o artifacts/bin/benchmark cmd/benchmark/main.go 

.PHONY: local-server
local-server:
	@ go run cmd/server/main.go

.PHONY: images
images: server-image client-image

.PHONY: server-image
server-image:
	@ docker build . -t http-over-uds-server:dev --target server

.PHONY: client-image
client-image:
	@ docker build . -t http-over-uds-client:dev --target client

.PHONY: benchmark
benchmark: clearn images
	@ docker-compose up \
		--abort-on-container-exit \
		--remove-orphans
	@ docker-compose down --remove-orphans -v

.PHONY: up
up:
	@ docker-compose up \
		--remove-orphans

.PHONY: clearn
clean:
	@ docker-compose down --remove-orphans -v
