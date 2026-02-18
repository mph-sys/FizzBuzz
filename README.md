# FizzBuzz Service

A Go-based REST API service that exposes a customizable FizzBuzz endpoint and tracks request statistics.

## Overview

This project is a CLI application built with [Cobra](https://github.com/spf13/cobra) that launches an HTTP server using the [Gin](https://github.com/gin-gonic/gin) framework. It allows users to generate FizzBuzz sequences with custom parameters and retrieves statistics about the most frequent requests.

## Installation & Usage

### Build

```bash
go build -o fizzbuzz-service
```

### Run HTTP Server

The application exposes a `http-server` command to start the REST API.
It requires the environment variables `MYSQL_USER` and `MYSQL_PASSWORD` to be set to connect to the database.

```bash
export MYSQL_USER=user
export MYSQL_PASSWORD=password
./fizzbuzz-service http-server --mysql-db dbname --mysql-host localhost
```

#### Flags

- `--mysql-db`, `-d` (string): MySQL DB name(required for persistence).
- `--mysql-host`, `-H` (string): MySQL host (default "localhost").
- `--bind-addr`, `-b` (string): Address to bind the server to (default ":8080").
- `--prometheus-bind-addr`, `-p` (string): Address to bind the prometheus metrics server to (default ":2112").

## Features

- **Customizable FizzBuzz**: Specify the two integers, the limit, and the two replacement strings.
- **Usage Statistics**: Tracks the number of hits for each request configuration and exposes the most used one.
- **Persistence**: Uses a SQL database to store request statistics.

## API Endpoints

### 1. Generate FizzBuzz Sequence

Returns a list of strings corresponding to the FizzBuzz sequence.

- **URL**: `/fizzbuzz/run`
- **Method**: `POST`
- **Query Parameters**:
    - `int1` (required): First multiple.
    - `int2` (required): Second multiple.
    - `limit` (required): The limit of the sequence (from 1 to limit).
    - `str1` (required): String to replace multiples of `int1`.
    - `str2` (required): String to replace multiples of `int2`.

**Example:**
```bash
curl -X POST "http://localhost:8080/fizzbuzz/run?int1=3&int2=5&limit=100&str1=fizz&str2=buzz"
```

### 2. Get Most Requested Stats

Returns the parameters used in the most frequent request.

- **URL**: `/fizzbuzz/stats/most-requested`
- **Method**: `GET`

**Example:**
```bash
curl "http://localhost:8080/fizzbuzz/stats/most-requested"
```

## Database Schema

The application requires a MySQL-compatible database with the following table structure (inferred from usage):

```sql
CREATE TABLE `stats` (
    `int1` INT,
    `int2` INT,
    `limit` INT,
    `str1` VARCHAR(255) COLLATE utf8mb4_unicode_ci,
    `str2` VARCHAR(255) COLLATE utf8mb4_unicode_ci,
    `hits` INT,
    PRIMARY KEY (`int1`,`int2`,`limit`,`str1`,`str2`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## Project Structure

The project follows a modular structure to separate concerns:

- **`cmd/`**: Application entry points. Contains the main executable and CLI command definitions (using Cobra).
- **`http/`**: HTTP layer implementation.
  - **`handlers/`**: Gin route handlers that process incoming requests.
  - **`models/`**: JSON request/response structures specific to the API.
  - **`service.go`**: Server configuration, routing setup, and startup logic.
- **`pkg/`**: Core business logic (Service layer).
  - **`models/`**: Domain models shared across the application.
  - Contains the pure logic for FizzBuzz generation and database interactions for statistics.
- **`api/`**: API documentation and specifications (OpenAPI).