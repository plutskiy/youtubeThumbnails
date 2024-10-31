# YouTube Thumbnail Downloader

![Tests](https://github.com/plutskiy/youtubeThumbnails/workflows/Tests/badge.svg)

This project provides a gRPC service for downloading YouTube video thumbnails based on video URLs. The server fetches thumbnails from YouTube, caches them in memory, and stores them in a local SQLite database.

## Features

- **gRPC Service**: Allows clients to request thumbnails by sending YouTube URLs.
- **Caching**: Stores thumbnails in memory for quick access.
- **Database Storage**: Saves thumbnails to an SQLite database to reduce duplicate requests to YouTube.

## Prerequisites

- Go 1.19 or higher
- SQLite3 (for database storage)

## Getting Started

1. **Install Dependencies**: Run 
```bash
go mod tidy
```

to install the required dependencies.

2. **Generate gRPC Code**: Before running the server, generate the gRPC code from the Protobuf file:
   ```bash
   protoc --go_out=. --go-grpc_out=. api/thumbnail.proto
   ```

3. **Run**: Pass YouTube video URLs as command-line arguments to download thumbnails.

   **Example**:
   ```bash
   go run main.go "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
   ```

   - Use the `-async` flag to download thumbnails asynchronously:
     ```bash
     go run main.go -async "https://www.youtube.com/watch?v=rWJ1tPCnVJI" "https://www.youtube.com/watch?v=rWJ1tPCnVJI"
     ```

5. **Thumbnail Storage**: Thumbnails are saved in an `image/` directory as `.jpg` files.

## Project Structure

- `api/`: Protobuf files defining the gRPC service.
- `cmd/`: Contains client and server implementations.
- `pkg/`: Contains application logic, including gRPC handlers.

## Database Setup

The project uses SQLite for storing thumbnail data. When first run, the server will create a table for storing thumbnails if it doesn’t already exist.

## Testing

To run tests, use:
```bash
go test ./...
```
