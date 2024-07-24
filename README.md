# githubECS

This project consists of three services to monitor GitHub repositories and commits:
1. **Repo Discovery Service:** Fetches and saves metadata about repositories from a specified GitHub organization.
2. **Commit Monitor Service:** Checks for new commits in the saved repositories periodically.
3. **Commit Manager Service:** Manages the fetching of commits from a specified start date.

## Architecture

The project is structured into three services that communicate via RabbitMQ:

1. **Repo Discovery Service:** Fetches repository metadata from GitHub and saves it to the database.
2. **Commit Monitor Service:** Periodically checks for new commits in the repositories.
3. **Commit Manager Service:** Manages the fetching of commits based on a specified start date.

## Setup

### Prerequisites

- Docker
- Docker Compose

### Environment Variables

Create a `.env` file in the root of the project and add the following variables:

```
DATABASE_URL=postgres://username:password@localhost:5432/githubecs
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

### Building and Running the Services

## Locally
Within the cmd directory is the entry points from which the three services can be started.

## Using Docker

1. Build the Docker images:

```sh
docker-compose build
```

2. Run the services:

```sh
docker-compose up
```

This will start all services and the necessary RabbitMQ instance.

## Endpoints

### Repo Discovery Service

- **POST /search**
    - Fetches repositories from GitHub and saves them to the database.
    - **Query Parameters:**
        - `query` (string): The search query to filter repositories.

- **GET /repositories/:full_name**
    - Retrieves a repository by its full name.
    - **Path Parameters:**
        - `full_name` (string): The full name of the repository.

- **GET /search**
    - Searches repositories by language.
    - **Query Parameters:**
        - `language` (string): The language to filter repositories by.

- **GET /top**
    - Retrieves the top N repositories by stars count.
    - **Query Parameters:**
        - `n` (int): The number of repositories to retrieve.

- **POST /reset_start_date**
    - Resets the start date for fetching commits.
    - **Request Body:**
        - `start_date` (string): The new start date in RFC3339 format.

## Usage

### API Endpoints

- **Search Repositories:** `POST /search?query=<query>`
- **Get Repository by Full Name:** `GET /repositories/:full_name`
- **Get Commits for a Repository:** `GET /repositories/:full_name/commits`
- **Search Repositories by Language:** `GET /search?language=<language>`
- **Get Top N Repositories by Stars Count:** `GET /top?n=<number>`
- **Reset Start Date for Commits:** `POST /reset_start_date?start_date=<start_date>`
