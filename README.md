# mailadk

Mail Identity Verifier bridge for the agentic ecosystem.

## Features
- Periodically polls an IMAP inbox for verification emails.
- Extracts and verifies SPA-issued JWTs (Base32 encoded).
- Signs and issues a Final JWT with verified email and WhatsApp LIP ID.
- Posts the final token to the SPA callback.

## Quick Start
1. Install dependencies:
   ```bash
   go mod tidy
   ```
2. Configure the application via `config.yaml` or `.env`:
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```
3. Run the application:
   ```bash
   go run cmd/mailadk/main.go
   ```

## Documentation
- [SPECIFICATION.md](SPECIFICATION.md)
- [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)
- [GEMINI.md](GEMINI.md)
