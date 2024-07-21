# githubECS

This project is a Golang service that discovers repositories matching user interests, saves them to a database, and monitors commits to these repositories.

## Features

- Discover GitHub repositories based on user interests.
- Save discovered repositories to a PostgreSQL database.
- Monitor and fetch new commits for saved repositories.
- RESTful API endpoints to retrieve repository and commit information.

## Installation

### Prerequisites

- Docker
- Docker Compose

### Environment Variables

Create a `.env` file in the project root with the following variables:

```
DATABASE_URL=your_database_url
```

### Build and Run

1. Clone the repository:

```sh
git clone https://github.com/yemmyharry/githubECS.git
cd githubECS
```

2. Build and run the Docker containers:

```sh
docker-compose up --build
```

## Usage

### API Endpoints

- **Search Repositories:** `POST /search?query=<query>`
- **Get Repository by Full Name:** `GET /repositories/:full_name`
- **Get Commits for a Repository:** `GET /repositories/:full_name/commits`
- **Search Repositories by Language:** `GET /search?language=<language>`
- **Get Top N Repositories by Stars Count:** `GET /top?n=<number>`
- **Reset Start Date for Commits:** `POST /reset_start_date?start_date=<start_date>`

### Example Requests

- **Search Repositories:**

```sh
curl -X POST "http://localhost:8080/search?query=rust"
```

- **Get Repository by Full Name:**

```sh
curl -X GET "http://localhost:8080/repositories/test-repo"
```

- **Get Commits for a Repository:**

```sh
curl -X GET "http://localhost:8080/repositories/test-repo/commits"
```

- **Search Repositories by Language:**

```sh
curl -X GET "http://localhost:8080/search?language=go"
```

- **Get Top N Repositories by Stars Count:**

```sh
curl -X GET "http://localhost:8080/top?n=5"
```

- **Reset Start Date for Commits:**

```sh
curl -X POST "http://localhost:8080/reset_start_date?start_date=2019-01-01T00:00:00Z"
```

## Development

### Running Tests

To run the tests, use the following command:

```sh
go test ./...
```
