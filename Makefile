GIT_HEAD = $(shell git rev-parse HEAD | head -c8)

docker: 
	docker build -t ghcr.io/luminous-gsm/fusion:local .

build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o build/fusion_linux_amd64 -v fusion.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o build/fusion_linux_arm64 -v fusion.go

run-docker-luminous-pov:
	docker run --name fusion-luminous-pov -e FUSION_UNIQUE_ID=UniqueId -e FUSION_API_SECRET_TOKEN=SuperSecretToken -e FUSION_CONSOLE_LOCATION=http://localhost:3200 ghcr.io/luminous-gsm/fusion:local

certificate:
	openssl req -x509 -newkey rsa:4096 -days 3650 -nodes -keyout ./certs/fusion.key -out ./certs/fusion.crt -subj "/CN=luminous-gsm.com" -addext "subjectAltName=DNS:luminous-gsm.com"

.PHONY: all docker