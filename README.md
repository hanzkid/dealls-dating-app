# Dating Apps Project

This project is a simplified dating application, designed as a technical assessment for the Senior Backend Engineer position at Dealls.

`Live url` : https://dating-apps-api.burhanyusuf.dev/
`API Documentation` : https://documenter.getpostman.com/view/40114065/2sAYBYepQk

## Project Structure

- `db/`: Houses database-related files, including seeders and migrations.
- `be/`: Contains the core backend service code.
- `docker-compose.yml`: Docker Compose file to set up and run the services.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Clone the repository:
    ```sh
    git clone https://github.com/hanzkid/dealls-dating-app.git
    cd dealls-dating-app
    ```

2. Build and launch the services:
    ```sh
    docker-compose up --build
    ```

3. Configure the Environment:

    Adjust the necessary environment variables in the `docker-compose.yml` file for the services. Additionally, copy the `.env.example` file in the `be` directory to `.env` and modify the values as needed.

4. Access the application:
    - Backend API: `http://localhost:7000`

5. Running unit test : 
    ```sh
    go test ./...
    ```

## Services

### Backend

The backend service is located in the `be/` directory and includes the following subfolders:

- `config/`: Contains configuration files for the application.
- `entity/`: Defines the data models used within the application.
- `helpers/`: Contains utility functions used throughout the application.
- `http/`: Manages HTTP server requests and processes.
- `repository/`: Contains code for database interactions.

#### Stack:

- Echo Framework: A high-performance, extensible, minimalist web framework for Go.
- Gorm: An ORM library for Golang, providing powerful tools for database interactions.
- Gomock: A mocking framework for Go, used for creating unit tests.
- Golang-jwt: A library for handling JSON Web Tokens (JWT) in Go applications.
- Godotenv: A library for loading environment variables from `.env` files.
- Golint: A tool for checking Go source code for style mistakes.

### Database

The database service is located in the `db/` directory and stores all the data required by the application.

#### Stack:

- PostgreSQL: A popular open-source database management system.
- Goose: A Golang-based tool for managing database migrations.

### Containerization

- Docker Compose: Used to orchestrate multiple containers.

#### Container:

- `app`: The container running the Go backend service.
- `postgres`: The container for the PostgreSQL database.
- `migrations`: The container executing migration scripts to apply database changes.
