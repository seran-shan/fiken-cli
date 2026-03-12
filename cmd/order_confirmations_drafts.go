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

var orderConfirmationsDraftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "Manage order confirmation drafts",
	Long:  "List, create, get, update, delete, attach, and finalize order confirmation drafts.",
}

var orderConfirmationsDraftsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List order confirmation drafts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointOrderConfirmationDrafts, slug)
		drafts, err := FetchAllPages[api.InvoiceishDraftResult](client, endpoint, nil, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching order confirmation drafts: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(drafts)
		}

		if len(drafts) == 0 {
			output.PrintInfo("No order confirmation drafts found.")
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
		fmt.Printf("\n%d order confirmation drafts\n", len(drafts))
		return nil
	},
}

var orderConfirmationsDraftsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an order confirmation draft",
	Long:  "Create a new order confirmation draft with a single line item.",
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
			Type:             "order_confirmation",
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

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointOrderConfirmationDrafts, slug), req)
		if err != nil {
			return fmt.Errorf("creating order confirmation draft: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing draft ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Order confirmation draft created (ID: %d)", id))
		return nil
	},
}

var orderConfirmationsDraftsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an order confirmation draft by ID",
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
		endpoint := fmt.Sprintf(api.EndpointOrderConfirmationDraft, slug, id)
		if _, err := client.Get(endpoint, &draft); err != nil {
			return fmt.Errorf("fetching order confirmation draft: %w", err)
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

var orderConfirmationsDraftsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an order confirmation draft",
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
		endpoint := fmt.Sprintf(api.EndpointOrderConfirmationDraft, slug, id)
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
			return fmt.Errorf("updating order confirmation draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Order confirmation draft %d updated", id))
		return nil
	},
}

var orderConfirmationsDraftsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an order confirmation draft",
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

		endpoint := fmt.Sprintf(api.EndpointOrderConfirmationDraft, slug, id)
		if _, err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting order confirmation draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Order confirmation draft %d deleted", id))
		return nil
	},
}

var orderConfirmationsDraftsAttachCmd = &cobra.Command{
	Use:   "attach <id>",
	Short: "Attach a document to an order confirmation draft",
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
		endpoint := fmt.Sprintf(api.EndpointOrderConfirmationDraftAttachments, slug, id)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("attaching to order confirmation draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to draft %d", id))
		return nil
	},
}

var orderConfirmationsDraftsAttachmentsCmd = &cobra.Command{
	Use:   "attachments [id]",
	Short: "List attachments for an order confirmation draft",
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
		var attachments []api.Attachment
		endpoint := fmt.Sprintf(api.EndpointOrderConfirmationDraftAttachments, slug, id)
		if _, err := client.Get(endpoint, &attachments); err != nil {
			return fmt.Errorf("fetching attachments: %w", err)
		}
		if jsonOutput {
			return output.PrintJSON(attachments)
		}
		if len(attachments) == 0 {
			output.PrintInfo("No attachments found.")
			return nil
		}
		table := output.NewTable("IDENTIFIER", "TYPE", "COMMENT", "DOWNLOAD URL")
		for _, a := range attachments {
			table.AddRow(
				a.Identifier,
				a.Type,
				a.Comment,
				a.DownloadUrl,
			)
		}
		table.Print()
		return nil
	},
}

var orderConfirmationsDraftsFinalizeCmd = &cobra.Command{
	Use:   "finalize <id>",
	Short: "Finalize an order confirmation draft and create an order confirmation",
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

		endpoint := fmt.Sprintf(api.EndpointOrderConfirmationDraftFinalize, slug, id)
		locationURL, err := client.PostEmpty(endpoint)
		if err != nil {
			return fmt.Errorf("finalizing order confirmation draft: %w", err)
		}

		confirmationID, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing confirmation ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Order confirmation draft %d finalized → Order confirmation %d created", id, confirmationID))
		return nil
	},
}

