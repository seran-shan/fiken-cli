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

// ContactRequest is the write schema for creating/updating contacts.
type ContactRequest struct {
	Name               string  `json:"name"`
	Email              string  `json:"email,omitempty"`
	OrganizationNumber string  `json:"organizationNumber,omitempty"`
	PhoneNumber        string  `json:"phoneNumber,omitempty"`
	Customer           bool    `json:"customer"`
	Supplier           bool    `json:"supplier"`
	Language           string  `json:"language,omitempty"`
	MemberNumber       int64   `json:"memberNumber,omitempty"`
	Address            Address `json:"address,omitempty"`
	Inactive           bool    `json:"inactive,omitempty"`
}

// Product represents a product or service in the catalog.
type Product struct {
	ProductId     int64   `json:"productId"`
	Name          string  `json:"name"`
	UnitPrice     int64   `json:"unitPrice,omitempty"`
	IncomeAccount string  `json:"incomeAccount,omitempty"`
	VatType       string  `json:"vatType,omitempty"`
	Active        bool    `json:"active"`
	ProductNumber string  `json:"productNumber,omitempty"`
	Stock         float64 `json:"stock,omitempty"`
	Note          string  `json:"note,omitempty"`
}

type ProductsResponse struct {
	PaginatedResponse
	Products []Product `json:"products"`
}

// ProductRequest is used to create or update a product.
type ProductRequest struct {
	Name          string  `json:"name"`
	UnitPrice     int64   `json:"unitPrice,omitempty"`
	IncomeAccount string  `json:"incomeAccount,omitempty"`
	VatType       string  `json:"vatType,omitempty"`
	Active        bool    `json:"active"`
	ProductNumber string  `json:"productNumber,omitempty"`
	Stock         float64 `json:"stock,omitempty"`
	Note          string  `json:"note,omitempty"`
}

// InvoiceRequest is the write schema for creating invoices directly (not via draft).
type InvoiceRequest struct {
	IssueDate       string               `json:"issueDate"`
	DueDate         string               `json:"dueDate"`
	Lines           []InvoiceLineRequest `json:"lines"`
	CustomerId      int64                `json:"customerId,omitempty"`
	BankAccountCode string               `json:"bankAccountCode,omitempty"`
	Cash            bool                 `json:"cash"`
	OrderReference  string               `json:"orderReference,omitempty"`
	OurReference    string               `json:"ourReference,omitempty"`
	YourReference   string               `json:"yourReference,omitempty"`
}

// InvoiceLineRequest is a line item for direct invoice creation.
type InvoiceLineRequest struct {
	Description string `json:"description"`
	Quantity    int64  `json:"quantity"`
	UnitPrice   int64  `json:"unitPrice"`
	VatType     string `json:"vatType"`
	ProductId   int64  `json:"productId,omitempty"`
	Discount    int64  `json:"discount,omitempty"`
}

// UpdateInvoiceRequest is for PATCH updates to an existing invoice.
type UpdateInvoiceRequest struct {
	NewDueDate   string `json:"newDueDate,omitempty"`
	SentManually bool   `json:"sentManually,omitempty"`
}

// InvoiceDraftRequest is the write schema for creating/updating invoice drafts.
// Uses daysUntilDueDate (int), NOT dueDate (string) like direct invoice creation.
type InvoiceDraftRequest struct {
	Type              string                    `json:"type"`
	CustomerId        int64                     `json:"customerId,omitempty"`
	DaysUntilDueDate  int64                     `json:"daysUntilDueDate,omitempty"`
	BankAccountNumber string                    `json:"bankAccountNumber,omitempty"`
	Lines             []InvoiceDraftLineRequest `json:"lines,omitempty"`
	OurReference      string                    `json:"ourReference,omitempty"`
	YourReference     string                    `json:"yourReference,omitempty"`
	OrderReference    string                    `json:"orderReference,omitempty"`
}

