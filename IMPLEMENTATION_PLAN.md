# Implementation Plan: mailadk (Mail Identity Verifier)

## Checklist
- [ ] Initialize project with go mod init.
- [ ] Configure environment variables (IMAP settings, Auth settings, Public/Private keys).
- [ ] Implement Base32 decoder.
- [ ] Implement IMAP poller using emersion/go-imap/v2.
- [ ] Implement JWT signature verification (SPA-issued).
- [ ] Implement JWT signing (Final token issuance).
- [ ] Implement HTTP POST callback logic.
- [ ] Add graceful shutdown and logging.

## Phase 1: Configuration
- Set up .env or env vars for IMAP_HOST, IMAP_USER, IMAP_PASS, IMAP_SUBJECT.
- Configure SPA_PUBLIC_KEY and VERIFIER_PRIVATE_KEY.
- Set SPA_CALLBACK_URL.

## Phase 2: IMAP Client
- Implement a background routine that uses imapclient.DialTLS.
- Use Search() to filter for messages with the configured subject.
- Fetch message headers and body parts.

## Phase 3: Token Logic
- Parse body for auth JSON.
- Decode Base32 to get the original JWT string.
- Verify EdDSA signature with the SPA_PUBLIC_KEY.
- Create new claims: {"lip_id": "...", "email": "...", "sub": "lip:email", "iat": ..., "exp": ...}.
- Sign new token with RS256/EdDSA private key.

## Phase 4: Callback & Cleanup
- POST {"token": "..."} to SPA.
- Ensure the SPA callback is accessible (SW must be active in the target tab).
- Mark processed emails as Deleted to keep the inbox clean.
