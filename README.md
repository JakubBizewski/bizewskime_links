# Short Link Generator

This is a short link generator written in Go.  
It provides an API endpoint for creating new short links and an endpoint to use the shot link.

This project is a learning exercise for me. This is my first Go application.  
It is hosted on [l.bizewski.me](https://l.bizewski.me).

## Technical and business considerations
- This application is build using the hexagonal architecture. I know it is a bit overkill for this application, but I wanted to learn more about it.
- Short codes are built using 3 random characters.
- If generated short code already exists, it will generate a new one. It will try 10 times before giving up.
- The application uses SQLite as a database. There is no need for a full-blown database for this application.

## TODO:
* [x] Implement a database repository. Right now, it is using a memory repository,
* [ ] Implement history of short links,
* [ ] Add rate limiting to prevent DoS attacks and brute-forcing.

## Deployment
```sh
docker run -d --restart unless-stopped --name links -p 8080:8080 -v /var/links:/app-storage -e ENCRYPTION_KEY=[your 32byte key] ghcr.io/jakubbizewski/jakubme_links:master
```

## Endpoints

### Create Short Link

```
POST /new
{
  "targetUrl": "https://www.google.com"
}

Response:
200 OK
{
  "shortLink": "abc123"
}


400 Bad Request
{
  "error": "Invalid URL"
}
```

### Use Short Link

```
GET /abc123

Response:
302 Found
Location: https://www.google.com
```

## Project Structure

- `adapters/`: Contains the web interface and memory repository implementations.
- `domain/`: Contains the domain model and ports for the application.
- `mocks/`: Contains mock implementations for testing.
- `main.go`: The entry point for the application.

## Running the Application

To run the application, use the following command:

```sh
go run main.go
```

#### Testing
```sh
go test ./...
```
