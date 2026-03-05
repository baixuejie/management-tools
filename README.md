# Key Management Tool

Personal key management tool with batch upload, copy templates, and usage tracking.

## Quick Start with Docker

1. Clone the repository
2. Run with Docker Compose:
   ```bash
   docker-compose up -d
   ```
3. Access the application at http://localhost
4. Default login: admin / admin123

## Configuration

Edit environment variables in `docker-compose.yml`:
- Database credentials
- JWT secret
- Encryption key (must be 32 bytes)
- Admin username and password hash

Generate password hash:
```bash
docker exec -it km-backend ./main -hash-password yourpassword
```

## Development

### Backend
```bash
cd backend
go run cmd/main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## Tech Stack

- Backend: Go + Gin + GORM + JWT
- Frontend: Vue 3 + Vant UI + Pinia
- Database: MySQL 8.0
- Deployment: Docker + Docker Compose
