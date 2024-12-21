# Server Info API

A simple Go web server that returns a JSON object with the client's IP address and hostname information. This is intended to be an example project to act like Azure's metadata service, but with a focus on IP address and hostname information. It is not supposed to be feature complete. That is why it is a simple project so that you have a good base to add your own features.

## Features

- Retrieves the client's IP address from HTTP request headers or remote address
- Performs reverse DNS lookups for IP addresses
- Stores IP to hostname mappings in SQLite database
- Supports CSV import of IP/hostname mappings
- Configurable server port via environment variable
- Available as multi-architecture Docker image (AMD64 and ARM64)

## Prerequisites

Go 1.17 or later installed on your system at a minimum. The docker image is built with Go 1.23.4.  

You can download Go from the [official website](https://golang.org/dl/).

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

## Docker Usage

You can run the service using Docker in two ways:

### Building the Docker Image Locally

Build and run the Docker image locally:

```bash
# Build the image
docker build -t server-info-api .

# Run the container with database persistence
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  server-info-api
```

### Using Pre-built Image from GitHub Container Registry

Pull and run the latest image from GitHub Container Registry:

```bash
# Pull the image (automatically selects correct architecture)
docker pull ghcr.io/networkrehab/server-info-api:main

# Run the container with database persistence
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  ghcr.io/networkrehab/server-info-api:main
```

### Docker Environment Variables

You can override the default port using the PORT environment variable:

```bash
# Override the default port
docker run -d \
  -p 8081:8081 \
  -e PORT=8081 \
  -v $(pwd)/data:/app/data \
  server-info-api
```

### Using with CSV Data

To import a CSV file with Docker:

```bash
# Mount both the CSV file and database directory
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data.csv:/app/data.csv \
  -v $(pwd)/data:/app/data \
  server-info-api \
  /app/server-info-api -import /app/data.csv
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

The server listens on port `8080` by default. To specify a different port, set the `PORT` environment variable:

```bash
export PORT=8081
go run main.go
```

### Importing Hostname Mappings from a CSV File

You can import IP and hostname mappings into the SQLite database from a CSV file using the `-import` flag. The CSV file should contain IP addresses and hostnames in the following format:

```csv
<IP address>,<hostname>
```

Each line represents a mapping between an IP address and a hostname.

#### Example CSV File

```csv
127.0.0.1,localhost
192.168.1.2,host2
10.0.0.1,host3.example.com
```

#### Importing the CSV File

To import the CSV file into the database, run:

```bash
go run main.go -import data.csv
```

This command will read the CSV file `data.csv` and populate the `host_ip_map.db` SQLite database with the IP-hostname mappings.

#### Running the Server After Import

After importing the CSV data, start the server normally:

```bash
go run main.go
```

Now, when a request is made to the server, it will use the hostname mappings from the database.

## Example

Make a request to the server:

```bash
curl http://localhost:8080/
```

Sample JSON response:

```json
{
  "ip": "192.168.1.100",
  "hostname": "myhost.local",
  "reverseiplookup": "myhost.mydomain.com",
  "fqdn2ip": "192.168.1.100"
}
```

## Database Persistence

The SQLite database is stored in `/app/data` within the container. To persist the database between container restarts:

1. Create a local directory for the database:
```bash
mkdir -p ./data
```

2. Mount the directory when running the container:
```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  server-info-api
```

## Architecture Support

The Docker image is available for multiple architectures:
- AMD64 (x86_64)
- ARM64 (aarch64)

Docker will automatically pull the correct image for your system architecture.

## Contributing

Contributions are welcome. Please open an issue or submit a pull request for any improvements.

## License

This project is licensed under the MIT License.
