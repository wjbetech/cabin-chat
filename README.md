# CABIN-CHAT

This application serves as my first foray into Go. The plan is to build a mini-Discord/IRC/WhatsApp type platform that will be locally hosted on my home lab, allowing for me and my close friends to log in with working auth and chat to each other, react, and send photos.

To run this project locally, you should/need:

- Go installed
- Docker Desktop installed and running
- Git Bash and/or PowerShell available

## Project Overview

`cabin-chat` is a small, dedicated, self-hosted chat server that I built for my friends and myself, hosted on my home lab, running 24/7. It is inspired by my growing up on irc channels, Discord and WhatsApp, and will integrate much of the functionality of those services; users, messages, reactions.

The project is being built in Go with PostgreSQL for persistent storage, and Docker Compose used to run the local DB during development.

The project is currently in `Phase 2 - Backend MVP`, and at this stage the Go server can:

- load the config from env vars
- start an HTTP server
- exposes a health endpoint
- optionally connect to PostgreSQL if a DB URL is configured

Current work focuses on wiring the core architecture first, then replacing in-memory behavior with PostgreSQL-backed implementation.

## Backend Architecture

The backend is currently organized into a few small packages:

- an entry point for the app at `cmd/cabin-chat`
- app config config defined in the `/pkg/env` directory
- A bootstrapped DB connection defined in the `/pkg/db` directory
- A set of core domain models, primarily defined as `User`, `Message` and `Reaction` in the `/pkg/chat` directory
- A store for the repo interfaces and in-memory store groundwork in the `/pkg/store` directory
- `documentation/roadmap.md` contains the current implementation checklist and build plan, which may be helpful to contributors

#### Current BE Status

- the Go server starts successfully
- the `/health` endpoint responds successfully
- PostgreSQL connection bootstrap is implemented
- local Docker Compose setup for PostgreSQL is being introduced
- in-memory and DB-backed persistence are WIP

#### Port Layout

The project currently uses separate ports for the Go server and PostgreSQL:

- `1111` — Go HTTP server running on the host machine
- `1122` — PostgreSQL exposed on the host machine through Docker
- `5432` — PostgreSQL's internal port inside the Docker container

#### Important!

- the Go app connects to PostgreSQL using `localhost:1122` because the Go app currently runs on the host machine, not inside Docker
- PostgreSQL still listens on `5432` inside the container, which is why the Docker port mapping forwards host `1122` to container `5432`

## env

The backend currently uses the following environment variables:

- `CABIN_CHAT_PORT` — the port the Go HTTP server listens on locally
- `CABIN_CHAT_JWT_SECRET` — the secret that will later be used for signing JWTs
- `CABIN_CHAT_DATABASE_URL` — the PostgreSQL connection string used by the Go server
- `CABIN_CHAT_UPLOAD_DIR` — the local directory where uploaded files will eventually be stored

Example values:

```env
CABIN_CHAT_PORT=port
CABIN_CHAT_JWT_SECRET=change-me
CABIN_CHAT_DATABASE_URL=postgres://cabin_admin:cabin_admin@localhost:port/cabin_chat?sslmode=disable
CABIN_CHAT_UPLOAD_DIR=./uploads
```

This app currently reads environment variables from the shell process. It does not automatically load values from `.env` unless you source that file first.

## Local Development

The current local development workflow is:

1. Start PostgreSQL with Docker Compose
2. Export the required environment variables
3. Run the Go server
4. Verify the health endpoint
5. Run formatting, build, and test checks regularly (run `bash scripts/cleanup.sh`)

A typical workflow might look like this:

```bash
docker compose up -d postgres
export CABIN_CHAT_DATABASE_URL="postgres://username:password@localhost:port/cabin_chat?sslmode=disable"
go run ./cmd/cabin-chat
```

## Roadmap

The project roadmap and current implementation checklist live in documentation/roadmap.md.

Current focus:

- finish the Backend MVP in Phase 2
- bring PostgreSQL fully online for local development
- create schema and persistence layers
- add auth and messaging endpoints