// InvoiceDraftLineRequest is a line item for invoice drafts.
type InvoiceDraftLineRequest struct {
	Description string `json:"description"`
	Quantity    int64  `json:"quantity"`
	UnitPrice   int64  `json:"unitPrice"`
	VatType     string `json:"vatType"`
	ProductId   int64  `json:"productId,omitempty"`
	Discount    int64  `json:"discount,omitempty"`
}

// InvoiceDraft is the response type for an invoice draft.
type InvoiceDraft struct {
	DraftId          int64                     `json:"draftId"`
	Uuid             string                    `json:"uuid,omitempty"`
	Type             string                    `json:"type"`
	CustomerId       int64                     `json:"customerId,omitempty"`
	Lines            []InvoiceDraftLineRequest `json:"lines,omitempty"`
	Net              int64                     `json:"net,omitempty"`
	Vat              int64                     `json:"vat,omitempty"`
	Gross            int64                     `json:"gross,omitempty"`
	LastModifiedDate string                    `json:"lastModifiedDate,omitempty"`
}

type InvoiceDraftsResponse struct {
	PaginatedResponse
	Drafts []InvoiceDraft `json:"drafts"`
}

// SaleRequest is the write schema for creating sales directly (not via draft).
// Uses OrderLineRequest (same as purchases) for line items.
type SaleRequest struct {
	Date           string             `json:"date"`
	Kind           string             `json:"kind"`
	Lines          []OrderLineRequest `json:"lines"`
	Currency       string             `json:"currency"`
	Paid           bool               `json:"paid"`
	CustomerId     int64              `json:"customerId,omitempty"`
	PaymentDate    string             `json:"paymentDate,omitempty"`
	PaymentAccount string             `json:"paymentAccount,omitempty"`
	DueDate        string             `json:"dueDate,omitempty"`
}

// SaleDraftRequest is the write schema for creating/updating sale drafts.
// Uses DraftLineRequest (NOT OrderLineRequest) for lines.
type SaleDraftRequest struct {
	Cash       bool               `json:"cash"`
	Lines      []DraftLineRequest `json:"lines,omitempty"`
	CustomerId int64              `json:"contactId,omitempty"`
	Paid       bool               `json:"paid"`
	Currency   string             `json:"currency,omitempty"`
}

// DraftLineRequest is the line type for sale and purchase drafts (not invoices).
type DraftLineRequest struct {
	Text     string `json:"text"`
	Account  string `json:"account,omitempty"`
	VatType  string `json:"vatType,omitempty"`
	NetPrice int64  `json:"netPrice,omitempty"`
	Vat      int64  `json:"vat,omitempty"`
	Gross    int64  `json:"gross,omitempty"`
}

// SaleDraft is the response type for a sale draft.
type SaleDraft struct {
	DraftId    int64              `json:"draftId"`
	Uuid       string             `json:"uuid,omitempty"`
	Cash       bool               `json:"cash"`
	Lines      []DraftLineRequest `json:"lines,omitempty"`
	CustomerId int64              `json:"contactId,omitempty"`
	Paid       bool               `json:"paid"`
}

type SaleDraftsResponse struct {
	PaginatedResponse
	Drafts []SaleDraft `json:"drafts"`
}

// Payment represents a payment on a sale or purchase.
type Payment struct {
	PaymentId   int64  `json:"paymentId"`
	Date        string `json:"date"`
	Account     string `json:"account"`
	Amount      int64  `json:"amount"`
	AmountInNok int64  `json:"amountInNok,omitempty"`
	Currency    string `json:"currency,omitempty"`
	Fee         int64  `json:"fee,omitempty"`
}

type PaymentsResponse struct {
	PaginatedResponse
	Payments []Payment `json:"payments"`
}

// PaymentRequest is used to create a payment.
type PaymentRequest struct {
	Date     string `json:"date"`
	Account  string `json:"account"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency,omitempty"`
}

// --- Purchase draft types ---

