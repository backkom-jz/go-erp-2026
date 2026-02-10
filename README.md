# Go ERP

A Go-based ERP backend starter with Gin, GORM, Viper, Zap, Redis, and JWT.

## Features
- Layered architecture: controller/service/repository/model
- Config management via Viper
- Structured logging with Zap
- GORM with MySQL/PostgreSQL support
- Redis integration (optional)
- JWT auth middleware scaffolding
- Unified API response format

## Quick Start
```bash
# Set environment (default: dev)
export APP_ENV=dev

# Run server
go run ./cmd/server
```

## Configuration
Config files:
- `configs/config.yaml`
- `configs/config.dev.yaml`

## Project Structure
```
cmd/
  server/
internal/
  api/
  service/
  repository/
  model/
  middleware/
  pkg/
  bootstrap/
configs/
docs/
migrations/
scripts/
```

## Notes
- Replace DSN in config before running.
- Redis is disabled by default; set `redis.enabled: true` to enable.
