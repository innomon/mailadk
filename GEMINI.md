# mailadk (Mail Identity Verifier)

## Overview
`mailadk` is a Go-based identity bridge that verifies user email addresses via IMAP and links them to WhatsApp Logical Identities (LIP ID). It acts as a trusted verifier in the agentic ecosystem.

## Core Mandates
- **Security First**: Never log or expose private keys, email credentials, or sensitive JWT claims.
- **Reliability**: Use robust IMAP polling with proper error handling and connection retries.
- **Verification Integrity**: Strictly verify the Single Page Application (SPA) signature before issuing a 'Final JWT'.
- **Clean Inbox**: Ensure processed emails are marked as seen or deleted to prevent double-processing.

## Tech Stack
- **Language**: Go 1.25+
- **IMAP**: `github.com/emersion/go-imap/v2`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **Configuration**: `spf13/viper` (YAML + Environment Variables), `joho/godotenv`
- **Logging**: Standard library `log` (for simplicity and visibility)

## Project Structure
- `cmd/mailadk/`: Entry point and main application logic.
- `internal/config/`: Configuration loading and validation.
- `internal/imap/`: IMAP polling and message retrieval logic.
- `internal/token/`: JWT verification and signing.
- `internal/callback/`: HTTP callback logic to the SPA.

## Development Workflow
- **Configuration**: Always support `config.yaml` with fallback to `.env` or environment variables (prefixed with `MAILADK_`).
- **Testing**: Include unit tests for token verification and Base32 decoding.
- **Documentation**: Keep `SPECIFICATION.md` and `IMPLEMENTATION_PLAN.md` updated with progress.
