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
- [x] Add database setup in `pkg/db`
  - [x] Add initial DB connection setup in `pkg/db`
  - [x] Create a docker-compose.yml file
  - [x] Build a persistent volume for the DB
  - [x] DB credentials and env setup
  - [x] Document startup steps for clonability
- [x] Create a schema for the DB
  - [x] Build tables for users, rooms, messages, and reactions
  - [x] Decide on SQL bootstrapping scripts or a migration approach
  - [x] Apply the initial schema successfully to local Postgres
- [x] Finalize UUID strategy for persistence
  - [x] Use UUID columns in Postgres for primary keys
  - [x] Decide to generate UUIDs in Go rather than in Postgres
  - [x] Add UUID generation to the write path for new records
- [x] Add DB-backed store implementations
  - [x] Create a Postgres-backed user store
    - [x] Create a user successfully
    - [x] Fetch that user by ID
    - [x] Fetch that user by username
    - [x] Ensure duplicate username returns the expected error
    - [x] Ensure missing users return the expected not-found error
  - [x] Create a Postgres-backed message store
    - [x] Create a message successfully
    - [x] Fetch messages by room ID
    - [x] Ensure messages are returned in a sensible chronological order
    - [x] Ensure missing room lookups return an empty result without crashing
  - [x] Create a Postgres-backed reaction store
    - [x] Create a reaction successfully
    - [x] Fetch reactions by message ID
    - [x] Ensure duplicate reactions are rejected by the unique constraint
    - [x] Ensure reactions are removed correctly
- [ ] Add password hashing with bcrypt
  - [x] Add the bcrypt dependency to the Go module
  - [x] Create a dedicated password helper package or file
  - [x] Add a function to hash a plain-text password
  - [x] Add a function to compare a plain-text password against a stored hash
  - [x] Ensure the hash output is never treated as reversible encryption
  - [x] Add tests for successful password hashing
  - [x] Add tests proving the same password does not produce a reusable plain-text value
  - [x] Add tests for successful password comparison
  - [x] Add tests for failed password comparison with the wrong password
  - [x] Decide where hashing should happen in the signup flow before persistence
- [x] Add JWT creation and validation
  - [x] Choose a JWT library for signing and parsing tokens
  - [x] Create a dedicated JWT helper file in `pkg/auth`
  - [x] Decide which claims Cabin Chat needs in the token
  - [x] Add a function to generate a signed access token for a user
  - [x] Add a function to parse and validate an incoming token string
  - [x] Ensure token validation fails for malformed tokens
  - [x] Ensure token validation fails when the signing secret is wrong
  - [x] Ensure token validation fails for expired tokens
  - [x] Add tests for successful token generation
  - [x] Add tests for successful token validation
  - [x] Add tests for invalid signature handling
  - [x] Add tests for malformed token handling
  - [x] Add tests for expired token handling
  - [x] Decide how authenticated user identity will be attached to request handling later
- [x] Add signup endpoint
  - [x] Create signup request and response types
  - [x] Add an HTTP handler for `POST /signup`
  - [x] Parse and validate the incoming JSON body
  - [x] Reject requests with missing username or password
  - [x] Hash the submitted password before building the user model
  - [x] Create a new user record through the user store
  - [x] Translate duplicate-username errors into the correct HTTP response
  - [x] Generate an access token for the newly created user
  - [x] Return a safe signup response without exposing the password hash
  - [x] Wire the signup handler into the router
  - [x] Add tests for successful signup
  - [x] Add tests for invalid request bodies
  - [x] Add tests for missing required fields
  - [x] Add tests for duplicate username handling
- [x] Add login endpoint
  - [x] Create login request and response types
  - [x] Add an HTTP handler for `POST /login`
  - [x] Parse and validate the incoming JSON body
  - [x] Reject requests with missing username or password
  - [x] Load the user by username from the user store
  - [x] Translate missing-user lookups into a generic invalid-credentials response
  - [x] Compare the submitted password against the stored password hash
  - [x] Translate password mismatch into the same generic invalid-credentials response
  - [x] Generate an access token for the authenticated user
  - [x] Return a safe login response without exposing the password hash
  - [x] Wire the login handler into the router
  - [x] Add tests for successful login
  - [x] Add tests for invalid request bodies
  - [x] Add tests for missing required fields
  - [x] Add tests for unknown username handling
  - [x] Add tests for wrong password handling
- [ ] Add authenticated profile or session-check endpoint
  - [x] Decide whether the endpoint will be `/me`, `/profile`, or `/session` (decided `/session`)
  - [x] Decide the response shape for the authenticated user payload (decided to mirror the safe public User shape closely)
  - [x] Create profile or session-check response types
  - [x] Create middleware to read the `Authorization` header
  - [x] Reject requests with a missing `Authorization` header
  - [x] Reject requests with an invalid `Authorization` header format
  - [x] Parse the bearer token from the header value
  - [x] Validate the JWT access token
  - [x] Reject requests with malformed, expired, or invalid tokens
  - [x] Extract the authenticated user ID from the validated token claims
  - [x] Attach the authenticated user ID to the request context
  - [ ] Add an HTTP handler for the authenticated profile or session-check endpoint
  - [ ] Load the authenticated user by ID from the user store
  - [ ] Translate missing-user lookups into the correct HTTP response
  - [ ] Return a safe authenticated-user response without exposing the password hash
  - [ ] Wire the auth middleware and endpoint into the router
  - [ ] Add tests for successful authenticated profile or session-check requests
  - [ ] Add tests for missing `Authorization` headers
  - [ ] Add tests for invalid `Authorization` header formats
  - [ ] Add tests for malformed tokens
  - [ ] Add tests for expired tokens
  - [ ] Add tests for invalid token signatures
  - [ ] Add tests for unknown authenticated user IDs
- [ ] Add message persistence
- [ ] Create the WebSocket hub
- [ ] Broadcast messages to connected clients
- [ ] Run `go fmt ./...`

Implementation notes:

- The app currently runs on the host machine while Postgres runs in Docker via Compose.
- The local database schema is now in place and applied successfully.
- UUID generation should happen in Go for now, not in Postgres. This keeps record creation logic simpler while the backend is still being built.
- Postgres-backed stores should be added incrementally, starting with users, then messages, then reactions.

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
