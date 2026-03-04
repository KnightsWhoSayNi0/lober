FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . /app/
RUN --mount=type=cache,target="/root/.cache/go-build" \
	go build -o server main.go

FROM alpine:3.22
RUN apk --update add ca-certificates && \
	rm -rf /var/cache/apk/*
RUN adduser -D lober
USER lober
COPY --from=builder /app /home/lober/app
WORKDIR /home/lober/app
EXPOSE 8000
EXPOSE 8443
CMD ["./server"]
