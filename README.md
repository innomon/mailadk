# mailadk

Mail Identity Verifier and RAG Ingestion bridge for the agentic ecosystem.

## 🚀 Utilities

### 1. `mailadk` (Identity Verifier)
Periodically polls an IMAP inbox for verification emails containing SPA-issued JWTs.
- Verifies the SPA signature (Base32 encoded).
- Signs a 'Final JWT' linking the email and WhatsApp LIP ID.
- Posts the final token to the SPA callback URL.

**Run it:**
```bash
go run cmd/mailadk/main.go
```

### 2. `mail-puller` (RAG Ingestion)
Stateful utility to fetch emails and ingest them into a Dyna-SLM RAG store via MCP.
- Fetches all emails since a `start_date` (first run) or `last_fetched_time` (subsequent runs).
- Calls the `ingest_asset` tool on the `dyna-slm` MCP server.
- Automatically handles the MCP server lifecycle.

**Run it:**
```bash
go run cmd/mail-puller/main.go
```

## 🛠 Configuration
Managed via `config.yaml` or environment variables (prefixed with `MAILADK_`).

```yaml
imap:
  host: "imap.example.com:993"
  user: "user@example.com"
  pass: "password"
  subject: "Identity Verification"

puller:
  start_date: "2026-03-01"
  state_file: "./data/state.json"
  mcp_command: "go run /home/innomon/orez/dyna-slm/cmd/mcp-server/main.go"
  mcp_model: "dyna-gemma3-270m"
  temp_dir: "./data/temp"
```

## 📖 Documentation
- [SPECIFICATION.md](SPECIFICATION.md)
- [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)
- [GEMINI.md](GEMINI.md)
