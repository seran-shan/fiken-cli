package api

import "time"

// PaginatedResponse wraps paginated list responses from Fiken.
type PaginatedResponse struct {
	Page        int `json:"page"`
	PageSize    int `json:"pageSize"`
	PageCount   int `json:"pageCount"`
	ResultCount int `json:"resultCount"`
}

// Company represents a Fiken company.
type Company struct {
	Name                string  `json:"name"`
	Slug                string  `json:"slug"`
	OrganizationNumber  string  `json:"organizationNumber"`
	VatType             string  `json:"vatType"`
	Address             Address `json:"address"`
	PhoneNumber         string  `json:"phoneNumber"`
	Email               string  `json:"email"`
	CreationDate        string  `json:"creationDate"`
	HasApiAccess        bool    `json:"hasApiAccess"`
	TestCompany         bool    `json:"testCompany"`
	AccountingStartDate string  `json:"accountingStartDate"`
}

type Address struct {
	StreetAddress      string `json:"streetAddress"`
	StreetAddressLine2 string `json:"streetAddressLine2,omitempty"`
	City               string `json:"city"`
	PostCode           string `json:"postCode"`
	Country            string `json:"country"`
}

type CompaniesResponse struct {
	PaginatedResponse
	Companies []Company `json:"companies"`
}

// Account represents an account in the chart of accounts.
type Account struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type AccountsResponse struct {
	PaginatedResponse
	Accounts []Account `json:"accounts"`
}

// AccountBalance represents a balance for an account.
type AccountBalance struct {
	Account Account `json:"account"`
	Balance int64   `json:"balance"`
}

type AccountBalancesResponse struct {
	PaginatedResponse
	AccountBalances []AccountBalance `json:"accountBalances"`
}

// BankAccount represents a bank account.
type BankAccount struct {
	BankAccountId     int64  `json:"bankAccountId"`
	Name              string `json:"name"`
	AccountCode       string `json:"accountCode"`
	BankAccountNumber string `json:"bankAccountNumber"`
	Iban              string `json:"iban,omitempty"`
	Bic               string `json:"bic,omitempty"`
	ForeignService    string `json:"foreignService,omitempty"`
	Type              string `json:"type"`
	Inactive          bool   `json:"inactive"`
}

type BankAccountsResponse struct {
	PaginatedResponse
	BankAccounts []BankAccount `json:"bankAccounts"`
}

