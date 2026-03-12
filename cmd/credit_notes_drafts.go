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

var creditNotesDraftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "Manage credit note drafts",
}

var creditNotesDraftsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List credit note drafts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointCreditNoteDrafts, slug)
		drafts, err := FetchAllPages[api.InvoiceishDraftResult](client, endpoint, nil, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching credit note drafts: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(drafts)
		}

		if len(drafts) == 0 {
			output.PrintInfo("No credit note drafts found.")
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
		fmt.Printf("\n%d credit note drafts\n", len(drafts))
		return nil
	},
}

var creditNotesDraftsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a credit note draft",
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID, _ := cmd.Flags().GetInt64("customer-id")
		contactPersonID, _ := cmd.Flags().GetInt64("contact-person-id")
		days, _ := cmd.Flags().GetInt64("days")
		hours, _ := cmd.Flags().GetFloat64("hours")
		currency, _ := cmd.Flags().GetString("currency")
		bankAccountCode, _ := cmd.Flags().GetString("bank-account-code")
		yourReference, _ := cmd.Flags().GetString("your-reference")
		ourReference, _ := cmd.Flags().GetString("our-reference")
		orderReference, _ := cmd.Flags().GetString("order-reference")
		projectID, _ := cmd.Flags().GetInt64("project-id")
		invoiceIssueDate, _ := cmd.Flags().GetString("invoice-issue-date")
		invoiceNumber, _ := cmd.Flags().GetString("invoice-number")
		description, _ := cmd.Flags().GetString("description")
		quantity, _ := cmd.Flags().GetInt64("quantity")
		unitPriceStr, _ := cmd.Flags().GetString("unit-price")
		vatType, _ := cmd.Flags().GetString("vat-type")
		productID, _ := cmd.Flags().GetInt64("product-id")
		discount, _ := cmd.Flags().GetInt64("discount")

		var missing []string
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
		if discount != 0 {
			line.Discount = discount
		}

		req := api.InvoiceishDraftRequest{
			Type:             "creditNote",
			CustomerId:       customerID,
			ContactPersonId:  contactPersonID,
			Days:             days,
			Hours:            hours,
			Currency:         currency,
			BankAccountCode:  bankAccountCode,
			YourReference:    yourReference,
			OurReference:     ourReference,
			OrderReference:   orderReference,
			ProjectId:        projectID,
			InvoiceIssueDate: invoiceIssueDate,
			InvoiceNumber:    invoiceNumber,
			Lines:            []api.InvoiceDraftLineRequest{line},
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointCreditNoteDrafts, slug), req)
		if err != nil {
			return fmt.Errorf("creating credit note draft: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing draft ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Credit note draft created (ID: %d)", id))
		return nil
	},
}

var creditNotesDraftsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a credit note draft by ID",
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

		var draft api.InvoiceishDraftResult
		endpoint := fmt.Sprintf(api.EndpointCreditNoteDraft, slug, id)
		if _, err := client.Get(endpoint, &draft); err != nil {
			return fmt.Errorf("fetching credit note draft: %w", err)
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
		table.AddRow("Currency", draft.Currency)
		table.AddRow("Last Modified", draft.LastModifiedDate)
		table.Print()

		if len(draft.Lines) > 0 {
			fmt.Println()
			lines := output.NewTable("DESCRIPTION", "QUANTITY", "UNIT PRICE", "VAT TYPE")
			for _, l := range draft.Lines {
				lines.AddRow(
					l.Description,
					fmt.Sprintf("%d", l.Quantity),
					output.FormatAmount(l.UnitPrice),
					l.VatType,
				)
			}
			lines.Print()
		}
		return nil
	},
}

var creditNotesDraftsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a credit note draft",
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

		var existing api.InvoiceishDraftResult
		endpoint := fmt.Sprintf(api.EndpointCreditNoteDraft, slug, id)
		if _, err := client.Get(endpoint, &existing); err != nil {
			return fmt.Errorf("fetching draft for update: %w", err)
		}

		req := api.InvoiceishDraftRequest{
			Type:             existing.Type,
			CustomerId:       existing.CustomerId,
			ContactPersonId:  existing.ContactPersonId,
			Days:             existing.Days,
			Hours:            existing.Hours,
			Currency:         existing.Currency,
			BankAccountCode:  existing.BankAccountCode,
			YourReference:    existing.YourReference,
			OurReference:     existing.OurReference,
			OrderReference:   existing.OrderReference,
			ProjectId:        existing.ProjectId,
			InvoiceIssueDate: existing.InvoiceIssueDate,
			InvoiceNumber:    existing.InvoiceNumber,
			Lines:            existing.Lines,
		}

		if cmd.Flags().Changed("customer-id") {
			req.CustomerId, _ = cmd.Flags().GetInt64("customer-id")
		}
		if cmd.Flags().Changed("contact-person-id") {
			req.ContactPersonId, _ = cmd.Flags().GetInt64("contact-person-id")
		}
		if cmd.Flags().Changed("days") {
			req.Days, _ = cmd.Flags().GetInt64("days")
		}
		if cmd.Flags().Changed("hours") {
			req.Hours, _ = cmd.Flags().GetFloat64("hours")
		}
		if cmd.Flags().Changed("currency") {
			req.Currency, _ = cmd.Flags().GetString("currency")
		}
		if cmd.Flags().Changed("bank-account-code") {
			req.BankAccountCode, _ = cmd.Flags().GetString("bank-account-code")
		}
		if cmd.Flags().Changed("your-reference") {
			req.YourReference, _ = cmd.Flags().GetString("your-reference")
		}
		if cmd.Flags().Changed("our-reference") {
			req.OurReference, _ = cmd.Flags().GetString("our-reference")
		}
		if cmd.Flags().Changed("order-reference") {
			req.OrderReference, _ = cmd.Flags().GetString("order-reference")
		}
		if cmd.Flags().Changed("project-id") {
			req.ProjectId, _ = cmd.Flags().GetInt64("project-id")
		}
		if cmd.Flags().Changed("invoice-issue-date") {
			req.InvoiceIssueDate, _ = cmd.Flags().GetString("invoice-issue-date")
		}
		if cmd.Flags().Changed("invoice-number") {
			req.InvoiceNumber, _ = cmd.Flags().GetString("invoice-number")
		}

		lineChanged := cmd.Flags().Changed("description") || cmd.Flags().Changed("quantity") ||
			cmd.Flags().Changed("unit-price") || cmd.Flags().Changed("vat-type") ||
			cmd.Flags().Changed("product-id") || cmd.Flags().Changed("discount")

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
			if cmd.Flags().Changed("discount") {
				line.Discount, _ = cmd.Flags().GetInt64("discount")
			}
			req.Lines = []api.InvoiceDraftLineRequest{line}
		}

		_, err = client.Put(endpoint, req)
		if err != nil {
			return fmt.Errorf("updating credit note draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Credit note draft %d updated", id))
		return nil
	},
}

var creditNotesDraftsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a credit note draft",
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

		endpoint := fmt.Sprintf(api.EndpointCreditNoteDraft, slug, id)
		if _, err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting credit note draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Credit note draft %d deleted", id))
		return nil
	},
}

var creditNotesDraftsAttachCmd = &cobra.Command{
	Use:   "attach <id>",
	Short: "Attach a document to a credit note draft",
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
		endpoint := fmt.Sprintf(api.EndpointCreditNoteDraftAttachments, slug, id)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("attaching to credit note draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to draft %d", id))
		return nil
	},
}

var creditNotesDraftsFinalizeCmd = &cobra.Command{
	Use:   "finalize <id>",
	Short: "Finalize a credit note draft and create a credit note",
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

		endpoint := fmt.Sprintf(api.EndpointCreditNoteDraftFinalize, slug, id)
		locationURL, err := client.PostEmpty(endpoint)
		if err != nil {
			return fmt.Errorf("finalizing credit note draft: %w", err)
		}

		creditNoteID, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing credit note ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Credit note draft %d finalized → Credit note %d created", id, creditNoteID))
		return nil
	},
}

