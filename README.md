# 1337b04rd

A minimalist hacker imageboard application written in Go (Golang) and PostgreSQL.

---

## 🧠 Overview

This project is a small web application for anonymous posting and commenting with support for image uploads, session persistence, and automatic archival of inactive threads.

---

## ✨ Features

✅ Create threads (posts) with text and/or images
✅ Anonymous sessions with avatars
✅ Add comments and replies (with image support)
✅ Archival logic for inactive threads
✅ Auto-generated UUID for sessions, posts, and comments
✅ Session name override (updates all posts/comments)
✅ Static frontend with HTML templates
✅ Clean logging and error handling
✅ Test coverage for service logic
✅ Session middleware using cookies
✅ Dockerized environment (App + PostgreSQL)

---

## 📦 Tech Stack

* **Go** (1.21+)
* **PostgreSQL** (15)
* **HTML + CSS** for frontend
* **Docker & Docker Compose**
* **Standard library only** (no third-party libraries allowed per requirements)

---

## 🧪 Running the Project

### 🔧 Prerequisites

* Docker & Docker Compose
* Go 1.21 (if running manually without Docker)

### 🚀 Option 1: Run with Docker

```bash
docker-compose up --build
```

Visit: [http://localhost:8080](http://localhost:8080)

### 🧪 Option 2: Run manually

```bash
cd cmd
go run 1337b04rd
```

Ensure PostgreSQL is running and `.env` variables are set or defaults will be used.

---

## 🗂️ Project Structure

```
1337b04rd/
├── cmd/1337b04rd         # Main entrypoint
├── config
├── data
├── internal/
│   ├── adapters/
│   │   ├── handler       # HTTP Handlers
│   │   ├── middleware    # Session middleware
│   │   └── repo/         # PostgreSQL Repos
│   ├── domain/           # Models and interfaces
│   └── service/          # Business logic
├── pkg/                  # Logger, utils (e.g., UUID)
├── static/               # HTML, CSS, assets
├── logging/              # App log file
├── Dockerfile
├── docker-compose.yml
└── init.sql              # DB schema
```

---

## 📜 API Routes

| Method | Endpoint               | Description                         |
| ------ | ---------------------- | ----------------------------------- |
| GET    | `/`                    | View catalog (non-archived threads) |
| GET    | `/archive`             | View archived threads               |
| GET    | `/posts/{id}`          | View thread with comments           |
| GET    | `/create`              | Form to create a new thread         |
| POST   | `/posts`               | Submit new thread                   |
| POST   | `/posts/{id}/comments` | Submit a comment (or reply)         |
| GET    | `/error`               | Render error page                   |

---

## 🧪 Testing

Run service-layer unit tests:

```bash
go test ./internal/service/... -v -cover
```

Check full test coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 📝 Notes

* **Sessions** are stored with a UUID and cookie, and last for **7 days**.
* Users may **override their display name** mid-session. All their posts/comments update accordingly.
* Session avatars are fetched randomly via API when a session is first created.
* Archival logic:

  * Threads with **no comments** are archived after **10 minutes**.
  * Threads with comments are archived **15 minutes** after the latest comment.
* Archival is performed during read/write operations or scheduled via timer (you can expand this).
* Filenames are validated, and images are uploaded to `/data`.

---

## ⚙️ CLI Help

```bash
./1337b04rd --help
```

```
hacker board

Usage:
  1337b04rd [--port <N>]
  1337b04rd --help

Options:
  --help       Show this screen.
  --port N     Port number.
```

---

## 🧩️ Future Ideas

* Implement graceful shutdown
* Background worker for archival
* Admin/mod panel
* CAPTCHA / spam protection
* Image size validation and resizing
* CSS polish for mobile

---