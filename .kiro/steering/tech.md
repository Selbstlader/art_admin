# Tech Stack

## Frontend (`art-design-pro/`)

- **Framework**: Vue 3.5 with Composition API (`<script setup>`)
- **Build**: Vite 7, TypeScript 5.6
- **State**: Pinia with persisted state
- **Router**: Vue Router 4
- **UI**: Element Plus (auto-imported)
- **HTTP**: Axios with custom wrapper
- **Styling**: SCSS with CSS variables for theming
- **i18n**: vue-i18n (zh/en)
- **Charts**: ECharts 5
- **Package Manager**: pnpm (required)

## Backend (`art_admin_backend/`)

- **Language**: Go 1.24
- **Framework**: Gin
- **ORM**: GORM (MySQL/SQLite)
- **Auth**: JWT (golang-jwt/v5)
- **WebSocket**: gorilla/websocket
- **Config**: Viper (YAML)
- **Logging**: Zap with lumberjack rotation
- **Docs**: Swagger (swaggo)
- **AI**: DeepSeek, Dify, Volcengine (Doubao), Baidu OCR

## Common Commands

### Frontend
```bash
cd art-design-pro
pnpm install          # Install dependencies
pnpm dev              # Dev server (port from .env)
pnpm build            # Production build
pnpm lint             # ESLint check
pnpm fix              # ESLint auto-fix
pnpm lint:prettier    # Format code
```

### Backend
```bash
cd art_admin_backend
go mod tidy           # Sync dependencies
go run cmd/server/main.go  # Run server (port 48080)
go build -o server cmd/server/main.go  # Build binary
swag init -g cmd/server/main.go -o docs  # Generate Swagger
```

## Environment

- Node >= 20.19.0
- pnpm >= 8.8.0
- Go >= 1.24
- MySQL or SQLite
