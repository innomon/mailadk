In the Go ecosystem, "popular and secure" is often split between libraries that **send** (SMTP/API) and those that **receive** (IMAP/POP3). Because the standard library's `net/smtp` is intentionally frozen and lacks modern security features like implicit SSL, the community has coalesced around a few key open-source projects.

### 1. For Sending: go-mail (`wneessen/go-mail`)

Currently considered the modern successor to older libraries like `gomail`. It is actively maintained and designed with security and Go idioms in mind.

* **Why it’s popular:** It has a very small dependency footprint and supports modern authentication (XOAUTH2, SCRAM, etc.).
* **Security:** Native support for implicit SSL/TLS and explicit STARTTLS. It handles RFC-compliant address validation out of the box.
* **Key Feature:** Middleware support, allowing you to easily add DKIM signing or logging without hacking the core library.

### 2. For Receiving: go-imap (`emersion/go-imap`)

If you need to fetch or manage emails, the `emersion` suite of tools is the industry standard for Go.

* **Why it’s popular:** It is the most robust implementation of IMAP4rev1 in Go. It powers many high-performance mail clients and servers (like Proton Mail's bridge).
* **Security:** Strongly typed, context-aware, and supports SASL authentication.
* **Version Note:** As of 2026, **v2** is the recommended version for new projects, offering better performance and a cleaner API.

---

### Comparison Table: Top Open-Source Options

| Library | Primary Use | Strengths | Maintenance |
| --- | --- | --- | --- |
| **`wneessen/go-mail`** | Sending | Modern, secure defaults, zero-dependency feel. | **High** |
| **`emersion/go-imap`** | Receiving | Industry standard, supports IMAP v4rev1 & v2. | **High** |
| **`emersion/go-smtp`** | Sending | Use this if you are building an SMTP *server* or a very low-level client. | High |
| **`jordan-wright/email`** | Sending | Very simple API; great for quick scripts or simple attachments. | Moderate |
| **`go-gomail/gomail`** | Sending | Very famous, but often considered "legacy." Use `go-mail` instead. | Low |

---

### Quick Implementation Examples

#### Sending a Secure Email (`go-mail`)

```go
import "github.com/wneessen/go-mail"

func main() {
    m := mail.NewMsg()
    m.From("sender@example.com")
    m.To("recipient@example.com")
    m.Subject("Gopher Security Update")
    m.SetBodyString(mail.TextPlain, "The system is secure.")

    // NewClient handles the security handshake automatically
    c, _ := mail.NewClient("smtp.example.com", mail.WithPort(465), mail.WithSMTPAuth(mail.SMTPAuthPlain),
        mail.WithUsername("user"), mail.WithPassword("pass"))
    
    c.DialAndSend(m)
}

```

#### Receiving Emails (`go-imap/v2`)

```go
import "github.com/emersion/go-imap/v2/protocols/imapclient"

func main() {
    // Connect to the server using TLS
    c, _ := imapclient.DialTLS("imap.example.com:993", nil)
    defer c.Logout()

    // Authenticate and list the first 10 messages in INBOX
    c.Login("user", "pass").Wait()
    // ... further mailbox operations
}

```

### Recommendation

* **For most apps:** Use **`wneessen/go-mail`** for sending and **`emersion/go-imap/v2`** for receiving. They represent the best balance of modern security and active community support.

