# How to run

## Run the application using Docker Compose
Ensure you have Docker and Docker Compose installed on your machine.
```
docker-compose -f .\docker\docker-compose.yaml up -d
```

## Access the application
- HTTP API: `http://localhost:8080`
- gRPC API: `localhost:9090`

## Access MongoDB
- Mongo Express: `http://localhost:8081`
- MongoDB CLI: Connect using `mongodb://admin:password@localhost:27017/backend_challenge`

## Run tests
```
cd ..
go test ./... -v
```

## Run the HTTP server
```
go run cmd/http/main.go
```

## Run the gRPC server
```
go run cmd/grpc/main.go
```

## Access the gRPC server
Use a gRPC client to connect to `localhost:9090`.

## Generate gRPC code
```
protoc --go_out=. --go-grpc_out=. internal\infrastructure\grpc\proto\user.proto
```

## Usage of JWT tokens
- Include JWT in the `Authorization` header of requests.
- Example:
```Authorization
Bearer <your_jwt_token>
```

## Postman Collection
- [postman_collection.json](./postman_collection.json)

## Sample HTTP requests public endpoints
```
curl --location 'http://localhost:8080/api/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "admin@example.com",
    "password": "password"
}'
```

## Sample HTTP response public endpoints
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjg3NzZlYWRkYzhlYmFjNjMwZDg2MWUwIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImV4cCI6MTc1Mjc0NTk4MSwiaWF0IjoxNzUyNjU5NTgxfQ.Q_vFdVZPEAp0DULnsHaLgiBGoxdwjvwolZ82dc6M6lU",
    "user": {
        "id": "68776eaddc8ebac630d861e0",
        "name": "Admin User",
        "email": "admin@example.com",
        "created_at": "2025-07-16T09:19:41.417Z"
    }
}
```

## Sample HTTP requests With JWT authorized
```
curl --location 'http://localhost:8080/api/users' \
--header 'Authorization: Bearer <your_jwt_token>'
```

## Sample HTTP response With JWT authorized
```
{
    "users": [
        {
            "id": "68776eaddc8ebac630d861e0",
            "name": "Admin User",
            "email": "admin@example.com",
            "created_at": "2025-07-16T09:19:41.417Z"
        }
    ],
    "total": 1
}
```

## Sample gRPC Example
- Use a gRPC client to connect to `localhost:9090`.
- Use Tool Like [Postman](https://www.postman.com/) to send requests.
- Using server reflection, you can explore available services and methods.

## Sample gRPC requests
```
localhost:9090  UserService/GetAllUsers
```

## Sample gRPC response
```
{
    "users": [
        {
            "id": "68776eaddc8ebac630d861e0",
            "name": "Admin User",
            "email": "admin@example.com",
            "created_at": "2025-07-16T09:19:41Z"
        }
    ]
}
```