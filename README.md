# Invoice

A professional, high-fidelity invoice generation and management platform built for modern freelancers and small teams.
Invoice provides a frictionless WYSIWYG editor to create, send, and track professional invoices instantly.

## Core Features

- **WYSIWYG Invoice Editor**: Direct on-template editing ensuring pixel-perfect output.
- **PDF Export**: Generate flawless, print-ready PDF documents built on a strict fixed grid layout.
- **Secure Client Sharing**: Generate unique URLs to track invoice views and payment initiation.
- **Role-Based Access Control (RBAC)**: Comprehensive user role and permission management.
- **Stateful Session Management**: Secure, PostgreSQL-backed session handling for authentication.
- **Modern UI Architecture**: Premium glassmorphism aesthetics powered by Tailwind CSS and custom design tokens.

## Technology Stack

- **Backend**: Go (Golang), standard library `net/http`
- **Database**: PostgreSQL
- **ORM**: GORM
- **Migrations**: Goose
- **Session Management**: SCS
- **Frontend**: HTML5 Templates, Tailwind CSS (Runtime Configuration), Vanilla JavaScript

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- [Air](https://github.com/cosmtrek/air) (for live reloading)
- [Goose](https://github.com/pressly/goose) (for database migrations)

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd invoice-builder
   ```

2. **Configure the Database**
   Ensure PostgreSQL is running and update the `DB_DSN` inside the `Makefile` and `cmd/web/main.go` to match your local
   database credentials.

3. **Apply Database Migrations**
   ```bash
   make migrate-up
   ```

## Development

The project includes a comprehensive `Makefile` to streamline development workflows.

- **Start the development server with live reloading**
  ```bash
  make air
  ```

- **Run the server standardly**
  ```bash
  make run
  ```

- **Build the executable**
  ```bash
  make build
  ```

- **Run tests**
  ```bash
  make test
  ```

### Database Migration Commands

- `make migrate-status` - Check the status of all migrations
- `make migrate-up` - Apply all pending migrations
- `make migrate-down` - Rollback the most recent migration
- `make migrate-reset` - Rollback all migrations to zero
- `make migrate-create name=<name>` - Generate a new SQL migration file


# invoice-builder
