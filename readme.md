# Todo API with MongoDB

## Overview
This is a RESTful API for a simple Todo application built using Go, Chi router, and MongoDB. The API allows users to create, read, update, and delete (CRUD) todo items.

## Features
- Fetch all todo items
- Create a new todo item
- Update an existing todo item
- Delete a todo item
- Graceful shutdown handling
- Middleware for logging requests

## Technologies Used
- Go (Golang)
- Chi Router (for routing)
- MongoDB (for database storage)
- mgo.v2 (MongoDB driver for Go)
- go-chi middleware (for logging)

## Project Structure
```
/ ├── main.go    # Entry point of the application
  ├── README.md  # Documentation
  ├── static/    # HTML templates (if needed)
  ├── go.mod     # Go module dependencies
```

## Prerequisites
Before running the project, ensure you have the following installed:
- Go (1.16+ recommended)
- MongoDB (running on localhost:27017 or modify `hostName` in the code)

## Installation & Setup
1. Clone the repository:
   ```sh
   git clone <repo-url>
   cd <repo-folder>
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Run the application:
   ```sh
   go run main.go
   ```
4. API will be available at `http://localhost:9000/`

## API Endpoints
### Fetch All Todos
- **URL:** `GET /todo/`
- **Response:**
  ```json
  {
    "data": [
      {
        "id": "611f13d9...
        "title": "Buy groceries",
        "completed": false,
        "createdAt": "2023-07-20T14:00:00Z"
      }
    ]
  }
  ```

### Create a Todo
- **URL:** `POST /todo/`
- **Body:**
  ```json
  {
    "title": "New Todo"
  }
  ```
- **Response:**
  ```json
  {
    "message": "todo created successfully",
    "todo_id": "6123abcd..."
  }
  ```

### Update a Todo
- **URL:** `PUT /todo/{id}`
- **Body:**
  ```json
  {
    "title": "Updated Task",
    "completed": true
  }
  ```
- **Response:**
  ```json
  {
    "message": "todo updated successfully"
  }
  ```

### Delete a Todo
- **URL:** `DELETE /todo/{id}`
- **Response:**
  ```json
  {
    "message": "todo deleted successfully"
  }
  ```

## Error Handling
- Invalid ID format or missing title will return a `400 Bad Request`
- Internal server/database errors return `500 Internal Server Error`

## Graceful Shutdown
- The server listens for `os.Interrupt` signals and gracefully shuts down all processes before exiting.

## Contributing
1. Fork the repository
2. Create a feature branch (`git checkout -b feature-xyz`)
3. Commit your changes (`git commit -m "Added feature xyz"`)
4. Push to the branch (`git push origin feature-xyz`)
5. Open a Pull Request

## License
This project is licensed under the MIT License.