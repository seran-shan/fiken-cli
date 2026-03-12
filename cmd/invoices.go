package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var invoicesCmd = &cobra.Command{
	Use:   "invoices",
	Short: "Manage invoices",
}

var invoicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List invoices",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointInvoices, slug)
		invoices, err := FetchAllPages[api.Invoice](client, endpoint, nil, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching invoices: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(invoices)
		}

		if len(invoices) == 0 {
			output.PrintInfo("No invoices found.")
			return nil
		}

		table := output.NewTable("ID", "INVOICE NUMBER", "ISSUE DATE", "DUE DATE", "CUSTOMER", "NET", "VAT", "GROSS", "PAID")
		for _, inv := range invoices {
			table.AddRow(
				fmt.Sprintf("%d", inv.InvoiceId),
				fmt.Sprintf("%d", inv.InvoiceNumber),
				inv.IssueDate,
				inv.DueDate,
				inv.Customer.Name,
				output.FormatAmount(inv.Net),
				output.FormatAmount(inv.Vat),
				output.FormatAmount(inv.Gross),
				BoolToYesNo(inv.Paid),
			)
		}
		table.Print()
		fmt.Printf("\n%d invoices\n", len(invoices))
		return nil
	},
}

var invoicesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an invoice",
	Long:  "Create a new invoice with a single line item.",
	RunE: func(cmd *cobra.Command, args []string) error {
		issueDate, _ := cmd.Flags().GetString("issue-date")
		dueDate, _ := cmd.Flags().GetString("due-date")
		customerID, _ := cmd.Flags().GetInt64("customer-id")
		bankAccountCode, _ := cmd.Flags().GetString("bank-account-code")
		description, _ := cmd.Flags().GetString("description")
		unitPriceStr, _ := cmd.Flags().GetString("unit-price")
		quantity, _ := cmd.Flags().GetInt64("quantity")
		vatType, _ := cmd.Flags().GetString("vat-type")
		cash, _ := cmd.Flags().GetBool("cash")
		orderReference, _ := cmd.Flags().GetString("order-reference")
		ourReference, _ := cmd.Flags().GetString("our-reference")
		yourReference, _ := cmd.Flags().GetString("your-reference")
		productID, _ := cmd.Flags().GetInt64("product-id")

		var missing []string
		if issueDate == "" {
			missing = append(missing, "--issue-date")
		}
		if dueDate == "" {
			missing = append(missing, "--due-date")
		}
		if customerID == 0 {
			missing = append(missing, "--customer-id")
		}
		if bankAccountCode == "" {
			missing = append(missing, "--bank-account-code")
		}
		if description == "" {
			missing = append(missing, "--description")
		}
		if unitPriceStr == "" {
			missing = append(missing, "--unit-price")
		}
		if quantity == 0 {
			missing = append(missing, "--quantity")
		}
		if vatType == "" {
			missing = append(missing, "--vat-type")
		}
		if len(missing) > 0 {
			return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
		}

		unitPrice, err := ParseAmountToCents(unitPriceStr)
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

		line := api.InvoiceLineRequest{
			Description: description,
			Quantity:    quantity,
			UnitPrice:   unitPrice,
			VatType:     vatType,
		}
		if productID != 0 {
			line.ProductId = productID
		}

		req := api.InvoiceRequest{
			IssueDate:       issueDate,
			DueDate:         dueDate,
			CustomerId:      customerID,
			BankAccountCode: bankAccountCode,
			Cash:            cash,
			OrderReference:  orderReference,
			OurReference:    ourReference,
			YourReference:   yourReference,
			Lines:           []api.InvoiceLineRequest{line},
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointInvoices, slug), req)
		if err != nil {
			return fmt.Errorf("creating invoice: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing invoice ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Invoice created (ID: %d)", id))
		return nil
	},
}

var invoicesGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an invoice by ID",
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

		var invoice api.Invoice
		endpoint := fmt.Sprintf(api.EndpointInvoice, slug, id)
		if _, err := client.Get(endpoint, &invoice); err != nil {
			return fmt.Errorf("fetching invoice: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(invoice)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", invoice.InvoiceId))
		table.AddRow("Invoice Number", fmt.Sprintf("%d", invoice.InvoiceNumber))
		table.AddRow("Issue Date", invoice.IssueDate)
		table.AddRow("Due Date", invoice.DueDate)
		table.AddRow("Customer", invoice.Customer.Name)
		table.AddRow("Net", output.FormatAmount(invoice.Net))
		table.AddRow("VAT", output.FormatAmount(invoice.Vat))
		table.AddRow("Gross", output.FormatAmount(invoice.Gross))
		table.AddRow("Currency", invoice.Currency)
		table.AddRow("Paid", BoolToYesNo(invoice.Paid))
		table.AddRow("KID", invoice.Kid)
		table.Print()
		return nil
	},
}

var invoicesUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an invoice",
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

		req := api.UpdateInvoiceRequest{}

		if cmd.Flags().Changed("new-due-date") {
			req.NewDueDate, _ = cmd.Flags().GetString("new-due-date")
		}
		if cmd.Flags().Changed("sent-manually") {
			req.SentManually, _ = cmd.Flags().GetBool("sent-manually")
		}

		endpoint := fmt.Sprintf(api.EndpointInvoice, slug, id)
		if err := client.Patch(endpoint, req, nil); err != nil {
			return fmt.Errorf("updating invoice: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Invoice %d updated", id))
		return nil
	},
}

var invoicesAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a document to an invoice",
	RunE: func(cmd *cobra.Command, args []string) error {
		invoiceID, _ := cmd.Flags().GetInt64("id")
		filePath, _ := cmd.Flags().GetString("file")

		if err := ValidateFile(filePath); err != nil {
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

		fields := map[string]string{
			"filename": filepath.Base(filePath),
		}
		endpoint := fmt.Sprintf(api.EndpointInvoiceAttachments, slug, invoiceID)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("attaching to invoice: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to invoice %d", invoiceID))
		return nil
	},
}

func init() {
	invoicesCreateCmd.Flags().String("issue-date", "", "Issue date (YYYY-MM-DD, required)")
	invoicesCreateCmd.Flags().String("due-date", "", "Due date (YYYY-MM-DD, required)")
	invoicesCreateCmd.Flags().Int64("customer-id", 0, "Customer contact ID (required)")
	invoicesCreateCmd.Flags().String("bank-account-code", "", "Bank account code (required)")
	invoicesCreateCmd.Flags().String("description", "", "Line item description (required)")
	invoicesCreateCmd.Flags().String("unit-price", "", "Unit price in decimal format e.g. '1000.00' (required)")
	invoicesCreateCmd.Flags().Int64("quantity", 0, "Line item quantity (required)")
	invoicesCreateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM' (required)")
	invoicesCreateCmd.Flags().Bool("cash", false, "Whether this is a cash invoice")
	invoicesCreateCmd.Flags().String("order-reference", "", "Order reference (optional)")
	invoicesCreateCmd.Flags().String("our-reference", "", "Our reference (optional)")
	invoicesCreateCmd.Flags().String("your-reference", "", "Your reference (optional)")
	invoicesCreateCmd.Flags().Int64("product-id", 0, "Product ID for the line item (optional)")

	invoicesUpdateCmd.Flags().String("new-due-date", "", "New due date (YYYY-MM-DD)")
	invoicesUpdateCmd.Flags().Bool("sent-manually", false, "Mark invoice as sent manually")

	invoicesAttachCmd.Flags().Int64("id", 0, "Invoice ID to attach to (required)")
	invoicesAttachCmd.MarkFlagRequired("id")
	invoicesAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	invoicesAttachCmd.MarkFlagRequired("file")

	invoicesCmd.AddCommand(invoicesListCmd)
	invoicesCmd.AddCommand(invoicesCreateCmd)
	invoicesCmd.AddCommand(invoicesGetCmd)
	invoicesCmd.AddCommand(invoicesUpdateCmd)
	invoicesCmd.AddCommand(invoicesAttachCmd)
	rootCmd.AddCommand(invoicesCmd)
}
