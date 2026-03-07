# Implementation Plan: mailadk (Mail Identity Verifier)

## Checklist
- [x] Initialize project with go mod init.
- [x] Configure environment variables (IMAP settings, Auth settings, Public/Private keys).
- [x] Implement Base32 decoder.
- [x] Implement IMAP poller using emersion/go-imap/v2.
- [x] Implement JWT signature verification (SPA-issued).
- [x] Implement JWT signing (Final token issuance).
- [x] Implement HTTP POST callback logic.
- [x] Add graceful shutdown and logging.

## Phase 1: Configuration (Completed)
- [x] Configure environment variables (IMAP settings, Auth settings, Public/Private keys).
- [x] Set up `config.yaml` support with fallback to `.env`.
- [x] Create basic `main.go` entry point.

## Phase 2: Base32 & Token Logic (Completed)
- [x] Implement Base32 decoder for the SPA-issued JWT.
- [x] Implement JWT signature verification (SPA-issued).
- [x] Implement JWT signing (Final token issuance).

## Phase 3: IMAP Client (Completed)
- [x] Implement a background routine that uses `imapclient.DialTLS`.
- [x] Use `Search()` to filter for messages with the configured subject.
- [x] Fetch message headers and body parts.

## Phase 4: Callback & Cleanup (Completed)
- [x] Implement HTTP POST callback logic to SPA.
- [x] Ensure the SPA callback is accessible.
- [x] Mark processed emails as Deleted or Seen to keep the inbox clean.
- [x] Add graceful shutdown and logging.
