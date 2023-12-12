FROM golang:1.21-alpine

ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /app-src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app

EXPOSE 8080
CMD [ "/app" ]
