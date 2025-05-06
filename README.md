# SWIFT-Codes Service

A self-contained microservice written in Go 1.22 that:

* Parses the official list of **SWIFT/BIC** codes from an Excel file,
* Stores the data in a PostgreSQL database,
* Exposes a REST API at `localhost:8080`.

You can run the service with a single Docker Compose command — no external setup required.

---

## Table of Contents

1. Quick Start
2. Example Requests
3. Testing
4. Project Structure
5. Environment Variables

---

## Quick Start

```bash
git clone https://github.com/YOUR_USERNAME/swift-codes-service
cd swift-codes-service
docker compose up --build        # first launch (~30s)
```

* `db` – PostgreSQL 16 with a `db-data` volume
* `api` – compiled in a **distroless** container and listens on port `:8080`

---

## Example Requests

| Description                          | Method / URI                             | Example |
|--------------------------------------|------------------------------------------|---------|
| HQ details + branch list             | `GET /v1/swift-codes/{swift}`            | http://localhost:8080/v1/swift-codes/AAISALTRXXX |
| Branch details only                  | `GET /v1/swift-codes/{swift}`            | http://localhost:8080/v1/swift-codes/BANKPLPWA01 |
| All SWIFT codes for a country        | `GET /v1/swift-codes/country/{ISO2}`     | http://localhost:8080/v1/swift-codes/country/PL |
| Add new SWIFT code                   | `POST /v1/swift-codes`                   | see payload below |
| Delete a SWIFT code                  | `DELETE /v1/swift-codes/{swift}`         | DELETE /v1/swift-codes/TESTPLPWXXX |

**Sample POST body**

```json
{
  "address": "Street 1",
  "bankName": "Demo Bank",
  "countryISO2": "PL",
  "countryName": "POLAND",
  "isHeadquarter": false,
  "swiftCode": "TESTPLPWXXX"
}
```

---

## Testing

The project includes two test types using **build tags**:

| Type            | Description                               | Command |
|------------------|--------------------------------------------|---------|
| **Unit**         | parser, repository, and HTTP handler tests using SQLite | `go test ./... -tags=unit` |
| **Integration**  | full end-to-end tests using PostgreSQL via Testcontainers | `go test ./... -tags=integration` |
| **Full suite**   | both unit and integration (recommended for CI) | `go test ./... -tags=unit,integration` |

Coverage report:

```bash
go test ./... -tags=unit,integration -coverprofile=cover.out
go tool cover -html=cover.out
```

> **Linux note:** To run integration tests without `sudo`, add your user to the `docker` group:  
> `sudo usermod -aG docker $USER && newgrp docker`

---

## Project Structure

```
.
├── cmd/
│   └── server/             → application entry point (main.go)
├── internal/               → application logic
│   ├── db/                 → PostgreSQL connection via Gorm
│   ├── models/             → model definitions (SwiftCode)
│   ├── parser/             → Excel parsing logic (.xlsx)
│   ├── repository/         → data access layer (CRUD)
│   └── transport/
│       └── http/           → HTTP handlers + Gin router
├── tests/                  → unit and integration tests
│   ├── http/               → HTTP handler tests (SQLite)
│   ├── integration/        → end-to-end tests using Testcontainers and Postgres
│   ├── parser/             → Excel parsing tests
│   └── repository/         → repository logic tests
├── Dockerfile              → multi-stage build (Go + distroless)
├── docker-compose.yml      → launches Postgres and the API locally
├── Interns_2025_SWIFT_CODES.xlsx  → input data file for SWIFT codes
├── go.mod / go.sum         → dependency definitions
└── README.md               → project documentation
```
