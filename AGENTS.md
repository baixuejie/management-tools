# Repository Guidelines

## Project Structure & Module Organization
- `backend/` contains the Go API service. Entry point is `backend/cmd/main.go`; core layers are under `backend/internal/` (`handlers`, `services`, `models`, `database`, `middleware`, `config`, `utils`).
- `frontend/` contains the Vue 3 app. Main app code is in `frontend/src/` (`views`, `api`, `stores`, `router`, `assets`).
- `docs/plans/` stores design notes and planning docs.
- Root-level infra files: `docker-compose.yml` (full stack), `start.sh` (local startup helper), and `README.md`/`LOCAL_SETUP.md` (setup guidance).

## Build, Test, and Development Commands
- `docker-compose up -d --build`: build and start MySQL, backend, and frontend containers.
- `cd backend && go run cmd/main.go`: run backend API locally on port `8080`.
- `cd frontend && npm install && npm run dev`: run frontend dev server on port `5173`.
- `cd frontend && npm run build`: produce production frontend assets.
- `cd backend && go test ./...`: run backend unit tests.

## Coding Style & Naming Conventions
- Go code must be `gofmt`-clean. Keep package names lowercase and exported identifiers in `PascalCase`.
- Keep handler/service/model separation intact; put new business logic in `services` instead of handlers.
- Vue/JS uses 2-space indentation and `PascalCase` component filenames (for example, `KeyManagement.vue`). Use `camelCase` for variables/functions and keep API calls centralized in `frontend/src/api/`.

## Testing Guidelines
- There are currently no committed frontend test suites; backend tests should be added as `*_test.go` files beside the code under test.
- For backend changes, include or update unit tests and run `cd backend && go test ./...`.
- For frontend changes, at minimum run `cd frontend && npm run build` and verify key flows manually (login, key spec CRUD, batch upload, mark key used).

## Commit & Pull Request Guidelines
- Follow the existing Conventional Commit pattern seen in history: `feat:`, `docs:`, `chore:`.
- Keep commits focused and descriptive (one logical change per commit).
- PRs should include: purpose summary, impacted modules, validation steps/commands, linked issue (if any), and screenshots for UI changes.

## Security & Configuration Tips
- Copy `backend/config.yaml.example` to `backend/config.yaml` for local config; do not commit secrets.
- Set strong values for `KM_JWT_SECRET` and `KM_ENCRYPTION_KEY` (must be 32 bytes).
- Store admin passwords as hashes (see README command for hash generation).
