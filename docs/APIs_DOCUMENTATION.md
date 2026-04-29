# ERP Microservices API Documentation

This project uses [swag](https://github.com/swaggo/swag) to generate OpenAPI 2.0 (Swagger) documentation from Go source code comments.

## Prerequisites

To generate the documentation, you need to have `swag` installed on your machine:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Generating Documentation

Run the provided PowerShell script from the `scripts` directory:

```powershell
./scripts/generate-docs.ps1
```

This script will:
1. Initialize/Update Swagger docs for `hr-service`.
2. Initialize/Update Swagger docs for `finance-service`.
3. Generate `docs/swagger.json`, `docs/swagger.yaml`, and `docs/docs.go` in each service directory.

## Viewing Documentation

Once generated, the documentation can be accessed via the following endpoints (when services are running):

### Direct Service Access
- **HR Service:** `http://localhost:8081/swagger/index.html` (If exposed)
- **Finance Service:** `http://localhost:8082/swagger/index.html` (If exposed)

### Through API Gateway
- **HR Service:** `http://localhost:8080/api/v1/hr/swagger/index.html`
- **Finance Service:** `http://localhost:8080/api/v1/finance/swagger/index.html`

> **Note:** To enable these endpoints, the swagger middleware must be configured in each service's `main.go`.

## Adding Documentation to New Endpoints

To document a new endpoint, add comments above the handler function following the [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format).

Example:
```go
// @Summary Get User
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
```
