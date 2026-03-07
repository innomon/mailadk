# mailadk (Mail Identity Verifier)

## Overview
`mailadk` is a Go-based identity bridge and ingestion suite that verifies email addresses via IMAP and ingests content into a RAG store. It acts as a trusted verifier in the agentic ecosystem.

## Core Mandates
- **Security First**: Never log or expose private keys, email credentials, or sensitive JWT claims.
- **Reliability**: Use robust IMAP polling with proper error handling and connection retries.
- **Verification Integrity**: Strictly verify the SPA signature before issuing a 'Final JWT'.
- **Clean Inbox**: Ensure processed emails are marked as seen or deleted to prevent double-processing.

## Tech Stack
- **Language**: Go 1.25+
- **IMAP**: `github.com/emersion/go-imap/v2`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **MCP**: `github.com/modelcontextprotocol/go-sdk/mcp`
- **Configuration**: `spf13/viper` (YAML + Env Vars), `joho/godotenv`

## Project Structure
### Binaries
- `cmd/mailadk/`: Identity verifier daemon (POLL -> VERIFY -> POST).
- `cmd/mail-puller/`: Stateful email ingestion utility (SYNC -> INGEST).

### Internal
- `internal/config/`: Configuration loading and validation.
- `internal/imap/`: IMAP polling and message retrieval logic.
- `internal/token/`: JWT verification and signing.
- `internal/callback/`: HTTP callback logic to the SPA.
- `internal/mcp/`: MCP client implementation to interact with Dyna-SLM.
- `internal/storage/`: State management and temporary file handling.

## Development Workflow
- **Configuration**: Support `config.yaml` with fallback to `.env` or environment variables (prefixed with `MAILADK_`).
- **MCP Integration**: `mail-puller` should ingest via the `ingest_asset` tool of the Dyna-SLM server.
- **Testing**: Include unit tests for core logic (tokens, base32).
- **Documentation**: Keep `SPECIFICATION.md` and `IMPLEMENTATION_PLAN.md` updated.
