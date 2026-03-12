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

var invoicesDraftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "Manage invoice drafts",
	Long:  "List, create, get, update, delete, attach, and finalize invoice drafts.",
}

var invoicesDraftsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List invoice drafts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointInvoiceDrafts, slug)
		drafts, err := FetchAllPages[api.InvoiceDraft](client, endpoint, nil, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching invoice drafts: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(drafts)
		}

		if len(drafts) == 0 {
			output.PrintInfo("No invoice drafts found.")
			return nil
		}

		table := output.NewTable("DRAFT ID", "TYPE", "CUSTOMER ID", "NET", "VAT", "GROSS", "LAST MODIFIED")
		for _, d := range drafts {
			table.AddRow(
				fmt.Sprintf("%d", d.DraftId),
				d.Type,
				fmt.Sprintf("%d", d.CustomerId),
				output.FormatAmount(d.Net),
				output.FormatAmount(d.Vat),
				output.FormatAmount(d.Gross),
				d.LastModifiedDate,
			)
		}
		table.Print()
		fmt.Printf("\n%d invoice drafts\n", len(drafts))
		return nil
	},
}

var invoicesDraftsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an invoice draft",
	Long:  "Create a new invoice draft with a single line item.",
	RunE: func(cmd *cobra.Command, args []string) error {
		draftType, _ := cmd.Flags().GetString("type")
		customerID, _ := cmd.Flags().GetInt64("customer-id")
		daysUntilDue, _ := cmd.Flags().GetInt64("days-until-due")
		description, _ := cmd.Flags().GetString("description")
		quantity, _ := cmd.Flags().GetInt64("quantity")
		unitPriceStr, _ := cmd.Flags().GetString("unit-price")
		vatType, _ := cmd.Flags().GetString("vat-type")
		bankAccountNumber, _ := cmd.Flags().GetString("bank-account-number")
		ourReference, _ := cmd.Flags().GetString("our-reference")
		yourReference, _ := cmd.Flags().GetString("your-reference")
		orderReference, _ := cmd.Flags().GetString("order-reference")
		productID, _ := cmd.Flags().GetInt64("product-id")

		var missing []string
		if draftType == "" {
			missing = append(missing, "--type")
		}
		if customerID == 0 {
			missing = append(missing, "--customer-id")
		}
		if daysUntilDue == 0 {
			missing = append(missing, "--days-until-due")
		}
		if description == "" {
			missing = append(missing, "--description")
		}
		if quantity == 0 {
			missing = append(missing, "--quantity")
		}
		if unitPriceStr == "" {
			missing = append(missing, "--unit-price")
		}
		if vatType == "" {
			missing = append(missing, "--vat-type")
		}
		if len(missing) > 0 {
			return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
		}

		if draftType != "invoice" && draftType != "cash_invoice" {
			return fmt.Errorf("invalid --type %q: must be \"invoice\" or \"cash_invoice\"", draftType)
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

		line := api.InvoiceDraftLineRequest{
			Description: description,
			Quantity:    quantity,
			UnitPrice:   unitPrice,
			VatType:     vatType,
		}
		if productID != 0 {
			line.ProductId = productID
		}

		req := api.InvoiceDraftRequest{
			Type:              draftType,
			CustomerId:        customerID,
			DaysUntilDueDate:  daysUntilDue,
			BankAccountNumber: bankAccountNumber,
			OurReference:      ourReference,
			YourReference:     yourReference,
			OrderReference:    orderReference,
			Lines:             []api.InvoiceDraftLineRequest{line},
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointInvoiceDrafts, slug), req)
		if err != nil {
			return fmt.Errorf("creating invoice draft: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing draft ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Invoice draft created (ID: %d)", id))
		return nil
	},
}

var invoicesDraftsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an invoice draft by ID",
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

		var draft api.InvoiceDraft
		endpoint := fmt.Sprintf(api.EndpointInvoiceDraft, slug, id)
		if _, err := client.Get(endpoint, &draft); err != nil {
			return fmt.Errorf("fetching invoice draft: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(draft)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("Draft ID", fmt.Sprintf("%d", draft.DraftId))
		table.AddRow("UUID", draft.Uuid)
		table.AddRow("Type", draft.Type)
		table.AddRow("Customer ID", fmt.Sprintf("%d", draft.CustomerId))
		table.AddRow("Net", output.FormatAmount(draft.Net))
		table.AddRow("VAT", output.FormatAmount(draft.Vat))
		table.AddRow("Gross", output.FormatAmount(draft.Gross))
		table.AddRow("Last Modified", draft.LastModifiedDate)
		table.Print()
		return nil
	},
}

var invoicesDraftsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an invoice draft",
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

		var existing api.InvoiceDraft
		endpoint := fmt.Sprintf(api.EndpointInvoiceDraft, slug, id)
		if _, err := client.Get(endpoint, &existing); err != nil {
			return fmt.Errorf("fetching draft for update: %w", err)
		}

		// Start from existing values
		req := api.InvoiceDraftRequest{
			Type:       existing.Type,
			CustomerId: existing.CustomerId,
			Lines:      existing.Lines,
		}

		if cmd.Flags().Changed("type") {
			req.Type, _ = cmd.Flags().GetString("type")
			if req.Type != "invoice" && req.Type != "cash_invoice" {
				return fmt.Errorf("invalid --type %q: must be \"invoice\" or \"cash_invoice\"", req.Type)
			}
		}
		if cmd.Flags().Changed("customer-id") {
			req.CustomerId, _ = cmd.Flags().GetInt64("customer-id")
		}
		if cmd.Flags().Changed("days-until-due") {
			req.DaysUntilDueDate, _ = cmd.Flags().GetInt64("days-until-due")
		}
		if cmd.Flags().Changed("bank-account-number") {
			req.BankAccountNumber, _ = cmd.Flags().GetString("bank-account-number")
		}
		if cmd.Flags().Changed("our-reference") {
			req.OurReference, _ = cmd.Flags().GetString("our-reference")
		}
		if cmd.Flags().Changed("your-reference") {
			req.YourReference, _ = cmd.Flags().GetString("your-reference")
		}
		if cmd.Flags().Changed("order-reference") {
			req.OrderReference, _ = cmd.Flags().GetString("order-reference")
		}

		// If any line-level flag changed, rebuild the single line
		lineChanged := cmd.Flags().Changed("description") || cmd.Flags().Changed("quantity") ||
			cmd.Flags().Changed("unit-price") || cmd.Flags().Changed("vat-type") ||
			cmd.Flags().Changed("product-id")

		if lineChanged {
			var line api.InvoiceDraftLineRequest
			if len(existing.Lines) > 0 {
				line = existing.Lines[0]
			}
			if cmd.Flags().Changed("description") {
				line.Description, _ = cmd.Flags().GetString("description")
			}
			if cmd.Flags().Changed("quantity") {
				line.Quantity, _ = cmd.Flags().GetInt64("quantity")
			}
			if cmd.Flags().Changed("unit-price") {
				unitPriceStr, _ := cmd.Flags().GetString("unit-price")
				unitPrice, err := ParseAmountToCents(unitPriceStr)
				if err != nil {
					return err
				}
				line.UnitPrice = unitPrice
			}
			if cmd.Flags().Changed("vat-type") {
				line.VatType, _ = cmd.Flags().GetString("vat-type")
			}
			if cmd.Flags().Changed("product-id") {
				line.ProductId, _ = cmd.Flags().GetInt64("product-id")
			}
			req.Lines = []api.InvoiceDraftLineRequest{line}
		}

		_, err = client.Put(endpoint, req)
		if err != nil {
			return fmt.Errorf("updating invoice draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Invoice draft %d updated", id))
		return nil
	},
}

var invoicesDraftsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an invoice draft",
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

		endpoint := fmt.Sprintf(api.EndpointInvoiceDraft, slug, id)
		if err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting invoice draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Invoice draft %d deleted", id))
		return nil
	},
}

var invoicesDraftsAttachCmd = &cobra.Command{
	Use:   "attach <id>",
	Short: "Attach a document to an invoice draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID %q: %w", args[0], err)
		}

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
		endpoint := fmt.Sprintf(api.EndpointInvoiceDraftAttachments, slug, id)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("attaching to invoice draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to draft %d", id))
		return nil
	},
}