// InboxDocument represents an item in the EHF inbox.
type InboxDocument struct {
	DocumentId  int64     `json:"documentId"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Filename    string    `json:"filename"`
	Status      string    `json:"status"`
	CreatedDate time.Time `json:"createdDate"`
}

type InboxResponse struct {
	PaginatedResponse
	Documents []InboxDocument `json:"documents"`
}

// Purchase represents a purchase/expense.
type Purchase struct {
	PurchaseId          int64       `json:"purchaseId"`
	TransactionId       int64       `json:"transactionId,omitempty"`
	Identifier          string      `json:"identifier,omitempty"`
	Date                string      `json:"date"`
	DueDate             string      `json:"dueDate,omitempty"`
	Kind                string      `json:"kind"`
	Lines               []OrderLine `json:"lines"`
	Supplier            Contact     `json:"supplier,omitempty"`
	Currency            string      `json:"currency"`
	PaymentAccount      string      `json:"paymentAccount,omitempty"`
	Paid                bool        `json:"paid"`
	TotalPaid           int64       `json:"totalPaid"`
	TotalPaidInCurrency int64       `json:"totalPaidInCurrency"`
}

type OrderLine struct {
	Description         string `json:"description"`
	Account             string `json:"account"`
	NetAmount           int64  `json:"netAmount"`
	VatAmount           int64  `json:"vatAmount"`
	GrossAmount         int64  `json:"grossAmount,omitempty"`
	NetAmountInCurrency int64  `json:"netAmountInCurrency,omitempty"`
	VatAmountInCurrency int64  `json:"vatAmountInCurrency,omitempty"`
	VatType             string `json:"vatType"`
}

// OrderLineRequest is used when creating purchases (write schema).
// Distinct from OrderLine which is used for reading existing purchases.
type OrderLineRequest struct {
	Description string `json:"description"`
	NetPrice    int64  `json:"netPrice"`
	Vat         int64  `json:"vat,omitempty"`
	Account     string `json:"account,omitempty"`
	VatType     string `json:"vatType"`
	ProjectId   int64  `json:"projectId,omitempty"`
}

type PurchasesResponse struct {
	PaginatedResponse
	Purchases []Purchase `json:"purchases"`
}

// PurchaseRequest is used to create a new purchase.
type PurchaseRequest struct {
	Date           string             `json:"date"`
	DueDate        string             `json:"dueDate,omitempty"`
	Kind           string             `json:"kind"`
	Paid           bool               `json:"paid"`
	Lines          []OrderLineRequest `json:"lines"`
	SupplierId     int64              `json:"supplierId,omitempty"`
	Currency       string             `json:"currency"`
	PaymentAccount string             `json:"paymentAccount,omitempty"`
	PaymentDate    string             `json:"paymentDate,omitempty"`
	Identifier     string             `json:"identifier,omitempty"`
	ProjectId      int64              `json:"projectId,omitempty"`
}

type ContactRef struct {
	ContactId       int64 `json:"contactId,omitempty"`
	ContactPersonId int64 `json:"contactPersonId,omitempty"`
}

// Sale represents a sale.
type Sale struct {
	SaleId    int64       `json:"saleId"`
	Date      string      `json:"date"`
	Kind      string      `json:"kind"`
	Lines     []OrderLine `json:"lines"`
	Customer  Contact     `json:"customer,omitempty"`
	Currency  string      `json:"currency"`
	DueDate   string      `json:"dueDate,omitempty"`
	Paid      bool        `json:"paid"`
	TotalPaid int64       `json:"totalPaid"`
}

type SalesResponse struct {
	PaginatedResponse
	Sales []Sale `json:"sales"`
}

// Invoice represents an invoice.
type Invoice struct {
	InvoiceId     int64       `json:"invoiceId"`
	InvoiceNumber int64       `json:"invoiceNumber"`
	IssueDate     string      `json:"issueDate"`
	DueDate       string      `json:"dueDate"`
	Lines         []OrderLine `json:"lines"`
	Customer      Contact     `json:"customer,omitempty"`
	Net           int64       `json:"net"`
	Vat           int64       `json:"vat"`
	Gross         int64       `json:"gross"`
	Currency      string      `json:"currency"`
	Paid          bool        `json:"paid"`
	Kid           string      `json:"kid,omitempty"`
}

type InvoicesResponse struct {
	PaginatedResponse
	Invoices []Invoice `json:"invoices"`
}

// JournalEntry represents a journal entry.
type JournalEntry struct {
	JournalEntryId int64         `json:"journalEntryId"`
	Date           string        `json:"date"`
	Description    string        `json:"description"`
	Lines          []JournalLine `json:"lines"`
}

type JournalLine struct {
	Account      string `json:"account"`
	DebitAmount  int64  `json:"debitAmount,omitempty"`
	CreditAmount int64  `json:"creditAmount,omitempty"`
}

// GeneralJournalEntryRequest is used to create a general journal entry (fri postering).
type GeneralJournalEntryRequest struct {
	Description    string                `json:"description,omitempty"`
	Open           bool                  `json:"open"`
	JournalEntries []JournalEntryRequest `json:"journalEntries"`
}

// JournalEntryRequest represents a single journal entry within a general journal entry.
type JournalEntryRequest struct {
	Description string                    `json:"description,omitempty"`
	Date        string                    `json:"date"`
	Lines       []JournalEntryLineRequest `json:"lines"`
}

// JournalEntryLineRequest is the write schema for journal entry lines.
// Amount means net amount for debit lines and gross amount (incl. VAT) for credit lines.
type JournalEntryLineRequest struct {
	Amount        int64  `json:"amount"`
	DebitAccount  string `json:"debitAccount,omitempty"`
	DebitVatCode  int64  `json:"debitVatCode,omitempty"`
	CreditAccount string `json:"creditAccount,omitempty"`
	CreditVatCode int64  `json:"creditVatCode,omitempty"`
}

type JournalEntriesResponse struct {
	PaginatedResponse
	JournalEntries []JournalEntry `json:"journalEntries"`
}

// Transaction represents a financial transaction.
type Transaction struct {
	TransactionId    int64  `json:"transactionId"`
	CreatedDate      string `json:"createdDate"`
	LastModifiedDate string `json:"lastModifiedDate"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	Deleted          bool   `json:"deleted"`
}

type TransactionsResponse struct {
	PaginatedResponse
	Transactions []Transaction `json:"transactions"`
}

// Contact represents a customer or supplier.
type Contact struct {
	ContactId          int64   `json:"contactId"`
	Name               string  `json:"name"`
	Email              string  `json:"email,omitempty"`
	OrganizationNumber string  `json:"organizationNumber,omitempty"`
	Customer           bool    `json:"customer"`
	Supplier           bool    `json:"supplier"`
	PhoneNumber        string  `json:"phoneNumber,omitempty"`
	MemberNumber       int64   `json:"memberNumber,omitempty"`
	Address            Address `json:"address,omitempty"`
	Language           string  `json:"language,omitempty"`
	Inactive           bool    `json:"inactive"`
}

type ContactsResponse struct {
	PaginatedResponse
	Contacts []Contact `json:"contacts"`
}
