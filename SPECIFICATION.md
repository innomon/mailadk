# Specification: mailadk (Mail Identity Verifier)

## 1. Overview
mailadk is a Go-based service that acts as an identity bridge between email and the agentic ecosystem (whatsadk, agentic). Its primary role is to verify a user's email address and link it to their WhatsApp Logical Identity (LIP ID).

## 2. Core Responsibilities
- **Inbox Polling**: Periodically connect to the configured IMAP server to fetch new messages.
- **Verification**: Extract a Base32-encoded JWT (signed by the SPA) from the email body.
- **Identity Merger**: Extract the sender's verified email address from the email headers.
- **JWT Issuance**: Create a new 'Final JWT' containing both the lip_id and the email.
- **Callback**: POST the Final JWT back to the SPA's /auth/callback endpoint.

## 3. Technology Stack
- **Language**: Go 1.22+
- **IMAP Library**: github.com/emersion/go-imap/v2 (for modern, performant email retrieval).
- **JWT Library**: github.com/golang-jwt/jwt/v5
- **Configuration**: Priority-based loading:
  1. Environment variables (prefixed with `MAILADK_`).
  2. `config.yaml` in the root or `./config/` directory.
  3. `.env` file for local development.
  4. Default values (e.g., `poll_interval: 30s`).

## 4. Authentication & Security
- **Inbound JWT Verification**: mailadk must verify the SPA's signature. This requires the SPA's public key.
- **Outbound JWT Signing**: mailadk signs the Final JWT using a private key (RS256 or EdDSA) trusted by whatsadk and the SPA.
- **Token Claims**:
  - lip_id: The user's WhatsApp LIP.
  - email: The verified sender's email.
  - sub: lip_id:email
  - iat, exp: Standard claims.

## 5. Flow
1. **Connect**: Use imapclient.DialTLS to connect to the mail server.
2. **Poll**: Every N seconds, fetch messages from the INBOX with the configured subject.
3. **Parse**:
   - Extract the From header to get the email address.
   - Parse the body for the {"auth": "..."} JSON.
4. **Verify**: Decode Base32 and verify the signature of the SPA-generated JWT.
5. **Post**: Use net/http to send a JSON payload {"token": "..."} to the SPA callback URL.
6. **Cleanup**: Mark the email as Deleted or Seen to avoid re-processing.
