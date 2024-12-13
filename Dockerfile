FROM golang:1.21-alpine as builder

ENV GIN_MODE=release

RUN apk add build-base

WORKDIR /app-src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /app

FROM alpine:3.21 as runner

ENV PORT=8080
ENV DB_PATH=/app-storage/links.db

COPY --from=builder /app /app
RUN mkdir /app-storage

EXPOSE 8080
CMD [ "/app" ]
