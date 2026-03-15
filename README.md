# Babel Shelf

A CRUD bookshelf app built in every language. Same API, same database, different implementations.

## Architecture

```
babel-shelf/
├── db/                  # Shared Postgres schema
│   └── init.sql
├── go-bookshelf/        # Go (port 8080)
├── ts-bookshelf/        # TypeScript/Bun (port 8081)
├── docker-compose.yml   # Orchestrates apps + Postgres
└── README.md
```

Each language gets its own directory (`<lang>-bookshelf/`) with a Dockerfile and devcontainer config. They all share the same Postgres database and expose the same REST API.

## API

| Method | Endpoint       | Description       | Status |
|--------|---------------|-------------------|--------|
| POST   | `/books`      | Create a book     | 201    |
| GET    | `/books`      | List all books    | 200    |
| GET    | `/books/{id}` | Get one book      | 200    |
| PUT    | `/books/{id}` | Update a book     | 200    |
| DELETE | `/books/{id}` | Delete a book     | 204    |

### Book object

```json
{
  "id": 1,
  "title": "Dune",
  "author": "Frank Herbert",
  "status": "want to read"
}
```

Status must be one of: `want to read`, `reading`, `finished`.

## Running

```bash
docker compose up --build
```

| Implementation | URL |
|---------------|-----|
| Go | `http://localhost:8080` |
| TypeScript | `http://localhost:8081` |

## Testing

Tests are integration tests that hit a real Postgres database (no mocks).

```bash
# Go — from inside devcontainer or with Go installed:
cd go-bookshelf
DATABASE_URL="postgres://shelf:shelf@db:5432/bookshelf?sslmode=disable" go test -v

# TypeScript — from inside devcontainer or with Bun installed:
cd ts-bookshelf
DATABASE_URL="postgres://shelf:shelf@db:5432/bookshelf?sslmode=disable" bun test
```

## Implementations

- [x] Go
- [x] TypeScript (Bun + Hono + pg)