// PurchaseDraftRequest is the write schema for creating/updating purchase drafts.
type PurchaseDraftRequest struct {
	Cash             bool               `json:"cash"`
	DueDate          string             `json:"dueDate,omitempty"`
	InvoiceNumber    string             `json:"invoiceNumber,omitempty"`
	Kid              string             `json:"kid,omitempty"`
	ProjectId        int64              `json:"projectId,omitempty"`
	InvoiceIssueDate string             `json:"invoiceIssueDate,omitempty"`
	Lines            []DraftLineRequest `json:"lines,omitempty"`
	SupplierId       int64              `json:"supplierId,omitempty"`
	Paid             bool               `json:"paid"`
	Currency         string             `json:"currency,omitempty"`
}

// PurchaseDraft is the response type for a purchase draft.
type PurchaseDraft struct {
	DraftId       int64              `json:"draftId"`
	Uuid          string             `json:"uuid,omitempty"`
	Cash          bool               `json:"cash"`
	Lines         []DraftLineRequest `json:"lines,omitempty"`
	SupplierId    int64              `json:"supplierId,omitempty"`
	DueDate       string             `json:"dueDate,omitempty"`
	InvoiceNumber string             `json:"invoiceNumber,omitempty"`
	Kid           string             `json:"kid,omitempty"`
	ProjectId     int64              `json:"projectId,omitempty"`
	Paid          bool               `json:"paid"`
}

type PurchaseDraftsResponse struct {
	PaginatedResponse
	Drafts []PurchaseDraft `json:"drafts"`
}

// --- Invoice management types ---

// SendInvoiceRequest is the request body for POST /invoices/send.
type SendInvoiceRequest struct {
	InvoiceId                  int64    `json:"invoiceId"`
	Method                     []string `json:"method,omitempty"`
	RecipientName              string   `json:"recipientName,omitempty"`
	RecipientEmail             string   `json:"recipientEmail,omitempty"`
	Message                    string   `json:"message,omitempty"`
	IncludeDocumentAttachments bool     `json:"includeDocumentAttachments,omitempty"`
}

// InvoiceCounter represents the invoice counter value.
type InvoiceCounter struct {
	Counter int64 `json:"counter"`
}

// --- Bank types ---

// BankAccountRequest is used to create a new bank account.
type BankAccountRequest struct {
	Name              string `json:"name"`
	BankAccountNumber string `json:"bankAccountNumber"`
	Type              string `json:"type,omitempty"`
	Bic               string `json:"bic,omitempty"`
	Iban              string `json:"iban,omitempty"`
}

// BankBalance represents a bank account balance at a given date.
type BankBalance struct {
	BankAccountId int64  `json:"bankAccountId"`
	Balance       int64  `json:"balance"`
	Date          string `json:"date,omitempty"`
	AccountCode   string `json:"accountCode,omitempty"`
}

type BankBalancesResponse struct {
	PaginatedResponse
	BankBalances []BankBalance `json:"bankBalances"`
}

// --- Contact person types ---

