GIT_HEAD = $(shell git rev-parse HEAD | head -c8)

docker: 
	docker build -t ghcr.io/luminous-gsm/fusion:local .

build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o build/wings_linux_amd64 -v wings.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o build/wings_linux_arm64 -v wings.go

run-docker-luminous-pov:
	docker run --name fusion-luminous-pov -e FUSION_UNIQUE_ID=UniqueId -e FUSION_API_SECRET_TOKEN=SuperSecretToken -e FUSION_CONSOLE_LOCATION=http://localhost:3200 ghcr.io/luminous-gsm/fusion:local

.PHONY: all docker