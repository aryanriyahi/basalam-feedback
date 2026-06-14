# Technical Specification: Feedback Board

## 1. Overview
[cite_start]A simple system to submit and track feedback. [cite_start]Features a user-facing submission page [cite: 6] [cite_start]and an admin dashboard to view and change feedback statuses[cite: 7].

## 2. Database Schema (PostgreSQL)
Table: `feedbacks`
- `id`: UUID (Primary Key)
- [cite_start]`title`: VARCHAR(255), Not Null [cite: 6]
- [cite_start]`message`: TEXT, Not Null [cite: 6]
- [cite_start]`status`: VARCHAR(50), Default: 'submitted' [cite: 9]
- `created_at`: TIMESTAMP, Default: NOW()

[cite_start]*Allowed Statuses*: `submitted` (ثبت شده), `reviewing` (در حال بررسی), `resolved` (رسیدگی شده)[cite: 9, 10].

## 3. API Endpoints Contract
### A. Public Endpoints
- `POST /api/feedbacks`
  - Body: `{ "title": "string", "message": "string" }`
  - Returns: `201 Created`

### B. Admin Endpoints
- `GET /api/feedbacks`
  - Returns: `200 OK`, Array of feedback objects (sorted by `created_at` descending).
- `PATCH /api/feedbacks/{id}/status`
  - Body: `{ "status": "reviewing" | "resolved" }`
  - Returns: `200 OK`

### C. Bonus Feature (Optional Auth)
- [cite_start]Implementation of a simple HTTP Basic Auth or Hardcoded Token middleware for Admin Endpoints to prevent unauthorized status changes.

## 4. Frontend Views
1. [cite_start]`/` -> `ui/index.html`: Contains a simple Bootstrap form (Title, Message) with a success/error toast[cite: 6].
2. `/admin` -> `ui/dashboard.html`: A responsive table listing all feedbacks. [cite_start]Includes a dropdown/button per row to update the status via the `PATCH` API[cite: 7].