// ContactPerson represents a person associated with a contact.
type ContactPerson struct {
	ContactPersonId int64  `json:"contactPersonId"`
	Name            string `json:"name"`
	Email           string `json:"email,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
}

// ContactPersonRequest is used to create or update a contact person.
type ContactPersonRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// --- Attachment type ---

// Attachment represents an attachment on a resource.
type Attachment struct {
	AttachmentId int64  `json:"attachmentId"`
	Filename     string `json:"filename,omitempty"`
	Type         string `json:"type,omitempty"`
	Date         string `json:"date,omitempty"`
}

// --- Group type ---

// Group represents a contact group.
type Group struct {
	Name string `json:"name"`
}

// --- Product sales report types ---

// ProductSalesReportRequest is the request body for POST /products/salesReport.
type ProductSalesReportRequest struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	ProductIds []int64 `json:"productIds,omitempty"`
}

// ProductSalesReportResult contains the result of a product sales report.
type ProductSalesReportResult struct {
	Lines []ProductSalesLineInfo `json:"lines,omitempty"`
}

// ProductSalesLineInfo is a single line in a product sales report.
type ProductSalesLineInfo struct {
	ProductId   int64  `json:"productId"`
	ProductName string `json:"productName,omitempty"`
	Count       int64  `json:"count"`
	Net         int64  `json:"net"`
	Vat         int64  `json:"vat"`
	Gross       int64  `json:"gross"`
}

// --- User type ---

// UserInfo represents the authenticated user.
type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// --- Invoiceish draft types (shared for offer/order-confirmation/invoice drafts) ---

type InvoiceishDraftRequest struct {
	Type             string                    `json:"type"`
	InvoiceIssueDate string                    `json:"invoiceIssueDate,omitempty"`
	InvoiceNumber    string                    `json:"invoiceNumber,omitempty"`
	CustomerId       int64                     `json:"customerId,omitempty"`
	ContactPersonId  int64                     `json:"contactPersonId,omitempty"`
	Days             int64                     `json:"days,omitempty"`
	Hours            float64                   `json:"hours,omitempty"`
	Currency         string                    `json:"currency,omitempty"`
	BankAccountCode  string                    `json:"bankAccountCode,omitempty"`
	YourReference    string                    `json:"yourReference,omitempty"`
	OurReference     string                    `json:"ourReference,omitempty"`
	OrderReference   string                    `json:"orderReference,omitempty"`
	ProjectId        int64                     `json:"projectId,omitempty"`
	Lines            []InvoiceDraftLineRequest `json:"lines,omitempty"`
}

type InvoiceishDraftResult struct {
	DraftId          int64                     `json:"draftId"`
	Uuid             string                    `json:"uuid,omitempty"`
	Type             string                    `json:"type"`
	InvoiceIssueDate string                    `json:"invoiceIssueDate,omitempty"`
	InvoiceNumber    string                    `json:"invoiceNumber,omitempty"`
	CustomerId       int64                     `json:"customerId,omitempty"`
	ContactPersonId  int64                     `json:"contactPersonId,omitempty"`
	Days             int64                     `json:"days,omitempty"`
	Hours            float64                   `json:"hours,omitempty"`
	Currency         string                    `json:"currency,omitempty"`
	BankAccountCode  string                    `json:"bankAccountCode,omitempty"`
	YourReference    string                    `json:"yourReference,omitempty"`
	OurReference     string                    `json:"ourReference,omitempty"`
	OrderReference   string                    `json:"orderReference,omitempty"`
	ProjectId        int64                     `json:"projectId,omitempty"`
	Lines            []InvoiceDraftLineRequest `json:"lines,omitempty"`
	Net              int64                     `json:"net,omitempty"`
	Vat              int64                     `json:"vat,omitempty"`
	Gross            int64                     `json:"gross,omitempty"`
	LastModifiedDate string                    `json:"lastModifiedDate,omitempty"`
}

type InvoiceishDraftsResponse struct {
	PaginatedResponse
	Drafts []InvoiceishDraftResult `json:"drafts"`
}

// --- Credit note types ---

type CreditNote struct {
	CreditNoteId        int64       `json:"creditNoteId"`
	CreditNoteNumber    int64       `json:"creditNoteNumber,omitempty"`
	Kid                 string      `json:"kid,omitempty"`
	Customer            Contact     `json:"customer,omitempty"`
	Net                 int64       `json:"net"`
	Vat                 int64       `json:"vat"`
	Gross               int64       `json:"gross"`
	NetInNok            int64       `json:"netInNok"`
	VatInNok            int64       `json:"vatInNok"`
	GrossInNok          int64       `json:"grossInNok"`
	CreditNoteText      string      `json:"creditNoteText,omitempty"`
	YourReference       string      `json:"yourReference,omitempty"`
	OurReference        string      `json:"ourReference,omitempty"`
	OrderReference      string      `json:"orderReference,omitempty"`
	Address             Address     `json:"address,omitempty"`
	Lines               []OrderLine `json:"lines,omitempty"`
	Currency            string      `json:"currency"`
	IssueDate           string      `json:"issueDate"`
	Settled             bool        `json:"settled"`
	AssociatedInvoiceId int64       `json:"associatedInvoiceId,omitempty"`
}

type CreditNotesResponse struct {
	PaginatedResponse
	CreditNotes []CreditNote `json:"creditNotes"`
}

type FullCreditNoteRequest struct {
	IssueDate      string `json:"issueDate"`
	InvoiceId      int64  `json:"invoiceId"`
	CreditNoteText string `json:"creditNoteText,omitempty"`
}

type PartialCreditNoteRequest struct {
	IssueDate       string                  `json:"issueDate"`
	InvoiceId       int64                   `json:"invoiceId,omitempty"`
	ContactId       int64                   `json:"contactId,omitempty"`
	ContactPersonId int64                   `json:"contactPersonId,omitempty"`
	CreditNoteText  string                  `json:"creditNoteText,omitempty"`
	OurReference    string                  `json:"ourReference,omitempty"`
	YourReference   string                  `json:"yourReference,omitempty"`
	OrderReference  string                  `json:"orderReference,omitempty"`
	ProjectId       int64                   `json:"projectId,omitempty"`
	Currency        string                  `json:"currency,omitempty"`
	Lines           []CreditNoteLineRequest `json:"lines"`
}

type CreditNoteLineRequest struct {
	Description string `json:"description"`
	Account     string `json:"account"`
	VatCode     string `json:"vatCode"`
	Amount      int64  `json:"amount"`
	Quantity    int64  `json:"quantity,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

type SendCreditNoteRequest struct {
	CreditNoteId               int64    `json:"creditNoteId"`
	Method                     []string `json:"method,omitempty"`
	IncludeDocumentAttachments bool     `json:"includeDocumentAttachments,omitempty"`
	RecipientName              string   `json:"recipientName,omitempty"`
	RecipientEmail             string   `json:"recipientEmail,omitempty"`
	Message                    string   `json:"message,omitempty"`
	EmailSendOption            string   `json:"emailSendOption,omitempty"`
	OrganizationNumber         string   `json:"organizationNumber,omitempty"`
	MobileNumber               string   `json:"mobileNumber,omitempty"`
}

// --- Offer types ---

type Offer struct {
	OfferId     int64       `json:"offerId"`
	OfferNumber int64       `json:"offerNumber,omitempty"`
	Date        string      `json:"date,omitempty"`
	Net         int64       `json:"net"`
	Vat         int64       `json:"vat"`
	Gross       int64       `json:"gross"`
	Customer    Contact     `json:"customer,omitempty"`
	Lines       []OrderLine `json:"lines,omitempty"`
	Currency    string      `json:"currency,omitempty"`
}

type OffersResponse struct {
	PaginatedResponse
	Offers []Offer `json:"offers"`
}

// --- Order confirmation types ---

type OrderConfirmation struct {
	ConfirmationId int64       `json:"confirmationId"`
	Date           string      `json:"date,omitempty"`
	Net            int64       `json:"net"`
	Vat            int64       `json:"vat"`
	Gross          int64       `json:"gross"`
	Customer       Contact     `json:"customer,omitempty"`
	Lines          []OrderLine `json:"lines,omitempty"`
	Currency       string      `json:"currency,omitempty"`
}

type OrderConfirmationsResponse struct {
	PaginatedResponse
	OrderConfirmations []OrderConfirmation `json:"orderConfirmations"`
}

// --- Project types ---

type Project struct {
	ProjectId   int64   `json:"projectId"`
	Number      string  `json:"number,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	StartDate   string  `json:"startDate,omitempty"`
	EndDate     string  `json:"endDate,omitempty"`
	Contact     Contact `json:"contact,omitempty"`
	Completed   bool    `json:"completed"`
}

type ProjectsResponse struct {
	PaginatedResponse
	Projects []Project `json:"projects"`
}

type ProjectRequest struct {
	Number      string `json:"number,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"startDate,omitempty"`
	EndDate     string `json:"endDate,omitempty"`
	ContactId   int64  `json:"contactId,omitempty"`
	Completed   bool   `json:"completed"`
}
