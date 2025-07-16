# Backend Golang Coding Test

## Objective
Build a simple RESTful API in Golang that manages a list of users. Use MongoDB for persistence, JWT for authentication, and follow clean code practices.

---

## Requirements

### 1. User Model
Each user should have:
- `ID` (auto-generated)
- `Name` (string)
- `Email` (string, unique)
- `Password` (hashed)
- `CreatedAt` (timestamp)

---

### 2. Authentication

#### Functions
[x] Register a new user.
[x] Authenticate user and return a JWT.

#### JWT
[x] Use JWT for protecting endpoints.
[x] Use middleware to validate tokens.
[x] Use HMAC (HS256) with a secret key.

---

### 3. User Functions

[x] Create a new user.
[x] Fetch user by ID.
[x] List all users.
[x] Update a user's name or email.
[x] Delete a user.

---

### 4. MongoDB Integration
[x] Use the official Go MongoDB driver.
[x] Store and retrieve users from MongoDB.

---

### 5. Middleware
[x] Logging middleware that logs HTTP method, path, and execution time.

---

### 6. Concurrency Task
[x] Run a background goroutine every 10 seconds that logs the number of users in the DB.

---

### 7. Testing
[x] Write unit tests

Use Go’s `testing` package. Mock MongoDB where possible.

---

## Bonus (Optional)

[x] Add Docker + `docker-compose` for API + MongoDB.
[x] Use Go interfaces to abstract MongoDB operations for testability.
[x] Add input validation (e.g., required fields, valid email).
[x] Implement graceful shutdown using `context.Context`.
- **gRPC Version**
  [x] Create a `.proto` file for `CreateUser` and `GetUser`.
  [x] Implement a gRPC server.
  - (Optional) Secure gRPC with token metadata.
- **Hexagonal Architecture**
  - Structure the project using hexagonal (ports & adapters) architecture:
    [x] Separate domain, application, and infrastructure layers.
    [x] Use interfaces for data access and external dependencies.
    [x] Keep business logic decoupled from frameworks and DB drivers.

---

## Submission Guidelines

- Submit a GitHub repo or zip file.
- Include a `README.md` with:
  - Project setup and run instructions
  - JWT token usage guide
  - Sample API requests/responses
  - Any assumptions or decisions made

---

## Evaluation Criteria

- Code quality, structure, and readability
- REST API correctness and completeness
- JWT implementation and security
- MongoDB usage and abstraction
- Bonus: gRPC, Docker, validation, shutdown
- Testing coverage and mocking
- Use of idiomatic Go

## How to run
- [HOWTORUN.md](./HOWTORUN.md)