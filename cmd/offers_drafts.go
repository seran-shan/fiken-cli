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

var offersDraftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "Manage offer drafts",
}

var offersDraftsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List offer drafts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointOfferDrafts, slug)
		drafts, err := FetchAllPages[api.InvoiceishDraftResult](client, endpoint, nil, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching offer drafts: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(drafts)
		}

		if len(drafts) == 0 {
			output.PrintInfo("No offer drafts found.")
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
		fmt.Printf("\n%d offer drafts\n", len(drafts))
		return nil
	},
}

var offersDraftsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an offer draft",
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
			Type:             "offer",
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

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointOfferDrafts, slug), req)
		if err != nil {
			return fmt.Errorf("creating offer draft: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing draft ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Offer draft created (ID: %d)", id))
		return nil
	},
}

var offersDraftsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an offer draft by ID",
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
		endpoint := fmt.Sprintf(api.EndpointOfferDraft, slug, id)
		if _, err := client.Get(endpoint, &draft); err != nil {
			return fmt.Errorf("fetching offer draft: %w", err)
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

var offersDraftsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an offer draft",
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
		endpoint := fmt.Sprintf(api.EndpointOfferDraft, slug, id)
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
			return fmt.Errorf("updating offer draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Offer draft %d updated", id))
		return nil
	},
}

var offersDraftsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an offer draft",
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

		endpoint := fmt.Sprintf(api.EndpointOfferDraft, slug, id)
		if _, err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting offer draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Offer draft %d deleted", id))
		return nil
	},
}

var offersDraftsAttachCmd = &cobra.Command{
	Use:   "attach <id>",
	Short: "Attach a document to an offer draft",
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
		endpoint := fmt.Sprintf(api.EndpointOfferDraftAttachments, slug, id)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("attaching to offer draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to draft %d", id))
		return nil
	},
}

var offersDraftsFinalizeCmd = &cobra.Command{
	Use:   "finalize <id>",
	Short: "Finalize an offer draft and create an offer",
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

		endpoint := fmt.Sprintf(api.EndpointOfferDraftFinalize, slug, id)
		locationURL, err := client.PostEmpty(endpoint)
		if err != nil {
			return fmt.Errorf("finalizing offer draft: %w", err)
		}

		offerID, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing offer ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Offer draft %d finalized → Offer %d created", id, offerID))
		return nil
	},
}

func init() {
	// Create flags
	offersDraftsCreateCmd.Flags().Int64("customer-id", 0, "Customer contact ID")
	offersDraftsCreateCmd.Flags().Int64("contact-person-id", 0, "Contact person ID")
	offersDraftsCreateCmd.Flags().Int64("days", 0, "Days")
	offersDraftsCreateCmd.Flags().Float64("hours", 0, "Hours")
	offersDraftsCreateCmd.Flags().String("currency", "NOK", "Currency code")
	offersDraftsCreateCmd.Flags().String("bank-account-code", "", "Bank account code")
	offersDraftsCreateCmd.Flags().String("your-reference", "", "Your reference")
	offersDraftsCreateCmd.Flags().String("our-reference", "", "Our reference")
	offersDraftsCreateCmd.Flags().String("order-reference", "", "Order reference")
	offersDraftsCreateCmd.Flags().Int64("project-id", 0, "Project ID")
	offersDraftsCreateCmd.Flags().String("invoice-issue-date", "", "Invoice issue date (YYYY-MM-DD)")
	offersDraftsCreateCmd.Flags().String("invoice-number", "", "Invoice number")
	offersDraftsCreateCmd.Flags().String("description", "", "Line item description (required)")
	offersDraftsCreateCmd.Flags().Int64("quantity", 0, "Line item quantity (required)")
	offersDraftsCreateCmd.Flags().String("unit-price", "", "Unit price in decimal format e.g. '1000.00' (required)")
	offersDraftsCreateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM' (required)")
	offersDraftsCreateCmd.Flags().Int64("product-id", 0, "Product ID for the line item")
	offersDraftsCreateCmd.Flags().Int64("discount", 0, "Discount percentage")

	// Update flags (same as create, all optional)
	offersDraftsUpdateCmd.Flags().Int64("customer-id", 0, "Customer contact ID")
	offersDraftsUpdateCmd.Flags().Int64("contact-person-id", 0, "Contact person ID")
	offersDraftsUpdateCmd.Flags().Int64("days", 0, "Days")
	offersDraftsUpdateCmd.Flags().Float64("hours", 0, "Hours")
	offersDraftsUpdateCmd.Flags().String("currency", "", "Currency code")
	offersDraftsUpdateCmd.Flags().String("bank-account-code", "", "Bank account code")
	offersDraftsUpdateCmd.Flags().String("your-reference", "", "Your reference")
	offersDraftsUpdateCmd.Flags().String("our-reference", "", "Our reference")
	offersDraftsUpdateCmd.Flags().String("order-reference", "", "Order reference")
	offersDraftsUpdateCmd.Flags().Int64("project-id", 0, "Project ID")
	offersDraftsUpdateCmd.Flags().String("invoice-issue-date", "", "Invoice issue date (YYYY-MM-DD)")
	offersDraftsUpdateCmd.Flags().String("invoice-number", "", "Invoice number")
	offersDraftsUpdateCmd.Flags().String("description", "", "Line item description")
	offersDraftsUpdateCmd.Flags().Int64("quantity", 0, "Line item quantity")
	offersDraftsUpdateCmd.Flags().String("unit-price", "", "Unit price in decimal format e.g. '1000.00'")
	offersDraftsUpdateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM'")
	offersDraftsUpdateCmd.Flags().Int64("product-id", 0, "Product ID for the line item")
	offersDraftsUpdateCmd.Flags().Int64("discount", 0, "Discount percentage")

	// Attach flags
	offersDraftsAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	offersDraftsAttachCmd.MarkFlagRequired("file")

	// Register subcommands
	offersDraftsCmd.AddCommand(offersDraftsListCmd)
	offersDraftsCmd.AddCommand(offersDraftsCreateCmd)
	offersDraftsCmd.AddCommand(offersDraftsGetCmd)
	offersDraftsCmd.AddCommand(offersDraftsUpdateCmd)
	offersDraftsCmd.AddCommand(offersDraftsDeleteCmd)
	offersDraftsCmd.AddCommand(offersDraftsAttachCmd)
	offersDraftsCmd.AddCommand(offersDraftsFinalizeCmd)

	offersCmd.AddCommand(offersDraftsCmd)
}
