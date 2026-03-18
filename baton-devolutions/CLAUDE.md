# CLAUDE.md

Instructions for AI assistants working with this Baton connector.

## What This Is

A ConductorOne Baton connector that syncs identity and access data from Devolutions Server (DVLS). Connectors implement the `ResourceSyncer` interface to expose users, groups, roles, and their relationships.

## Build & Test

```bash
go build ./cmd/baton-devolutions   # Build connector
go test ./...                      # Run tests
go test -v ./... -count=1          # Verbose, no cache
```

## Architecture

- **REST API Client** (`pkg/client/`): HTTP client for DVLS REST API with Application Identity auth (appKey + appSecret). Token auto-refreshes on expiry (5 min TTL).
- **Connector** (`pkg/connector/`): Resource syncers for Users, Groups, Roles, and Vaults.
- **Config** (`pkg/config/`): CLI configuration fields (server-url, app-key, app-secret).

## Resource Types

| Type | Trait | Description |
|------|-------|-------------|
| User | TRAIT_USER | DVLS users with email, username, status |
| Group | TRAIT_GROUP | User groups with membership |
| Role | TRAIT_ROLE | Permission sets (Contributor/Operator/Reader) |
| Vault | (none) | Vaults with permission-based access |

## Configuration

```bash
baton-devolutions --server-url="https://dvls.example.com" --app-key="..." --app-secret="..."
```

Or via environment variables:
- `BATON_SERVER_URL`
- `BATON_APP_KEY`
- `BATON_APP_SECRET`
