package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var creditNotesCmd = &cobra.Command{
	Use:   "credit-notes",
	Short: "Manage credit notes",
}

var creditNotesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List credit notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		params := url.Values{}
		if v, _ := cmd.Flags().GetString("issue-date"); v != "" {
			params.Set("issueDate", v)
		}
		if cmd.Flags().Changed("settled") {
			settled, _ := cmd.Flags().GetBool("settled")
			params.Set("settled", strconv.FormatBool(settled))
		}
		if v, _ := cmd.Flags().GetInt64("customer-id"); v != 0 {
			params.Set("customerId", strconv.FormatInt(v, 10))
		}

		endpoint := fmt.Sprintf(api.EndpointCreditNotes, slug)
		creditNotes, err := FetchAllPages[api.CreditNote](client, endpoint, params, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching credit notes: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(creditNotes)
		}

		if len(creditNotes) == 0 {
			output.PrintInfo("No credit notes found.")
			return nil
		}

		table := output.NewTable("ID", "NUMBER", "ISSUE DATE", "CUSTOMER", "NET", "VAT", "GROSS", "CURRENCY", "SETTLED")
		for _, cn := range creditNotes {
			table.AddRow(
				fmt.Sprintf("%d", cn.CreditNoteId),
				fmt.Sprintf("%d", cn.CreditNoteNumber),
				cn.IssueDate,
				cn.Customer.Name,
				output.FormatAmount(cn.Net),
				output.FormatAmount(cn.Vat),
				output.FormatAmount(cn.Gross),
				cn.Currency,
				BoolToYesNo(cn.Settled),
			)
		}
		table.Print()
		fmt.Printf("\n%d credit notes\n", len(creditNotes))
		return nil
	},
}

var creditNotesGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a credit note by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var cn api.CreditNote
		endpoint := fmt.Sprintf(api.EndpointCreditNote, slug, id)
		if _, err := client.Get(endpoint, &cn); err != nil {
			return fmt.Errorf("fetching credit note: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(cn)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", cn.CreditNoteId))
		table.AddRow("Number", fmt.Sprintf("%d", cn.CreditNoteNumber))
		table.AddRow("Issue Date", cn.IssueDate)
		table.AddRow("Customer", cn.Customer.Name)
		table.AddRow("Net", output.FormatAmount(cn.Net))
		table.AddRow("VAT", output.FormatAmount(cn.Vat))
		table.AddRow("Gross", output.FormatAmount(cn.Gross))
		table.AddRow("Net (NOK)", output.FormatAmount(cn.NetInNok))
		table.AddRow("VAT (NOK)", output.FormatAmount(cn.VatInNok))
		table.AddRow("Gross (NOK)", output.FormatAmount(cn.GrossInNok))
		table.AddRow("Currency", cn.Currency)
		table.AddRow("Settled", BoolToYesNo(cn.Settled))
		table.AddRow("KID", cn.Kid)
		table.AddRow("Credit Note Text", cn.CreditNoteText)
		table.AddRow("Our Reference", cn.OurReference)
		table.AddRow("Your Reference", cn.YourReference)
		table.AddRow("Order Reference", cn.OrderReference)
		if cn.AssociatedInvoiceId != 0 {
			table.AddRow("Associated Invoice", fmt.Sprintf("%d", cn.AssociatedInvoiceId))
		}
		table.Print()

		if len(cn.Lines) > 0 {
			fmt.Println()
			lines := output.NewTable("DESCRIPTION", "ACCOUNT", "NET", "VAT", "VAT TYPE")
			for _, l := range cn.Lines {
				lines.AddRow(
					l.Description,
					l.Account,
					output.FormatAmount(l.NetAmount),
					output.FormatAmount(l.VatAmount),
					l.VatType,
				)
			}
			lines.Print()
		}
		return nil
	},
}

var creditNotesCreateFullCmd = &cobra.Command{
	Use:   "create-full",
	Short: "Create a full credit note from an invoice",
	RunE: func(cmd *cobra.Command, args []string) error {
		issueDate, _ := cmd.Flags().GetString("issue-date")
		invoiceID, _ := cmd.Flags().GetInt64("invoice-id")
		creditNoteText, _ := cmd.Flags().GetString("credit-note-text")

		var missing []string
		if issueDate == "" {
			missing = append(missing, "--issue-date")
		}
		if invoiceID == 0 {
			missing = append(missing, "--invoice-id")
		}
		if len(missing) > 0 {
			return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		req := api.FullCreditNoteRequest{
			IssueDate:      issueDate,
			InvoiceId:      invoiceID,
			CreditNoteText: creditNoteText,
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointCreditNotesFull, slug), req)
		if err != nil {
			return fmt.Errorf("creating full credit note: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing credit note ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Full credit note created (ID: %d)", id))
		return nil
	},
}

