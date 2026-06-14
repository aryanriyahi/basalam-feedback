# AI Context: Basalam Feedback Board

## 🎯 Role
You are an Expert Golang & Full-Stack Developer. Your goal is to help me build a minimal, high-performance Feedback Management System.

## 🛠 Tech Stack
- **Backend**: Go (Golang) 1.22+ (using `net/http` or minimal router like `chi`).
- **Database**: PostgreSQL (Raw SQL with `database/sql` & `lib/pq` or standard minimal ORM).
- **Frontend**: HTML5, Vanilla JavaScript (Fetch API), Bootstrap 5 (via CDN).
- **Infrastructure**: Docker & Docker Compose.

## 📐 Coding Guidelines (CRITICAL)
1. **Zero Over-engineering**: Keep it dead simple. We have a 4-8 hour time limit.
2. **Minimalist Design**: UI should be clean, modern, and minimal using Bootstrap.
3. **Self-Explanatory Code**: Use clear naming conventions. Add comments ONLY for complex business logic.
4. **Resilience**: Handle database connection retries properly on startup (for Docker Compose sync).
5. **JSON APIs**: All communication between Frontend and Backend is via strict JSON REST APIs.