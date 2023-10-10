# Go Notes API

This is a simple RESTful API for managing notes built with Go and the Gin framework.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [TODO](#todo)
- [Contributing](#contributing)
- [License](#license)

## Features

- Create, read, update, and delete (CRUD) notes.
- Rate limiting for API requests.
- Validation for note inputs.
- Cross-Origin Resource Sharing (CORS) support.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (v1.16 or higher) installed on your machine.
- MySQL database server running locally or at the specified host.

## Getting Started

To get started with this project, follow these steps:

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/solomonbaez/SB-Go-NAPI
   ```

2. Change into the project directory:

   ```bash
   cd SB-Go-NAPI
   ```

3. Install project dependencies:

   ```bash
   go mod tidy
   ```

4. Set up your MySQL database and update the database connection configuration in the `api/config/cfg.go` file:

   ```go
   const (
       DBUSER     = "your-username"
       DBPASSWORD = "your-password"
       DBNET      = "tcp"
       DBHOST     = "127.0.0.1:3306"
       DBPORT     = "3306"
       DBNAME     = "your-db-name"
       DBLIMIT    = 1 // rate limit - default: 1 request / second
   )
   ```

5. Build and run the API:

   ```bash
   go run main.go
   ```

6. Your API should now be running at `http://localhost:8000`.

## API Endpoints

The following API endpoints are available:

- `POST /notes` - Create a new note.
- `GET /notes` - Retrieve a list of all notes.
- `GET /notes/:id` - Retrieve a specific note by ID.
- `PUT /notes/:id` - Update a specific note by ID.
- `DELETE /notes/:id` - Delete a specific note by ID.

## Usage Examples

Here are some example requests using `curl` to interact with the API:

- Create a new note:

  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"title": "New Note", "content": "This is a new note."}' http://localhost:8000/notes
  ```

- Retrieve all notes:

  ```bash
  curl http://localhost:8000/notes
  ```

- Retrieve a specific note by ID:

  ```bash
  curl http://localhost:8000/notes/1
  ```

- Update a note by ID:

  ```bash
  curl -X PUT -H "Content-Type: application/json" -d '{"title": "Updated Note", "content": "This note has been updated."}' http://localhost:8000/notes/1
  ```

- Delete a note by ID:

  ```bash
  curl -X DELETE http://localhost:8000/notes/1
  ```
## TODO
- Create test suite (tests/fuzz).

## Contributing

Contributions are welcome! If you'd like to contribute to this project, please follow these steps:

1. Fork the repository on GitHub.
2. Clone your forked repository to your local machine.
3. Create a new branch for your feature or bug fix.
4. Make your changes and commit them.
5. Push your changes to your fork on GitHub.
6. Open a pull request to the main repository.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