func init() {
	orderConfirmationsDraftsCreateCmd.Flags().Int64("customer-id", 0, "Customer contact ID")
	orderConfirmationsDraftsCreateCmd.Flags().Int64("contact-person-id", 0, "Contact person ID")
	orderConfirmationsDraftsCreateCmd.Flags().Int64("days", 0, "Days")
	orderConfirmationsDraftsCreateCmd.Flags().Float64("hours", 0, "Hours")
	orderConfirmationsDraftsCreateCmd.Flags().String("currency", "NOK", "Currency code")
	orderConfirmationsDraftsCreateCmd.Flags().String("bank-account-code", "", "Bank account code")
	orderConfirmationsDraftsCreateCmd.Flags().String("your-reference", "", "Your reference")
	orderConfirmationsDraftsCreateCmd.Flags().String("our-reference", "", "Our reference")
	orderConfirmationsDraftsCreateCmd.Flags().String("order-reference", "", "Order reference")
	orderConfirmationsDraftsCreateCmd.Flags().Int64("project-id", 0, "Project ID")
	orderConfirmationsDraftsCreateCmd.Flags().String("invoice-issue-date", "", "Invoice issue date (YYYY-MM-DD)")
	orderConfirmationsDraftsCreateCmd.Flags().String("invoice-number", "", "Invoice number")
	orderConfirmationsDraftsCreateCmd.Flags().String("description", "", "Line item description (required)")
	orderConfirmationsDraftsCreateCmd.Flags().Int64("quantity", 0, "Line item quantity (required)")
	orderConfirmationsDraftsCreateCmd.Flags().String("unit-price", "", "Unit price in decimal format e.g. '1000.00' (required)")
	orderConfirmationsDraftsCreateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM' (required)")
	orderConfirmationsDraftsCreateCmd.Flags().Int64("product-id", 0, "Product ID for the line item")
	orderConfirmationsDraftsCreateCmd.Flags().Int64("discount", 0, "Discount percentage")

	orderConfirmationsDraftsUpdateCmd.Flags().Int64("customer-id", 0, "Customer contact ID")
	orderConfirmationsDraftsUpdateCmd.Flags().Int64("contact-person-id", 0, "Contact person ID")
	orderConfirmationsDraftsUpdateCmd.Flags().Int64("days", 0, "Days")
	orderConfirmationsDraftsUpdateCmd.Flags().Float64("hours", 0, "Hours")
	orderConfirmationsDraftsUpdateCmd.Flags().String("currency", "", "Currency code")
	orderConfirmationsDraftsUpdateCmd.Flags().String("bank-account-code", "", "Bank account code")
	orderConfirmationsDraftsUpdateCmd.Flags().String("your-reference", "", "Your reference")
	orderConfirmationsDraftsUpdateCmd.Flags().String("our-reference", "", "Our reference")
	orderConfirmationsDraftsUpdateCmd.Flags().String("order-reference", "", "Order reference")
	orderConfirmationsDraftsUpdateCmd.Flags().Int64("project-id", 0, "Project ID")
	orderConfirmationsDraftsUpdateCmd.Flags().String("invoice-issue-date", "", "Invoice issue date (YYYY-MM-DD)")
	orderConfirmationsDraftsUpdateCmd.Flags().String("invoice-number", "", "Invoice number")
	orderConfirmationsDraftsUpdateCmd.Flags().String("description", "", "Line item description")
	orderConfirmationsDraftsUpdateCmd.Flags().Int64("quantity", 0, "Line item quantity")
	orderConfirmationsDraftsUpdateCmd.Flags().String("unit-price", "", "Unit price in decimal format e.g. '1000.00'")
	orderConfirmationsDraftsUpdateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM'")
	orderConfirmationsDraftsUpdateCmd.Flags().Int64("product-id", 0, "Product ID for the line item")
	orderConfirmationsDraftsUpdateCmd.Flags().Int64("discount", 0, "Discount percentage")

	orderConfirmationsDraftsAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	orderConfirmationsDraftsAttachCmd.MarkFlagRequired("file")

	orderConfirmationsDraftsCmd.AddCommand(orderConfirmationsDraftsListCmd)
	orderConfirmationsDraftsCmd.AddCommand(orderConfirmationsDraftsCreateCmd)
	orderConfirmationsDraftsCmd.AddCommand(orderConfirmationsDraftsGetCmd)
	orderConfirmationsDraftsCmd.AddCommand(orderConfirmationsDraftsUpdateCmd)
	orderConfirmationsDraftsCmd.AddCommand(orderConfirmationsDraftsDeleteCmd)
	orderConfirmationsDraftsCmd.AddCommand(orderConfirmationsDraftsAttachCmd)
	orderConfirmationsDraftsCmd.AddCommand(orderConfirmationsDraftsFinalizeCmd)
	orderConfirmationsDraftsCmd.AddCommand(orderConfirmationsDraftsAttachmentsCmd)

	orderConfirmationsCmd.AddCommand(orderConfirmationsDraftsCmd)
}
