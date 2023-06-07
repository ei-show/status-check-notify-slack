FROM golang:1.19-alpine

WORKDIR /app

COPY src/ .

CMD ["go", "run", "/app/main.go"]