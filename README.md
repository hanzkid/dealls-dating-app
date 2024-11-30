# Dating Apps Project

## Project Structure

- `db/`: Contains database-related files: seeders and migrations.
- `be/`: Contains main backend service code.
- `docker-compose.yml`: Docker Compose file to set up and run the services.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Clone the repository:
    ```sh
    https://github.com/hanzkid/dealls-dating-app.git
    cd dealls-dating-app
    ```

2. Build and start the services:
    ```sh
    docker-compose up --build
    ```

3. Access the application:
    - Backend API: `http://localhost:<backend_port>`
    - Database: `http://localhost:<database_port>`

## Services

### Backend

The backend service is defined in the `be/` directory. It handles the core logic of the dating application.

### Database

The database service is defined in the `db/` directory. It stores all the data required by the application.

## Configuration

The `docker-compose.yml` file contains the configuration for all services. You can modify this file to change the settings as per your requirements.
