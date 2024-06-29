# Go HTTP Server Project

This project demonstrates building a simple HTTP server in Go by gradually adding features in stages. Each stage introduces new functionalities, starting from binding to a port to supporting concurrent connections, reading request headers and bodies, serving files, and handling HTTP compression.

The full code for each stage is available in the `server.go` file, with comments indicating the different stages.

## Table of Contents

- [Stages](#stages)
  - [Stage 1: Bind to a Port](#stage-1-bind-to-a-port)
  - [Stage 2: Respond with 200 OK](#stage-2-respond-with-200-ok)
  - [Stage 3: Extract URL Path](#stage-3-extract-url-path)
  - [Stage 4: Respond with a Body](#stage-4-respond-with-a-body)
  - [Stage 5: Read Headers](#stage-5-read-headers)
  - [Stage 6: Support Concurrent Connections](#stage-6-support-concurrent-connections)
  - [Stage 7: Return a File](#stage-7-return-a-file)
  - [Stage 8: Read Request Body](#stage-8-read-request-body)
  - [HTTP Compression](#http-compression)
    - 1.[Compression Headers](#compression-headers)
    - 2.[Multiple Compression Schemes](#multiple-compression-schemes)
    - 3.[Gzip Compression](#gzip-compression)
- [Usage](#usage)
- [License](#license)

## Stages

### Stage 1: Bind to a Port

In this stage, the server is configured to bind to a specific port and handle any errors that occur if the binding fails. This is the foundation of our HTTP server.

### Stage 2: Respond with 200 OK

After successfully binding to a port, the server can accept a connection and respond with a simple `200 OK` status to acknowledge the request.

### Stage 3: Extract URL Path

In this stage, the server parses the incoming request to extract the URL path. This allows the server to understand which resource the client is requesting.

### Stage 4: Respond with a Body

The server now responds with a body, providing different content based on the URL path. This includes basic responses like "Hello, World!" or echoing a message back to the client.

### Stage 5: Read Headers

This stage introduces reading and parsing headers from the HTTP request. The server can now understand and respond based on various headers, such as `User-Agent`.

### Stage 6: Support Concurrent Connections

The server is enhanced to handle multiple concurrent connections using goroutines. This allows it to serve multiple clients simultaneously without blocking.

### Stage 7: Return a File

In this stage, the server is capable of returning a file in response to a request. This includes setting the appropriate headers and serving the file content.

### Stage 8: Read Request Body

The server can now read the body of incoming requests. This is essential for handling POST requests where the client sends data to the server.

### HTTP Compression

#### 1:Compression Headers

The server is updated to handle HTTP compression headers, allowing it to respond with compressed content if the client supports it.

#### 2:Multiple Compression Schemes

Support for multiple compression schemes is added, enabling the server to choose the best compression method based on client capabilities.

#### : 3:Gzip Compression

In this stage, the server is updated to support gzip compression for the responses. When a client sends a request with the `Accept-Encoding: gzip` header, the server compresses the response body using gzip before sending it back to the client. This reduces the size of the response and improves the performance for clients that support gzip compression.

## Usage

1. Clone the repository.
2. Navigate to the project directory.
3. Run the server using the command:

   ```sh
   go run server.go
