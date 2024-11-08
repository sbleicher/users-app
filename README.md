# User Application

Small angular and go application for doing CRUD operations on users

## Running the frontend
```shell
cd users-frontend
npm install
npm start
```
Testing the frontend:
```shell
cd users-frontend
npm test
```

## Running the backend
```shell
docker compose build --no-cache
docker compose up -d
```
Testing the backend:
```shell
cd users-backend
go test -v ./...
```

To run without docker (Not tested):
```shell
cd users-backend
go mod download && go mod verify
go clean -cache && go build -o main .

chmod +x main

# This will require you have a postgres database setup locally
DATABASE_URL=postgres://<username>:<password>@localhost:5432/users?sslmode=disable ./main
```

## Swagger
Hosted at: http://localhost:8080/swagger/index.html

```shell
cd users-backend/
swag init # To generate a new set of swagger documents
```

## Next Steps
Here are some things I would look to improve if I spent some more time on this.

Frontend:
- Pagination on users table
- Readonly View Page
- "Are you sure" popup or dialog for frontend delete
- Mobile and tablet responsive
- Refactor some of the frontend tests
- E2E testing for happy and error paths
- Update endpoint specific to what changed
- Linter
- Filter and sort on the users list
- Configs

Backend:
- Pagination on get all users endpoint
- Update endpoint specific to what changed
- Refactoring the handler tests
- Configs
