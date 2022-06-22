# Stage 1 (Build)
FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS builder

ARG VERSION=0.1
RUN apk add --update --no-cache git make upx
WORKDIR /app/
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o fusion \
    fusion.go
RUN upx fusion
RUN echo "ID=\"distroless\"" > /etc/os-release

# Stage 2 (Final)
FROM gcr.io/distroless/static:latest
COPY --from=builder /etc/os-release /etc/os-release

COPY --from=builder /app/fusion /usr/bin/
COPY ./fusion.yaml /etc/fusion/config.yml

EXPOSE 8899

CMD [ "/usr/bin/fusion", "--config", "/etc/fusion/config.yml" ]