var creditNotesCreatePartialCmd = &cobra.Command{
	Use:   "create-partial",
	Short: "Create a partial credit note",
	RunE: func(cmd *cobra.Command, args []string) error {
		issueDate, _ := cmd.Flags().GetString("issue-date")
		description, _ := cmd.Flags().GetString("description")
		account, _ := cmd.Flags().GetString("account")
		vatCode, _ := cmd.Flags().GetString("vat-code")
		amountStr, _ := cmd.Flags().GetString("amount")

		var missing []string
		if issueDate == "" {
			missing = append(missing, "--issue-date")
		}
		if description == "" {
			missing = append(missing, "--description")
		}
		if account == "" {
			missing = append(missing, "--account")
		}
		if vatCode == "" {
			missing = append(missing, "--vat-code")
		}
		if amountStr == "" {
			missing = append(missing, "--amount")
		}
		if len(missing) > 0 {
			return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
		}

		amount, err := ParseAmountToCents(amountStr)
		if err != nil {
			return err
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		invoiceID, _ := cmd.Flags().GetInt64("invoice-id")
		contactID, _ := cmd.Flags().GetInt64("contact-id")
		contactPersonID, _ := cmd.Flags().GetInt64("contact-person-id")
		creditNoteText, _ := cmd.Flags().GetString("credit-note-text")
		ourReference, _ := cmd.Flags().GetString("our-reference")
		yourReference, _ := cmd.Flags().GetString("your-reference")
		orderReference, _ := cmd.Flags().GetString("order-reference")
		projectID, _ := cmd.Flags().GetInt64("project-id")
		currency, _ := cmd.Flags().GetString("currency")
		quantity, _ := cmd.Flags().GetInt64("quantity")
		comment, _ := cmd.Flags().GetString("comment")

		line := api.CreditNoteLineRequest{
			Description: description,
			Account:     account,
			VatCode:     vatCode,
			Amount:      amount,
			Quantity:    quantity,
			Comment:     comment,
		}

		req := api.PartialCreditNoteRequest{
			IssueDate:       issueDate,
			InvoiceId:       invoiceID,
			ContactId:       contactID,
			ContactPersonId: contactPersonID,
			CreditNoteText:  creditNoteText,
			OurReference:    ourReference,
			YourReference:   yourReference,
			OrderReference:  orderReference,
			ProjectId:       projectID,
			Currency:        currency,
			Lines:           []api.CreditNoteLineRequest{line},
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointCreditNotesPartial, slug), req)
		if err != nil {
			return fmt.Errorf("creating partial credit note: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing credit note ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Partial credit note created (ID: %d)", id))
		return nil
	},
}

var creditNotesSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a credit note",
	RunE: func(cmd *cobra.Command, args []string) error {
		creditNoteID, _ := cmd.Flags().GetInt64("credit-note-id")
		if creditNoteID == 0 {
			return fmt.Errorf("missing required flag: --credit-note-id")
		}

		methodStr, _ := cmd.Flags().GetString("method")
		recipientName, _ := cmd.Flags().GetString("recipient-name")
		recipientEmail, _ := cmd.Flags().GetString("recipient-email")
		message, _ := cmd.Flags().GetString("message")
		includeAttachments, _ := cmd.Flags().GetBool("include-document-attachments")
		emailSendOption, _ := cmd.Flags().GetString("email-send-option")
		orgNumber, _ := cmd.Flags().GetString("organization-number")
		mobileNumber, _ := cmd.Flags().GetString("mobile-number")

		methods := strings.Split(methodStr, ",")

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		req := api.SendCreditNoteRequest{
			CreditNoteId:               creditNoteID,
			Method:                     methods,
			RecipientName:              recipientName,
			RecipientEmail:             recipientEmail,
			Message:                    message,
			IncludeDocumentAttachments: includeAttachments,
			EmailSendOption:            emailSendOption,
			OrganizationNumber:         orgNumber,
			MobileNumber:               mobileNumber,
		}

		if err := client.Post(fmt.Sprintf(api.EndpointCreditNotesSend, slug), req, nil); err != nil {
			return fmt.Errorf("sending credit note: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Credit note %d sent", creditNoteID))
		return nil
	},
}

var creditNotesCounterCmd = &cobra.Command{
	Use:   "counter",
	Short: "Get or set the credit note counter",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}
		var counter api.CreditNoteCounter
		if _, err := client.Get(fmt.Sprintf(api.EndpointCreditNoteCounter, slug), &counter); err != nil {
			return fmt.Errorf("fetching credit note counter: %w", err)
		}
		if jsonOutput {
			return output.PrintJSON(counter)
		}
		fmt.Printf("Credit note counter: %d\n", counter.Value)
		return nil
	},
}

var creditNotesCounterSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the credit note counter",
	RunE: func(cmd *cobra.Command, args []string) error {
		value, _ := cmd.Flags().GetInt64("value")
		client, err := getClient()
		if err != nil {
			return err
		}
		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}
		req := api.CreditNoteCounter{Value: int32(value)}
		if err := client.Post(fmt.Sprintf(api.EndpointCreditNoteCounter, slug), req, nil); err != nil {
			return fmt.Errorf("setting credit note counter: %w", err)
		}
		output.PrintSuccess(fmt.Sprintf("Credit note counter set to %d", value))
		return nil
	},
}

func init() {
	creditNotesListCmd.Flags().String("issue-date", "", "Filter by issue date (YYYY-MM-DD)")
	creditNotesListCmd.Flags().Bool("settled", false, "Filter by settled status")
	creditNotesListCmd.Flags().Int64("customer-id", 0, "Filter by customer ID")

	creditNotesCreateFullCmd.Flags().String("issue-date", "", "Issue date (YYYY-MM-DD, required)")
	creditNotesCreateFullCmd.Flags().Int64("invoice-id", 0, "Invoice ID to credit (required)")
	creditNotesCreateFullCmd.Flags().String("credit-note-text", "", "Credit note text (optional)")

	creditNotesCreatePartialCmd.Flags().String("issue-date", "", "Issue date (YYYY-MM-DD, required)")
	creditNotesCreatePartialCmd.Flags().String("description", "", "Line item description (required)")
	creditNotesCreatePartialCmd.Flags().String("account", "", "Account code (required)")
	creditNotesCreatePartialCmd.Flags().String("vat-code", "", "VAT code (required)")
	creditNotesCreatePartialCmd.Flags().String("amount", "", "Amount in decimal format e.g. '1000.00' (required)")
	creditNotesCreatePartialCmd.Flags().Int64("invoice-id", 0, "Associated invoice ID (optional)")
	creditNotesCreatePartialCmd.Flags().Int64("contact-id", 0, "Contact ID (optional)")
	creditNotesCreatePartialCmd.Flags().Int64("contact-person-id", 0, "Contact person ID (optional)")
	creditNotesCreatePartialCmd.Flags().String("credit-note-text", "", "Credit note text (optional)")
	creditNotesCreatePartialCmd.Flags().String("our-reference", "", "Our reference (optional)")
	creditNotesCreatePartialCmd.Flags().String("your-reference", "", "Your reference (optional)")
	creditNotesCreatePartialCmd.Flags().String("order-reference", "", "Order reference (optional)")
	creditNotesCreatePartialCmd.Flags().Int64("project-id", 0, "Project ID (optional)")
	creditNotesCreatePartialCmd.Flags().String("currency", "NOK", "Currency code (default NOK)")
	creditNotesCreatePartialCmd.Flags().Int64("quantity", 1, "Line item quantity (default 1)")
	creditNotesCreatePartialCmd.Flags().String("comment", "", "Line item comment (optional)")

	creditNotesSendCmd.Flags().Int64("credit-note-id", 0, "Credit note ID to send (required)")
	creditNotesSendCmd.Flags().String("method", "auto", "Delivery method (comma-separated: auto,email,ehf,efaktura,vipps,sms,letter)")
	creditNotesSendCmd.Flags().String("recipient-name", "", "Recipient name (optional)")
	creditNotesSendCmd.Flags().String("recipient-email", "", "Recipient email (optional)")
	creditNotesSendCmd.Flags().String("message", "", "Message to include (optional)")
	creditNotesSendCmd.Flags().Bool("include-document-attachments", false, "Include document attachments")
	creditNotesSendCmd.Flags().String("email-send-option", "", "Email send option (optional)")
	creditNotesSendCmd.Flags().String("organization-number", "", "Organization number (optional)")
	creditNotesSendCmd.Flags().String("mobile-number", "", "Mobile number (optional)")

	creditNotesCounterSetCmd.Flags().Int64("value", 0, "Counter value to set (required)")
	creditNotesCounterSetCmd.MarkFlagRequired("value")

	creditNotesCmd.AddCommand(creditNotesListCmd)
	creditNotesCmd.AddCommand(creditNotesGetCmd)
	creditNotesCmd.AddCommand(creditNotesCreateFullCmd)
	creditNotesCmd.AddCommand(creditNotesCreatePartialCmd)
	creditNotesCmd.AddCommand(creditNotesSendCmd)
	creditNotesCmd.AddCommand(creditNotesCounterCmd)
	creditNotesCounterCmd.AddCommand(creditNotesCounterSetCmd)
	rootCmd.AddCommand(creditNotesCmd)
}
