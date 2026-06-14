# Basalam Feedback Board

Minimal Feedback Board (Go + PostgreSQL + Bootstrap).

Run with Docker Compose:

```bash
docker compose up --build
```

The app will be available at http://localhost:8080 and the admin dashboard at http://localhost:8080/admin (use `admin`/`admin` by default).

Run locally (requires a running Postgres instance):

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/feedbacks?sslmode=disable
export BASIC_AUTH_USER=admin
export BASIC_AUTH_PASSWORD=admin
go build -o feedbackboard ./cmd
./feedbackboard
```

Notes:
- The app auto-creates the `feedbacks` table and the `pgcrypto` extension on startup.
- Public endpoint: `POST /api/feedbacks` — returns `201 Created` on success.
- Admin endpoints (Basic Auth protected): `GET /api/feedbacks`, `PATCH /api/feedbacks/{id}/status`.
