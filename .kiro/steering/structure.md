# Project Structure

## Frontend (`art-design-pro/src/`)

```
src/
├── api/              # API request functions (one file per domain)
├── assets/
│   ├── icons/        # Icon fonts
│   ├── img/          # Images by category
│   └── styles/       # Global SCSS (variables, mixins, themes)
├── components/
│   ├── core/         # Reusable UI components (cards, charts, forms, tables)
│   └── custom/       # Business-specific components
├── composables/      # Vue composables (useAuth, useTable, useTheme, etc.)
├── config/           # App configuration (images, header, fast-enter)
├── directives/       # Custom directives (auth, highlight, ripple, roles)
├── enums/            # TypeScript enums
├── locales/          # i18n translations (zh.json, en.json)
├── mock/             # Mock data for development
├── router/           # Vue Router configuration
├── store/modules/    # Pinia stores (user, menu, setting, worktab, table)
├── types/            # TypeScript type definitions
├── typings/          # Global type declarations (.d.ts)
├── utils/            # Utilities organized by domain:
│   ├── browser/      # BOM, cookies
│   ├── http/         # Axios wrapper, error handling
│   ├── navigation/   # Route helpers, worktab management
│   ├── storage/      # LocalStorage abstraction
│   ├── theme/        # Theme switching
│   └── ui/           # Colors, loading, tabs
└── views/            # Page components by feature module
```

## Backend (`art_admin_backend/`)

```
art_admin_backend/
├── cmd/server/       # Application entry point (main.go)
├── config/           # Configuration files (config.yaml)
├── docs/             # Swagger documentation
├── internal/
│   ├── api/          # HTTP handlers
│   │   ├── middleware/   # Auth, CORS, logging, recovery
│   │   ├── v1/           # Versioned API handlers
│   │   ├── mood/         # Mood module handlers
│   │   ├── project/      # Project module handlers
│   │   └── travel/       # Travel module handlers
│   ├── dto/          # Data Transfer Objects
│   │   ├── request/      # Request DTOs
│   │   └── response/     # Response DTOs
│   ├── model/        # GORM models (database entities)
│   ├── pkg/          # Internal packages
│   │   ├── config/       # Config loader
│   │   ├── database/     # DB connection & migrations
│   │   ├── jwt/          # JWT utilities
│   │   ├── logger/       # Zap logger setup
│   │   ├── response/     # Standard API responses
│   │   ├── websocket/    # WebSocket hub
│   │   └── [ai-services] # DeepSeek, Dify, Volcengine, BaiduOCR
│   ├── repository/   # Data access layer (one repo per model)
│   ├── router/       # Route registration
│   └── service/      # Business logic layer
├── migrations/       # SQL migration files
└── logs/             # Application logs
```

## Architecture Patterns

### Frontend
- **API Layer**: `src/api/*.ts` → wraps HTTP client
- **State**: Pinia stores with persistence
- **Views**: Feature-based folders under `src/views/`
- **Path Aliases**: `@/` = `src/`, `@views/`, `@utils/`, `@stores/`, `@styles/`

### Backend
- **Layered Architecture**: Handler → Service → Repository → Model
- **Dependency Injection**: Services initialized in main.go, passed to handlers
- **Response Format**: Standardized via `internal/pkg/response`
- **API Versioning**: `/api/v1/` prefix for versioned endpoints
