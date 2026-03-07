# Implementation Plan: mailadk (Mail Identity Verifier)

## Checklist
- [x] Initialize project with go mod init.
- [x] Configure environment variables (IMAP settings, Auth settings, Public/Private keys).
- [ ] Implement Base32 decoder.
- [ ] Implement IMAP poller using emersion/go-imap/v2.
- [ ] Implement JWT signature verification (SPA-issued).
- [ ] Implement JWT signing (Final token issuance).
- [ ] Implement HTTP POST callback logic.
- [ ] Add graceful shutdown and logging.

## Phase 1: Configuration (Completed)
- [x] Configure environment variables (IMAP settings, Auth settings, Public/Private keys).
- [x] Set up `config.yaml` support with fallback to `.env`.
- [x] Create basic `main.go` entry point.

## Phase 2: Base32 & Token Logic
- [ ] Implement Base32 decoder for the SPA-issued JWT.
- [ ] Implement JWT signature verification (SPA-issued).
- [ ] Implement JWT signing (Final token issuance).

## Phase 3: IMAP Client
- [ ] Implement a background routine that uses `imapclient.DialTLS`.
- [ ] Use `Search()` to filter for messages with the configured subject.
- [ ] Fetch message headers and body parts.

## Phase 4: Callback & Cleanup
- [ ] Implement HTTP POST callback logic to SPA.
- [ ] Ensure the SPA callback is accessible.
- [ ] Mark processed emails as Deleted or Seen to keep the inbox clean.
- [ ] Add graceful shutdown and logging.