var invoicesDraftsFinalizeCmd = &cobra.Command{
	Use:   "finalize <id>",
	Short: "Finalize an invoice draft and create an invoice",
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

		endpoint := fmt.Sprintf(api.EndpointInvoiceDraftFinalize, slug, id)
		locationURL, err := client.PostEmpty(endpoint)
		if err != nil {
			return fmt.Errorf("finalizing invoice draft: %w", err)
		}

		invoiceID, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing invoice ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Invoice draft %d finalized → Invoice %d created", id, invoiceID))
		return nil
	},
}

func init() {
	// Create flags
	invoicesDraftsCreateCmd.Flags().String("type", "", "Draft type: 'invoice' or 'cash_invoice' (required)")
	invoicesDraftsCreateCmd.Flags().Int64("customer-id", 0, "Customer contact ID (required)")
	invoicesDraftsCreateCmd.Flags().Int64("days-until-due", 0, "Days until due date (required)")
	invoicesDraftsCreateCmd.Flags().String("description", "", "Line item description (required)")
	invoicesDraftsCreateCmd.Flags().Int64("quantity", 0, "Line item quantity (required)")
	invoicesDraftsCreateCmd.Flags().String("unit-price", "", "Unit price in decimal format e.g. '1000.00' (required)")
	invoicesDraftsCreateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM' (required)")
	invoicesDraftsCreateCmd.Flags().String("bank-account-number", "", "Bank account number (optional)")
	invoicesDraftsCreateCmd.Flags().String("our-reference", "", "Our reference (optional)")
	invoicesDraftsCreateCmd.Flags().String("your-reference", "", "Your reference (optional)")
	invoicesDraftsCreateCmd.Flags().String("order-reference", "", "Order reference (optional)")
	invoicesDraftsCreateCmd.Flags().Int64("product-id", 0, "Product ID for the line item (optional)")

	// Update flags (same as create, all optional)
	invoicesDraftsUpdateCmd.Flags().String("type", "", "Draft type: 'invoice' or 'cash_invoice'")
	invoicesDraftsUpdateCmd.Flags().Int64("customer-id", 0, "Customer contact ID")
	invoicesDraftsUpdateCmd.Flags().Int64("days-until-due", 0, "Days until due date")
	invoicesDraftsUpdateCmd.Flags().String("description", "", "Line item description")
	invoicesDraftsUpdateCmd.Flags().Int64("quantity", 0, "Line item quantity")
	invoicesDraftsUpdateCmd.Flags().String("unit-price", "", "Unit price in decimal format e.g. '1000.00'")
	invoicesDraftsUpdateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM'")
	invoicesDraftsUpdateCmd.Flags().String("bank-account-number", "", "Bank account number")
	invoicesDraftsUpdateCmd.Flags().String("our-reference", "", "Our reference")
	invoicesDraftsUpdateCmd.Flags().String("your-reference", "", "Your reference")
	invoicesDraftsUpdateCmd.Flags().String("order-reference", "", "Order reference")
	invoicesDraftsUpdateCmd.Flags().Int64("product-id", 0, "Product ID for the line item")

	// Attach flags
	invoicesDraftsAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	invoicesDraftsAttachCmd.MarkFlagRequired("file")

	// Register subcommands on drafts parent
	invoicesDraftsCmd.AddCommand(invoicesDraftsListCmd)
	invoicesDraftsCmd.AddCommand(invoicesDraftsCreateCmd)
	invoicesDraftsCmd.AddCommand(invoicesDraftsGetCmd)
	invoicesDraftsCmd.AddCommand(invoicesDraftsUpdateCmd)
	invoicesDraftsCmd.AddCommand(invoicesDraftsDeleteCmd)
	invoicesDraftsCmd.AddCommand(invoicesDraftsAttachCmd)
	invoicesDraftsCmd.AddCommand(invoicesDraftsFinalizeCmd)

	// Register drafts on invoices
	invoicesCmd.AddCommand(invoicesDraftsCmd)
}