func init() {
	creditNotesDraftsCreateCmd.Flags().Int64("customer-id", 0, "Customer contact ID")
	creditNotesDraftsCreateCmd.Flags().Int64("contact-person-id", 0, "Contact person ID")
	creditNotesDraftsCreateCmd.Flags().Int64("days", 0, "Days")
	creditNotesDraftsCreateCmd.Flags().Float64("hours", 0, "Hours")
	creditNotesDraftsCreateCmd.Flags().String("currency", "NOK", "Currency code")
	creditNotesDraftsCreateCmd.Flags().String("bank-account-code", "", "Bank account code")
	creditNotesDraftsCreateCmd.Flags().String("your-reference", "", "Your reference")
	creditNotesDraftsCreateCmd.Flags().String("our-reference", "", "Our reference")
	creditNotesDraftsCreateCmd.Flags().String("order-reference", "", "Order reference")
	creditNotesDraftsCreateCmd.Flags().Int64("project-id", 0, "Project ID")
	creditNotesDraftsCreateCmd.Flags().String("invoice-issue-date", "", "Invoice issue date (YYYY-MM-DD)")
	creditNotesDraftsCreateCmd.Flags().String("invoice-number", "", "Invoice number")
	creditNotesDraftsCreateCmd.Flags().String("description", "", "Line item description (required)")
	creditNotesDraftsCreateCmd.Flags().Int64("quantity", 0, "Line item quantity (required)")
	creditNotesDraftsCreateCmd.Flags().String("unit-price", "", "Unit price in decimal format e.g. '1000.00' (required)")
	creditNotesDraftsCreateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM' (required)")
	creditNotesDraftsCreateCmd.Flags().Int64("product-id", 0, "Product ID for the line item")
	creditNotesDraftsCreateCmd.Flags().Int64("discount", 0, "Discount percentage")

	creditNotesDraftsUpdateCmd.Flags().Int64("customer-id", 0, "Customer contact ID")
	creditNotesDraftsUpdateCmd.Flags().Int64("contact-person-id", 0, "Contact person ID")
	creditNotesDraftsUpdateCmd.Flags().Int64("days", 0, "Days")
	creditNotesDraftsUpdateCmd.Flags().Float64("hours", 0, "Hours")
	creditNotesDraftsUpdateCmd.Flags().String("currency", "", "Currency code")
	creditNotesDraftsUpdateCmd.Flags().String("bank-account-code", "", "Bank account code")
	creditNotesDraftsUpdateCmd.Flags().String("your-reference", "", "Your reference")
	creditNotesDraftsUpdateCmd.Flags().String("our-reference", "", "Our reference")
	creditNotesDraftsUpdateCmd.Flags().String("order-reference", "", "Order reference")
	creditNotesDraftsUpdateCmd.Flags().Int64("project-id", 0, "Project ID")
	creditNotesDraftsUpdateCmd.Flags().String("invoice-issue-date", "", "Invoice issue date (YYYY-MM-DD)")
	creditNotesDraftsUpdateCmd.Flags().String("invoice-number", "", "Invoice number")
	creditNotesDraftsUpdateCmd.Flags().String("description", "", "Line item description")
	creditNotesDraftsUpdateCmd.Flags().Int64("quantity", 0, "Line item quantity")
	creditNotesDraftsUpdateCmd.Flags().String("unit-price", "", "Unit price in decimal format e.g. '1000.00'")
	creditNotesDraftsUpdateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM'")
	creditNotesDraftsUpdateCmd.Flags().Int64("product-id", 0, "Product ID for the line item")
	creditNotesDraftsUpdateCmd.Flags().Int64("discount", 0, "Discount percentage")

	creditNotesDraftsAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	creditNotesDraftsAttachCmd.MarkFlagRequired("file")

	creditNotesDraftsCmd.AddCommand(creditNotesDraftsListCmd)
	creditNotesDraftsCmd.AddCommand(creditNotesDraftsCreateCmd)
	creditNotesDraftsCmd.AddCommand(creditNotesDraftsGetCmd)
	creditNotesDraftsCmd.AddCommand(creditNotesDraftsUpdateCmd)
	creditNotesDraftsCmd.AddCommand(creditNotesDraftsDeleteCmd)
	creditNotesDraftsCmd.AddCommand(creditNotesDraftsAttachCmd)
	creditNotesDraftsCmd.AddCommand(creditNotesDraftsFinalizeCmd)

	creditNotesCmd.AddCommand(creditNotesDraftsCmd)
}
