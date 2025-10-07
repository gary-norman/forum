# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Codex** is a web forum application built with Go that allows users to communicate, create posts, react to content, and organize discussions in channels. The application features user authentication, post/comment systems, reactions (likes/dislikes), channel management, moderation, and image uploads.

## Build & Run Commands

### Local Development (Native)
```bash
# Build the application
make build

# Run the server (requires prior build)
make run

# Build and run in one step
make build && make run

# Run database migrations
bin/codex migrate

# Seed the database with initial data
bin/codex seed
```

### Docker Development
```bash
# Interactive menu for common tasks
make menu

# Configure Docker settings (image name, container name, port)
make configure

# Reset Docker configuration and choose db (dev/prod)
make reset-config

# Build Docker image
make build-image

# Run Docker container
make run-container
```

### Script Management
```bash
# Make scripts executable with correct permissions
make install-scripts

# Backup scripts before updates
make backup-scripts

# Verify script checksums to detect changes
make verify-scripts
```

## Architecture

### Directory Structure

- **`cmd/server/`** - Application entry point
  - `main.go` - Server initialization, CLI commands (migrate/seed), graceful shutdown
  - `migrations.go` - Database migration runner
  - `seed.go` - Database seeding logic

- **`internal/app/`** - Application initialization and configuration
  - Central `App` struct that holds all database models and configuration
  - Loads `.env` file for DB_ENV and DB_PATH settings
  - Initializes SQLite connection and all model instances

- **`internal/db/`** - Database connection and initialization

- **`internal/dao/`** - Generic Data Access Object with type-safe CRUD operations
  - Uses Go generics `DAO[T models.DBModel]`
  - Provides `All()`, `GetByID()`, `Insert()`, `Update()`, `Delete()`

- **`internal/sqlite/`** - SQLite-specific implementations for each model
  - `*-sql.go` files contain model-specific queries
  - Models: Users, Posts, Comments, Reactions, Channels, Memberships, Cookies, Flags, Loyalty, Saved, Mods, Rules, Images, MutedChannels

- **`internal/models/`** - Data models and business logic
  - `*-models.go` files define structs with `db` tags for column mapping
  - Password hashing, UUID generation, validation logic

- **`internal/service/`** - Business logic layer (currently minimal)
  - User creation with password hashing
  - Follows repository pattern for testability

- **`internal/http/`**
  - `handlers/` - HTTP request handlers for auth, posts, comments, channels, users, search, reactions, moderation
  - `middleware/` - Authentication, context injection, logging
  - `routes/` - Router setup and handler dependency injection
    - `registry.go` - Dependency injection for handlers (important!)
    - `routes.go` - Route definitions using Go 1.22+ enhanced servemux

- **`internal/view/`** - HTML template rendering
  - `render.go` - Template initialization and helper functions
  - Custom template functions: `compareAsInts`, `increment`, `dict`, `reactionStatus`, etc.
  - Templates located in `assets/templates/` (`.html` and `.tmpl` files)

- **`migrations/`** - SQL schema files
  - `001_schema.sql` - Main database schema
  - `002_triggers.sql` - Database triggers
  - `003_indexes.sql` - Performance indexes

- **`assets/`** - Static files (CSS, icons, fonts, templates)

- **`scripts/`** - Shell scripts for configuration and deployment
  - `menu.sh` - Interactive menu for Docker operations
  - `configure.sh` - Docker configuration
  - `reset-config.sh` - Reset Docker settings

### Key Architectural Patterns

**Dependency Injection via Handler Registry**
The application uses manual dependency injection in `internal/http/routes/registry.go`. Handlers are created in a specific order to resolve dependencies:
1. Create flat handlers (Session, Reaction, Auth)
2. Create handlers with single-level deps (Comment, Channel)
3. Create complex handlers (User, Post, Home)
4. Wire everything into RouteHandler struct

**Database Models**
All models implement the `DBModel` interface with a `TableName()` method. Struct fields use `db` tags for SQLite column mapping. The generic DAO uses reflection to map between structs and database rows.

**Template Rendering**
Templates are parsed once at startup. The `TempHelper` struct provides access to the App instance for template functions like `reactionStatus`. Templates are composed of partials (`.tmpl` files) and assembled dynamically.

**Authentication Flow**
Session tokens stored in cookies → `auth` middleware extracts user → context injection via `WithUser` middleware → handlers access user from `r.Context()`

**Routes (Go 1.22+ Enhanced Servemux)**
Uses method-based routing: `POST /register`, `GET /post/{postId}`, etc.
Middleware wrapping: `mw.WithUser(http.HandlerFunc(r.Home.GetHome), r.App)`

## Database

- **Type:** SQLite3
- **Driver:** `github.com/mattn/go-sqlite3`
- **Configuration:** `.env` file (DB_ENV, DB_PATH)
- **Schema Location:** `migrations/001_schema.sql`
- **Migrations:** Run with `bin/codex migrate` CLI command
- **Seeding:** Run with `bin/codex seed` CLI command

## Environment Configuration

Create a `.env` file with:
```
DB_ENV=dev
DB_PATH=./identifier.sqlite
```

Use `make configure` or `make reset-config` to set up Docker environment variables.

## Important Notes

- The server runs on port `8888` by default
- A pprof server runs on `localhost:6060` for profiling
- Graceful shutdown on SIGTERM/SIGINT with 10-second timeout
- Image uploads stored in `db/userdata/images/{channel-images,user-images,post-images}/`
- Custom color scheme using Catppuccin Mocha palette (see `internal/colors/`)
- Script backups automatically timestamped in `scripts/backups/`
