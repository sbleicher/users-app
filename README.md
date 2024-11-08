# User Application

Small angular and go application for doing CRUD operations on users

## Running the frontend
```shell
cd users-frontend
npm install
npm start

# For Tests
npm test
```

## Running the backend
```shell
docker compose build --no-cache
docker compose up -d

#For Tests
go test -v ./...
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
