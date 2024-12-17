# Server Info API

A simple Go web server that returns a JSON object with the client's IP address and hostname information.

## Features

- Retrieves the client's IP address from HTTP request headers or remote address.
- Returns the information in a formatted JSON response.
- Configurable server port via an environment variable.

## Prerequisites

- Go 1.17 or later installed on your system.

## Installation

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/networkrehab/server-info-api.git
cd server-info-api
```

Initialize the Go module and download dependencies:

```bash
go mod init github.com/networkrehab/server-info-api
go mod tidy
```

## Usage

Run the server using the `go run` command:

```bash
go run main.go
```

Alternatively, build the executable and run it:

```bash
go build -o server-info-api
./server-info-api
```

The server listens on port `8081` by default. To specify a different port, set the `PORT` environment variable:

```bash
export PORT=8080
go run main.go
```

## Example

Make a request to the server to retrieve your IP and hostname:

```bash
curl http://localhost:8081/
```

Sample JSON response:

```json
{
  "IP": "127.0.0.1",
  "Hostname": "localhost"
}
```

## Contributing

Contributions are welcome. Please open an issue or submit a pull request for any improvements.

## License

This project is licensed under the MIT License.
