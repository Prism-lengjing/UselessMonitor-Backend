# UselessMonitor-Backend

A lightweight Gin + SQLite service that tracks monitor entries with key-based access control.

## Getting Started

1. Set required environment variables (use a `.env` file or export directly):
   ```bash
   READ_KEY=example-read-key
   ADMIN_KEY=example-admin-key
   ```
2. Install Go dependencies and run the server:
   ```bash
   # adjust GOPROXY if your environment blocks the default proxy
   GOPROXY=https://proxy.golang.org,direct go mod tidy
   go run .
   ```
3. The API listens on port `8080` by default.

## API Overview

- `GET /monitor` — list monitors (requires `READ_KEY` or `ADMIN_KEY`).
- `POST /monitor` — create a monitor (requires `ADMIN_KEY`).
- `PUT /monitor/:id` — update a monitor (requires `ADMIN_KEY`).
- `DELETE /monitor/:id` — delete a monitor (requires `ADMIN_KEY`).
- `GET /status` — system status summary (requires `READ_KEY` or `ADMIN_KEY`).

See [`apidoc.md`](apidoc.md) for full request and response examples.
