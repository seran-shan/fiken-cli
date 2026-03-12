# fiken-cli

A command-line client for the [Fiken.no](https://fiken.no) accounting API. Manage your Norwegian business accounting from the terminal.

## Features

- 🏢 List and manage companies
- 📊 Chart of accounts and balances
- 🛒 Create and manage purchases with receipt attachments
- 💰 Create and manage sales with document attachments
- 📄 Create invoices and attach supporting documents
- 📒 Create general journal entries with attachments
- 🔄 View transactions (auto-generated from purchases/sales/journal entries)
- 📥 EHF inbox management and document upload
- 🏦 Bank account overview
- 📋 Dashboard with key metrics
- 🔄 JSON output for scripting
- ⚡ Built-in rate limiting and pagination

## Installation

### From source

```bash
git clone https://github.com/seran-shan/fiken-cli.git
cd fiken-cli
make install
```

### Build locally

```bash
make build
./fiken --help
```

## Quick Start

### 1. Get your API token

Go to [Fiken API Settings](https://fiken.no/innstillinger/api) and create a Personal API Token.

### 2. Authenticate

```bash
fiken auth token <your-token>
```

### 3. List your companies

```bash
fiken companies
```

### 4. Set a default company (optional)

```bash
fiken companies default <company-slug>
```

### 5. View your dashboard

```bash
fiken status
```

## Command Reference

### Authentication

```bash
fiken auth token <token>    # Save API token
fiken auth token            # Show token status
fiken auth status           # Verify token works
fiken auth logout           # Remove stored token
```

### Companies

```bash
fiken companies             # List all companies
fiken companies default     # Show default company
fiken companies default <slug>  # Set default company
```

### Accounts

```bash
fiken accounts              # List chart of accounts
fiken accounts --from 1000 --to 2000  # Filter by account range
```

### Balances

```bash
fiken balances              # List account balances
```

### Bank Accounts

```bash
fiken bank list             # List bank accounts
```

### Inbox (EHF)

```bash
fiken inbox                 # List all inbox documents
fiken inbox --status pending    # Filter by status
fiken inbox upload --file <path>  # Upload document to inbox
fiken inbox upload --file receipt.pdf --name "Office supplies" --description "Q1 office supplies"
```

### Purchases

```bash
fiken purchases list        # List purchases
fiken purchases create \    # Create a new purchase
  --date 2025-01-15 \
  --kind cash_purchase \
  --account 1920 \
  --amount 1500.00 \
  --vat-type HIGH \
  --supplier "Office Supplies AS"
fiken purchases create \    # Create purchase with receipt attached
  --date 2025-01-15 \
  --kind cash_purchase \
  --account 1920 \
  --amount 1500.00 \
  --vat-type HIGH \
  --file receipt.pdf
fiken purchases attach --id 123 --file receipt.pdf  # Attach receipt to existing purchase
```

### Sales

```bash
fiken sales list            # List sales
fiken sales attach --id 123 --file document.pdf  # Attach document to existing sale
```

### Invoices

```bash
fiken invoices list         # List invoices
fiken invoices attach --id 123 --file document.pdf  # Attach document to existing invoice
```

### Journal Entries

```bash
fiken journal create \      # Create a general journal entry
  --date 2025-01-15 \
  --description "Year-end adjustment" \
  --debit-account 1920 \
  --credit-account 3000 \
  --amount 5000.00
fiken journal attach --id 123 --file document.pdf  # Attach document to journal entry
```

### Transactions

```bash
fiken transactions list     # List all transactions
fiken transactions list --last-modified 2025-01-01  # Filter by date
fiken transactions get 456  # Get a specific transaction
```

### Status Dashboard

```bash
fiken status                # Overview of pending items
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON (default: table) |
| `--no-input` | Non-interactive mode |
| `--company <slug>` | Select company (auto-detected if only one) |
| `--keyring-backend <backend>` | Keyring backend (default: `auto`) |

## Credential Storage

Credentials (API token, default company) are stored securely using your OS keyring via [99designs/keyring](https://github.com/99designs/keyring).

### Supported backends

| Backend | OS | Flag value |
|---------|----|------------|
| Secret Service (GNOME Keyring / KDE Wallet) | Linux | `secret-service` |
| Keychain | macOS | `keychain` |
| Windows Credential Manager | Windows | `wincred` |
| [pass](https://www.passwordstore.org/) | Linux/macOS | `pass` |
| Encrypted file | Any (fallback) | `file` |

By default (`auto`), the best available backend is used. The encrypted file backend is the last-resort fallback and will prompt for a password.

### Choosing a backend

```bash
# Use a specific backend
fiken --keyring-backend file auth token <token>

# Or set via environment variable
export FIKEN_KEYRING_BACKEND=file
fiken auth token <token>
```

### Migration from plaintext storage

If you previously stored your token in `~/.config/fiken/token`, it will be automatically migrated to the keyring on first use. The plaintext file is deleted after successful migration.

## API Details

- Base URL: `https://api.fiken.no/api/v2`
- Auth: Bearer token (Personal API Token)
- Amounts are in cents (øre): `100000` = `1 000,00 kr`
- Rate limit: max 4 requests/second (enforced by client)
- Pagination: automatic for large result sets
- Amounts in CLI are entered as decimals (e.g., `1500.00`) and auto-converted to cents

## Examples

### List companies as JSON

```bash
fiken companies --json
```

### Script: get all account codes

```bash
fiken accounts --json | jq '.[].code'
```

### Use with a specific company

```bash
fiken purchases list --company my-company-slug
```

### Quick status check

```bash
fiken status --json | jq '.inbox_count'
```

### Upload a receipt to inbox

```bash
fiken inbox upload --file receipt.pdf --name "Lunch receipt"
```

### Create a purchase with receipt

```bash
fiken purchases create --date 2025-01-15 --kind cash_purchase \
  --account 1920 --amount 250.00 --vat-type HIGH --file receipt.jpg
```

### Create a journal entry

```bash
fiken journal create --date 2025-01-15 --description "Depreciation" \
  --debit-account 6000 --credit-account 1200 --amount 10000.00
```

### Attach receipt to existing purchase

```bash
fiken purchases attach --id 12345 --file receipt.pdf
```

## Development

```bash
make build    # Build binary
make test     # Run tests
make fmt      # Format code
make lint     # Run linter
make clean    # Clean build artifacts
```

## License

MIT – see [LICENSE](LICENSE).
