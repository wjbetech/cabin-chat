## Overview & Vision

Cabin Chat is a locally hosted chat application for a small friend group.

The project will be built in this order:

1. Plan the system and document the path from backend to deployment.
2. Build the Go backend to a clean, maintainable standard.
3. Build a modern frontend with Vite, React, TypeScript, TailwindCSS, and shadcn.
4. Connect the frontend to the backend locally.
5. Deploy and harden the full application on the home server.

The target experience is a modern WhatsApp-style chat app with:

- login and authentication
- real-time messaging
- reactions
- photo uploads
- local self-hosted deployment on the home lab

## Principles

- Build the backend first, then the frontend, then integration, then deployment.
- Keep the MVP small and working before adding optional features.
- Use simple auth first: username or email, password, bcrypt, and JWT.
- Treat every phase as incomplete until it can be tested locally.
- Prefer clear folder boundaries so each concern has one home.

## Core Architecture

### Environment and config

Configuration will use `CABIN_CHAT_*` environment variables loaded through `pkg/env`.

The server startup flow will be:

`config -> logger -> stores -> services -> router -> WebSocket hub -> HTTP server`

### Planned backend layout

- `cmd/cabin-chat` for the application entry point
- `internal/api` for HTTP handlers and route wiring
- `pkg/env` for configuration loading
- `pkg/auth` for password hashing, JWT, and auth helpers
- `pkg/chat` for chat domain logic and WebSocket hub code
- `pkg/db` for database setup and implementations
- `pkg/store` for repository interfaces

### Planned backend features

- environment-driven configuration
- health endpoint
- signup and login endpoints
- bcrypt password hashing
- JWT issuing and validation
- user and message persistence with Postgres
- WebSocket-based real-time chat
- reactions and media upload support
- tests using `go test ./...`

### Planned frontend stack

- Vite
- React
- TypeScript
- TailwindCSS
- shadcn/ui

The frontend goal is a cute, modern, WhatsApp-style chat interface rather than a generic placeholder SPA.

## Out Of Scope For The First MVP

- Discord auth
- OAuth providers
- mobile app clients
- push notifications
- advanced admin or moderation tooling
- multi-server or federated chat

## Build Phases

### Phase 1 - Foundation

Depends on: none

Goals:

- get the Go application compiling
- load configuration from environment variables
- start a basic HTTP server
- expose a health endpoint

Tasks:

- [x] Confirm the repo layout
- [x] Confirm the environment variable naming convention (`CABIN_CHAT_*`)
- [x] Finish `pkg/env/config.go`
- [x] Create `cmd/cabin-chat/main.go`
- [x] Start an HTTP server in `main.go`
- [x] Add a `GET /health` endpoint
- [x] Run `go fmt ./...`

Done when:

- `go build ./...` succeeds [x]
- `go test ./...` succeeds [x]
- the server starts locally [x]
- `GET /health` returns `200 OK` [x]

### Phase 2 - Backend MVP

Depends on: Phase 1

Goals:

- create the first real backend features
- support account creation and login
- persist users and messages
- support real-time messaging locally

Tasks:

- [x] Define the user model
- [x] Define the message model
- [x] Define the reaction model
- [x] Create repository interfaces in `pkg/store`
- [ ] Add database setup in `pkg/db`
- [ ] Add password hashing with bcrypt
- [ ] Add JWT creation and validation
- [ ] Add signup endpoint
- [ ] Add login endpoint
- [ ] Add authenticated profile or session-check endpoint
- [ ] Add message persistence
- [ ] Create the WebSocket hub
- [ ] Broadcast messages to connected clients
- [ ] Run `go fmt ./...`

Done when:

- `go build ./...` succeeds
- `go test ./...` succeeds
- a user can sign up locally
- a user can log in locally
- two local clients can authenticate and exchange messages

### Phase 3 - Frontend MVP

Depends on: Phase 2 API shape being stable enough to target

Goals:

- create the frontend project
- build the first usable UI
- model the screens needed for auth and chat

Tasks:

- [ ] Scaffold a Vite + React + TypeScript app
- [ ] Add TailwindCSS
- [ ] Add shadcn/ui
- [ ] Set up the app layout and routing
- [ ] Build the landing or login entry view
- [ ] Build the signup view
- [ ] Build the authenticated chat shell
- [ ] Build the room or conversation list UI
- [ ] Build the message list UI
- [ ] Build the message composer
- [ ] Add mocked data so the UI can be tested before integration

Done when:

- the frontend runs locally
- the login and chat screens render correctly
- the UI can be navigated without backend integration
- mocked messages display in the chat interface

### Phase 4 - Full-Stack Integration

Depends on: Phases 2 and 3

Goals:

- connect the frontend to the live backend
- use real auth tokens
- replace mocked chat data with live API and WebSocket data

Tasks:

- [ ] Connect frontend signup to the backend
- [ ] Connect frontend login to the backend
- [ ] Store the JWT on the client
- [ ] Send authenticated REST requests with the JWT
- [ ] Open a WebSocket connection from the frontend
- [ ] Validate the user on WebSocket connect
- [ ] Render real messages from the backend
- [ ] Send real messages from the frontend
- [ ] Add media upload flow end to end
- [ ] Add reaction flow end to end
- [ ] Test failure cases and unauthenticated cases

Done when:

- the frontend works against the local backend end to end
- login persists for a usable session
- messages appear in real time across clients
- uploaded media is visible in chat
- reactions update correctly in the UI

### Phase 5 - Deployment And Hardening

Depends on: Phase 4

Goals:

- package the full application for the home lab
- secure the deployment
- confirm the app works for real users outside the local dev machine

Tasks:

- [ ] Add a Dockerfile for the backend
- [ ] Add a Docker strategy for the frontend build or static serving
- [ ] Add a `docker-compose.yml` stack
- [ ] Configure Postgres for deployment
- [ ] Configure environment variables for production
- [ ] Remove placeholder secrets
- [ ] Configure reverse proxy and TLS
- [ ] Lock down auth and session settings
- [ ] Confirm upload storage paths and permissions
- [ ] Add backup and recovery notes
- [ ] Test multi-user access from the home lab
- [ ] Review logs and basic security posture

Done when:

- the application is reachable from the home lab environment
- friends can log in securely
- HTTPS is enabled
- the app supports real chat usage with auth, messages, reactions, and media

## Verification Checklist

- [ ] `go build ./...`
- [ ] `go test ./...`
- [ ] backend starts locally
- [ ] frontend starts locally
- [ ] auth works end to end
- [ ] WebSocket messaging works end to end
- [ ] media upload works end to end
- [ ] deployment works on the home server

## Notes For Future Iteration

- Add OAuth only after the core username or email plus password flow is stable.
- Consider serving the frontend from the Go service only after the separate development workflow is working cleanly.
- Keep this file updated as the single checklist for project progress.
