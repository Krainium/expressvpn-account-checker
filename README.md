# 🔐 ExpressVPN Account Checker


A fast, terminal-based ExpressVPN credential checker. Point it at a combo list, set how many workers you want, let it run. Every valid account gets its full details extracted — plan, expiry, OVPN credentials, payment method — saved to disk the moment it's found.

---

## ✨ What It Does

- ✅ Checks ExpressVPN accounts against the live API
- 🏷️ Classifies each account — **Premium**, **Trial**, **Free**, **Dead**
- 📋 Extracts plan type, expiry date, days remaining, auto-renew status, payment method
- 🔑 Pulls OVPN + PPTP credentials from every valid account
- 📄 Generates a ready-to-import `.ovpn` config file per hit in `ovpns/`
- 💾 Saves all hits instantly to `ExpressVpnHits.txt` as they're found
- ⚡ Runs up to 20 workers in parallel for fast bulk checking
- 🔄 Auto-retries on rate limits — no babysitting required
- 📊 Live progress bar showing Premium / Trial / Free / Dead counts in real time
- 🌍 Built-in server browser — 100+ ExpressVPN servers across 90+ countries to pick for OVPN generation

---

## 🖥️ Menu

```
  [1]  Single Check
  [2]  File Check
  [3]  Browse Servers
  [4]  Exit
```

---

## 🚀 Usage

### Single Account

Run the tool, pick `[1]`, enter the email and password. The result prints immediately.

### Bulk Combo List

Pick `[2]`, provide the path to your combo file, choose worker count (1–20). The checker processes every line concurrently while the live bar tracks progress.

**Combo format — one account per line:**
```
email@example.com:password123
another@email.com:hunter2
```

### Server Browser

Pick `[3]` to browse the full server list by country. Selecting a server means every OVPN config generated in that session will use that server's hostname.

---

## 📊 Output

**Terminal — premium hit:**
```
[PREMIUM] user@example.com:password
         Status    : PREMIUM / PAID
         Plan      : 12 Month
         Expires   : 2026-12-01 (224 days)
         Auto-Renew: true
         Payment   : ANDROID
         OVPN      : xv-username / xv-password
         Config    : ovpns/user@example.com.ovpn
```

**Terminal — trial hit:**
```
[TRIAL]   user@example.com:password
         Status    : FREE TRIAL
         Expires   : 2025-06-15 (12 days)
```

**Live progress bar during bulk run:**
```
[████████░░░░░░░░░░░░] 412/900 (45%)  Premium 18 | Trial 4 | Free 23 | Dead 367 | Retry 2
```

---

## 📁 Output Files

| File | What's in it |
|------|--------------|
| `ExpressVpnHits.txt` | Every valid account, one per line, with all extracted fields |
| `ovpns/<email>.ovpn` | Ready-to-use OpenVPN config for that account |

**Hit line format:**
```
email:pass | Status = PREMIUM | OVPNUsername = ... | OVPNPassword = ... | Plan = 12 Month | ExpireDate = 2026-12-01 | DaysLeft = 224 | AutoRenew = true | PaymentMethod = ANDROID
```

**Generated `.ovpn` includes:**
- Account email, plan, expiry embedded as comments
- OVPN username/password in `<auth-user-pass>` block
- Full ExpressVPN CA3 certificate
- UDP + TCP remotes with random selection fallback

---

## 🔧 Build

Requires **Go 1.21+**. No external dependencies — pure stdlib.

```bash
git clone https://github.com/krainium/expressvpn-account-checker
cd expressvpn-account-checker
go build -o main main.go
./main
```
```bash
or run directly
go run main.go

**Static binary (no CGO):**
```bash
CGO_ENABLED=0 go build -o expressvpn .
```

---

## 📦 Project Structure

```
.
├── main.go              # Full source — single file, no dependencies
├── go.mod               # Module file
├── ExpressVpnHits.txt   # Created automatically on first hit
└── ovpns/               # Created automatically — one .ovpn per hit
```

---

## ⚙️ Technical Notes

- Replicates the official iOS ExpressVPN client API flow (v21.21.0)
- Credentials are encrypted with RSA envelope encryption before transmission
- Subscription data comes from a batch endpoint and is parsed for license status
- Rate limit responses (HTTP 429) trigger automatic retry with 500ms back-off, up to 50 attempts per account
- Workers share a semaphore channel — no goroutine leaks
- All file writes are thread-safe

---

*krainium*
