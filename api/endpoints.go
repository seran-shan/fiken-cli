package api

const (
	BaseURL = "https://api.fiken.no/api/v2"

	// Company endpoints
	EndpointCompanies = "/companies"

	// Endpoints under /companies/{slug}
	EndpointAccounts        = "/companies/%s/accounts"
	EndpointAccountBalances = "/companies/%s/accountBalances"
	EndpointBankAccounts    = "/companies/%s/bankAccounts"
	EndpointInbox           = "/companies/%s/inbox"
	EndpointPurchases       = "/companies/%s/purchases"
	EndpointSales           = "/companies/%s/sales"
	EndpointInvoices        = "/companies/%s/invoices"
	EndpointJournalEntries  = "/companies/%s/journalEntries"
	EndpointTransactions    = "/companies/%s/transactions"
	EndpointContacts        = "/companies/%s/contacts"

	// General journal entry creation endpoint (different from EndpointJournalEntries which is for reading)
	EndpointGeneralJournalEntries = "/companies/%s/generalJournalEntries"

	// Attachment endpoints (use fmt.Sprintf with company slug and entity int64 ID)
	EndpointPurchaseAttachments     = "/companies/%s/purchases/%d/attachments"
	EndpointSaleAttachments         = "/companies/%s/sales/%d/attachments"
	EndpointInvoiceAttachments      = "/companies/%s/invoices/%d/attachments"
	EndpointJournalEntryAttachments = "/companies/%s/journalEntries/%d/attachments"

	// Single entity endpoints
	EndpointTransaction = "/companies/%s/transactions/%d"

	// Contact endpoints
	EndpointContact = "/companies/%s/contacts/%d"

	// Product endpoints
	EndpointProducts = "/companies/%s/products"
	EndpointProduct  = "/companies/%s/products/%d"

	// Invoice endpoints (single entity + drafts)
	EndpointInvoice                 = "/companies/%s/invoices/%d"
	EndpointInvoiceDrafts           = "/companies/%s/invoices/drafts"
	EndpointInvoiceDraft            = "/companies/%s/invoices/drafts/%d"
	EndpointInvoiceDraftAttachments = "/companies/%s/invoices/drafts/%d/attachments"
	EndpointInvoiceDraftFinalize    = "/companies/%s/invoices/drafts/%d/createInvoice"

	// Sale endpoints (single entity + delete/settle + payments + drafts)
	EndpointSale                 = "/companies/%s/sales/%d"
	EndpointSaleDelete           = "/companies/%s/sales/%d/delete"
	EndpointSaleSettle           = "/companies/%s/sales/%d/settled"
	EndpointSalePayments         = "/companies/%s/sales/%d/payments"
	EndpointSalePayment          = "/companies/%s/sales/%d/payments/%d"
	EndpointSaleDrafts           = "/companies/%s/sales/drafts"
	EndpointSaleDraft            = "/companies/%s/sales/drafts/%d"
	EndpointSaleDraftAttachments = "/companies/%s/sales/drafts/%d/attachments"
	EndpointSaleDraftFinalize    = "/companies/%s/sales/drafts/%d/createSale"

	// Purchase endpoints (single entity + payments)
	EndpointPurchase         = "/companies/%s/purchases/%d"
	EndpointPurchasePayments = "/companies/%s/purchases/%d/payments"
	EndpointPurchasePayment  = "/companies/%s/purchases/%d/payments/%d"
	EndpointPurchaseDelete           = "/companies/%s/purchases/%d/delete"
	EndpointPurchaseDrafts           = "/companies/%s/purchases/drafts"
	EndpointPurchaseDraft            = "/companies/%s/purchases/drafts/%d"
	EndpointPurchaseDraftAttachments = "/companies/%s/purchases/drafts/%d/attachments"
	EndpointPurchaseDraftFinalize    = "/companies/%s/purchases/drafts/%d/createPurchase"

	// Journal endpoints (additions)
	EndpointJournalEntry = "/companies/%s/journalEntries/%d"

	// Invoice endpoints (additions)
	EndpointInvoiceSend    = "/companies/%s/invoices/send"
	EndpointInvoiceCounter = "/companies/%s/invoices/counter"

	// Bank endpoints (additions)
	EndpointBankAccount  = "/companies/%s/bankAccounts/%d"
	EndpointBankBalances = "/companies/%s/bankBalances"

	// Transaction endpoints (additions)
	EndpointTransactionDelete = "/companies/%s/transactions/%d/delete"

	// Contact endpoints (additions)
	EndpointContactAttachments = "/companies/%s/contacts/%d/attachments"
	EndpointContactPersons     = "/companies/%s/contacts/%d/contactPerson"
	EndpointContactPerson      = "/companies/%s/contacts/%d/contactPerson/%d"

	// Groups
	EndpointGroups = "/companies/%s/groups"

	// Product endpoints (additions)
	EndpointProductSalesReport = "/companies/%s/products/salesReport"

	// Inbox endpoints (additions)
	EndpointInboxDocument = "/companies/%s/inbox/%d"

	// Company endpoints (additions)
	EndpointCompany = "/companies/%s"

	// User
	EndpointUser = "/user"

	// Account/Balance endpoints (additions)
	EndpointAccountBalance = "/companies/%s/accountBalances/%s"
	EndpointAccount        = "/companies/%s/accounts/%s"

	// Credit note endpoints
	EndpointCreditNotes                = "/companies/%s/creditNotes"
	EndpointCreditNote                 = "/companies/%s/creditNotes/%d"
	EndpointCreditNotesFull            = "/companies/%s/creditNotes/full"
	EndpointCreditNotesPartial         = "/companies/%s/creditNotes/partial"
	EndpointCreditNotesSend            = "/companies/%s/creditNotes/send"
	EndpointCreditNoteDrafts           = "/companies/%s/creditNotes/drafts"
	EndpointCreditNoteDraft            = "/companies/%s/creditNotes/drafts/%d"
	EndpointCreditNoteDraftAttachments = "/companies/%s/creditNotes/drafts/%d/attachments"
	EndpointCreditNoteDraftFinalize    = "/companies/%s/creditNotes/drafts/%d/createCreditNote"

	// Offer endpoints
	EndpointOffers                 = "/companies/%s/offers"
	EndpointOffer                  = "/companies/%s/offers/%d"
	EndpointOfferDrafts            = "/companies/%s/offers/drafts"
	EndpointOfferDraft             = "/companies/%s/offers/drafts/%d"
	EndpointOfferDraftAttachments  = "/companies/%s/offers/drafts/%d/attachments"
	EndpointOfferDraftFinalize     = "/companies/%s/offers/drafts/%d/createOffer"

	// Order confirmation endpoints
	EndpointOrderConfirmations                = "/companies/%s/orderConfirmations"
	EndpointOrderConfirmation                 = "/companies/%s/orderConfirmations/%d"
	EndpointOrderConfirmationDrafts           = "/companies/%s/orderConfirmations/drafts"
	EndpointOrderConfirmationDraft            = "/companies/%s/orderConfirmations/drafts/%d"
	EndpointOrderConfirmationDraftAttachments = "/companies/%s/orderConfirmations/drafts/%d/attachments"
	EndpointOrderConfirmationDraftFinalize    = "/companies/%s/orderConfirmations/drafts/%d/createOrderConfirmation"

	// Project endpoints
	EndpointProjects = "/companies/%s/projects"
	EndpointProject  = "/companies/%s/projects/%d"
)

// Pagination defaults
const (
	DefaultPageSize = 25
	MaxPageSize     = 100
)

// Response headers
const (
	HeaderPage        = "Fiken-Api-Page"
	HeaderPageSize    = "Fiken-Api-Page-Size"
	HeaderPageCount   = "Fiken-Api-Page-Count"
	HeaderResultCount = "Fiken-Api-Result-Count"
)
