# Kafka Go usage

## Contents

- [About](#about-project)
- [Kafka UI](#kafka-ui)
- [Installation](#installation)
- [Usage](#usage)

## About project

This project is a simple example of using the Kafka message broker with Golang applications. This example implements a cluster of two Kafka servers - they have a common topic and 2 partitions that are shared between the servers, in addition one server stores replications of the other and can completely replace it in case of failure. Two consumers have been implemented on Go, as well as a produser with an HTTP server responsible for sending messages to Kafka.

## Kafka UI

To use the Kafka UI after running the service via Docker, go to `http://localhost:8090`

## Installation
> IMPORTANT If you have made changes to the configuration files, check them against the parameters specified in docker-compose.yml

To start a project using Docker, run the following command from the root of the project:
```bash
docker-compose up -d
```
To stop, use:
```bash
docker-compose down
```

Use the following commands to start the producer and consumers:
#### Run Producer
```bash
go run cmd/producer/main.go
```

#### Start the first consumer
```bash
go run cmd/consumer/main.go
```
#### Start the second consumer
```bash
go run cmd/consumer2/main.go
```

## Usage
#### Request to send messages:
`POST /message`

#### Body
```json
{
    "message": "Hello World!"
}
```

### ⭐️ If you like my project, don't spare your stars 🙃