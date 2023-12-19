FROM golang:1.21-alpine

ENV GIN_MODE=release
ENV PORT=8080
ENV DB_PATH=/app-storage/links.db

RUN apk add build-base

WORKDIR /app-src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /app
RUN mkdir /app-storage

EXPOSE 8080
CMD [ "/app" ]
