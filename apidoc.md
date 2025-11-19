# UselessMonitor API Documentation

## Authentication

All endpoints require an `Authorization` header containing either the configured `READ_KEY` or `ADMIN_KEY`.

- `READ_KEY` can view monitors and system status.
- `ADMIN_KEY` can view, create, update, and delete monitors.

## Monitor Endpoints

### `GET /monitor`

Retrieves all monitor entries and their latest status.

**Headers**
- `Authorization` (string, required): `READ_KEY` or `ADMIN_KEY`.

**Success Response** (`200 OK`)
```json
[
  {
    "id": 1,
    "name": "API Health Check",
    "type": "API",
    "status": "Healthy",
    "last_check": "2024-06-01T12:00:00Z"
  }
]
```

**Error Responses**
- `401 Unauthorized` when the header is missing.
- `403 Forbidden` when the key does not match.

**Example**
```bash
curl -H "Authorization: $READ_KEY" http://localhost:8080/monitor
```

---

### `POST /monitor`

Creates a new monitor entry.

**Headers**
- `Authorization` (string, required): `ADMIN_KEY`.

**Request Body**
```json
{
  "name": "API Health Check",
  "type": "API",
  "status": "Healthy"
}
```

**Success Response** (`201 Created`)
```json
{
  "id": 1,
  "name": "API Health Check",
  "type": "API",
  "status": "Healthy",
  "last_check": "2024-06-01T12:00:00Z"
}
```

**Error Responses**
- `400 Bad Request` when the payload is invalid.
- `401 Unauthorized` when the header is missing.
- `403 Forbidden` when the key does not match the admin key.
- `500 Internal Server Error` when persistence fails.

**Example**
```bash
curl -X POST \
  -H "Authorization: $ADMIN_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "API Health Check", "type": "API", "status": "Healthy"}' \
  http://localhost:8080/monitor
```

---

### `PUT /monitor/:id`

Updates a monitor entry by ID.

**Headers**
- `Authorization` (string, required): `ADMIN_KEY`.

**Request Body**
Any combination of updatable fields:
```json
{
  "name": "API Health Check",
  "type": "API",
  "status": "Unhealthy"
}
```

**Success Response** (`200 OK`)
```json
{
  "id": 1,
  "name": "API Health Check",
  "type": "API",
  "status": "Unhealthy",
  "last_check": "2024-06-01T12:05:00Z"
}
```

**Error Responses**
- `400 Bad Request` when the payload is invalid.
- `401 Unauthorized` when the header is missing.
- `403 Forbidden` when the key does not match the admin key.
- `404 Not Found` when the monitor does not exist.
- `500 Internal Server Error` when persistence fails.

**Example**
```bash
curl -X PUT \
  -H "Authorization: $ADMIN_KEY" \
  -H "Content-Type: application/json" \
  -d '{"status": "Unhealthy"}' \
  http://localhost:8080/monitor/1
```

---

### `DELETE /monitor/:id`

Deletes a monitor entry by ID.

**Headers**
- `Authorization` (string, required): `ADMIN_KEY`.

**Success Response** (`200 OK`)
```json
{"message": "Monitor deleted"}
```

**Error Responses**
- `401 Unauthorized` when the header is missing.
- `403 Forbidden` when the key does not match the admin key.
- `500 Internal Server Error` when deletion fails.

**Example**
```bash
curl -X DELETE -H "Authorization: $ADMIN_KEY" http://localhost:8080/monitor/1
```

---

## Status Endpoint

### `GET /status`

Returns a summary of the system health based on stored monitors.

**Headers**
- `Authorization` (string, required): `READ_KEY` or `ADMIN_KEY`.

**Success Response** (`200 OK`)
```json
{
  "status": "OK",
  "monitors": 3,
  "healthy_monitors": 3
}
```

**Error Responses**
- `401 Unauthorized` when the header is missing.
- `403 Forbidden` when the key does not match.
- `500 Internal Server Error` when the summary cannot be generated.

**Example**
```bash
curl -H "Authorization: $READ_KEY" http://localhost:8080/status
```
