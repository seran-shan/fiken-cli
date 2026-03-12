# fiken-cli

A command-line client for the [Fiken.no](https://fiken.no) accounting API. Manage your Norwegian business accounting from the terminal.

## Features

- 🏢 List and manage companies
- 📊 Chart of accounts and balances
- 🛒 Full purchase management including drafts lifecycle and payments
- 👥 Full contact management (list, create, get, update, delete)
- 📦 Full product management (list, create, get, update, delete)
- 💰 Create and manage sales, sale drafts, attachments, settlements, and sale payments
- 📄 Create and manage invoices, invoice drafts, and attachments
- 📒 Journal entries with listing, get, and attachments
- 🔄 View and delete transactions
- 📥 EHF inbox management and document upload
- 🏦 Bank account management and balances
- 👤 User info and company details
- 📊 Product sales reports
- 👥 Contact persons and groups management
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
fiken companies get <slug>  # Get company details by slug
fiken companies default     # Show default company
fiken companies default <slug>  # Set default company
```

### Accounts

```bash
fiken accounts              # List chart of accounts
fiken accounts --from 1000 --to 2000  # Filter by account range
fiken accounts get 1920     # Get a specific account
```

### Balances

```bash
fiken balances                   # List account balances (date defaults to today)
fiken balances --date 2025-01-15  # List balances for a specific date
fiken balances get 1920          # Get balance for a specific account
```

### Bank Accounts

```bash
fiken bank list                  # List bank accounts
fiken bank get 123               # Get a specific bank account
fiken bank create --name "Main Account" --bank-account-number 12345678901
fiken bank create --name "Savings" --bank-account-number 98765432109 --type tax
fiken bank balances              # List bank balances
fiken bank balances --date 2025-01-15  # List balances for a specific date
```

### Inbox (EHF)

```bash
fiken inbox                 # List all inbox documents
fiken inbox --status pending     # Filter by status
fiken inbox get 123         # Get a specific inbox document
fiken inbox upload --file <path>  # Upload document to inbox
fiken inbox upload --file receipt.pdf --name "Office supplies" --description "Q1 office supplies"
```

### Purchases

```bash
fiken purchases list             # List purchases
fiken purchases get 123          # Get a specific purchase
fiken purchases delete 123 --description "Created by mistake"  # Soft-delete
fiken purchases create \         # Create a new purchase
  --date 2025-01-15 \
  --kind cash_purchase \
  --description "Office supplies" \
  --account 6800 \
  --amount 1500.00 \
  --vat-type HIGH
fiken purchases attach --id 123 --file receipt.pdf
fiken purchases payments list 123
fiken purchases payments create 123 --date 2025-01-15 --account 1920 --amount 1500.00
fiken purchases payments get 123 456

# Purchase Drafts
fiken purchases drafts list
fiken purchases drafts create --cash=false --description "Office supplies" \
  --account 6800 --amount 1500.00 --vat-type HIGH --supplier-id 123
fiken purchases drafts get 123
fiken purchases drafts update 123 --amount 2000.00
fiken purchases drafts delete 123
fiken purchases drafts attach 123 --file receipt.pdf
fiken purchases drafts finalize 123
```

### Contacts

```bash
fiken contacts list
fiken contacts create --name "Acme AS" --customer --email post@acme.no
fiken contacts get 123
fiken contacts update 123 --phone "+47 999 99 999"
fiken contacts delete 123
fiken contacts attach 123 --file document.pdf

# Contact Persons
fiken contacts persons list 123
fiken contacts persons create 123 --name "John Doe" --email john@acme.no
fiken contacts persons get 123 --person-id 456
fiken contacts persons update 123 --person-id 456 --email newemail@acme.no
fiken contacts persons delete 123 --person-id 456
```

### Products

```bash
fiken products list
fiken products create --name "Consulting" --income-account 3000 --vat-type HIGH --unit-price 1250.00
fiken products get 123
fiken products update 123 --note "Updated product note"
fiken products delete 123
fiken products sales-report --from 2025-01-01 --to 2025-12-31
```

### Sales

```bash
fiken sales list            # List sales
fiken sales create --date 2025-01-15 --kind cash_sale --description "Consulting" \
  --account 3000 --amount 1250.00 --vat-type HIGH
fiken sales get 123
fiken sales delete 123 --description "Created by mistake"
fiken sales settle 123 --settled-date 2025-01-31
fiken sales attach --id 123 --file document.pdf  # Attach document to existing sale
fiken sales drafts list
fiken sales drafts create --cash=false --description "Consulting" --account 3000 --amount 1250.00 --vat-type HIGH
fiken sales drafts get 123
fiken sales drafts update 123 --amount 1500.00
fiken sales drafts delete 123
fiken sales drafts attach 123 --file document.pdf
fiken sales drafts finalize 123
fiken sales payments list 123
fiken sales payments create 123 --date 2025-01-15 --account 1920 --amount 1250.00
fiken sales payments get 123 456
```

### Invoices

```bash
fiken invoices list         # List invoices
fiken invoices create --issue-date 2025-01-15 --due-date 2025-01-29 --customer-id 123 \
  --bank-account-code 1920:12345 --description "Consulting" --unit-price 1250.00 --quantity 1 --vat-type HIGH
fiken invoices get 123
fiken invoices update 123 --new-due-date 2025-02-05
fiken invoices attach --id 123 --file document.pdf
fiken invoices send --invoice-id 123 --method auto
fiken invoices send --invoice-id 123 --method email,ehf --recipient-email client@example.com
fiken invoices counter                        # Get current counter
fiken invoices counter set --value 1000       # Set counter
fiken invoices attachments 123                # List invoice attachments
fiken invoices drafts list
fiken invoices drafts create --type invoice --customer-id 123 --days-until-due 14 \
  --description "Consulting" --quantity 1 --unit-price 1250.00 --vat-type HIGH
fiken invoices drafts get 123
fiken invoices drafts update 123 --quantity 2
fiken invoices drafts delete 123
fiken invoices drafts attach 123 --file document.pdf
fiken invoices drafts finalize 123
```

### Journal Entries

```bash
fiken journal list                 # List journal entries
fiken journal list --date 2025-01-15  # Filter by date
fiken journal get 123              # Get a specific journal entry
fiken journal attachments 123      # List attachments for a journal entry
fiken journal create \             # Create a general journal entry
  --date 2025-01-15 \
  --description "Year-end adjustment" \
  --debit-account 1920 \
  --credit-account 3000 \
  --amount 5000.00
fiken journal attach --id 123 --file document.pdf
```

### Transactions

```bash
fiken transactions list     # List all transactions
fiken transactions list --last-modified 2025-01-01  # Filter by date
fiken transactions get 456  # Get a specific transaction
fiken transactions delete 456 --description "Created by mistake"
```

### Groups

```bash
fiken groups                     # List contact groups
fiken groups --json              # Output as JSON
```

### User

```bash
fiken user                       # Show authenticated user info
fiken user --json                # Output as JSON
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
