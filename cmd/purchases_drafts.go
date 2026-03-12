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

var purchasesDraftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "Manage purchase drafts",
}

var purchasesDraftsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List purchase drafts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointPurchaseDrafts, slug)
		drafts, err := FetchAllPages[api.PurchaseDraft](client, endpoint, nil, 25, 4)
		if err != nil {
			return fmt.Errorf("fetching purchase drafts: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(drafts)
		}

		if len(drafts) == 0 {
			output.PrintInfo("No purchase drafts found.")
			return nil
		}

		table := output.NewTable("DRAFT ID", "CASH", "SUPPLIER ID", "DUE DATE", "INVOICE NUMBER", "PAID")
		for _, d := range drafts {
			table.AddRow(
				fmt.Sprintf("%d", d.DraftId),
				BoolToYesNo(d.Cash),
				fmt.Sprintf("%d", d.SupplierId),
				d.DueDate,
				d.InvoiceNumber,
				BoolToYesNo(d.Paid),
			)
		}
		table.Print()

		fmt.Printf("\n%d purchase drafts\n", len(drafts))
		return nil
	},
}

var purchasesDraftsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a purchase draft",
	Long:  "Create a new purchase draft with a single draft line.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cash, _ := cmd.Flags().GetBool("cash")
		description, _ := cmd.Flags().GetString("description")
		account, _ := cmd.Flags().GetString("account")
		amountStr, _ := cmd.Flags().GetString("amount")
		vatType, _ := cmd.Flags().GetString("vat-type")
		supplierID, _ := cmd.Flags().GetInt64("supplier-id")
		dueDate, _ := cmd.Flags().GetString("due-date")
		invoiceNumber, _ := cmd.Flags().GetString("invoice-number")
		kid, _ := cmd.Flags().GetString("kid")
		projectID, _ := cmd.Flags().GetInt64("project-id")
		invoiceIssueDate, _ := cmd.Flags().GetString("invoice-issue-date")
		paid, _ := cmd.Flags().GetBool("paid")
		currency, _ := cmd.Flags().GetString("currency")

		var missing []string
		if !cmd.Flags().Changed("cash") {
			missing = append(missing, "--cash")
		}
		if description == "" {
			missing = append(missing, "--description")
		}
		if account == "" {
			missing = append(missing, "--account")
		}
		if amountStr == "" {
			missing = append(missing, "--amount")
		}
		if vatType == "" {
			missing = append(missing, "--vat-type")
		}
		if len(missing) > 0 {
			return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
		}

		amountCents, err := ParseAmountToCents(amountStr)
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

		req := api.PurchaseDraftRequest{
			Cash:             cash,
			SupplierId:       supplierID,
			DueDate:          dueDate,
			InvoiceNumber:    invoiceNumber,
			Kid:              kid,
			ProjectId:        projectID,
			InvoiceIssueDate: invoiceIssueDate,
			Paid:             paid,
			Currency:         currency,
			Lines: []api.DraftLineRequest{
				{
					Text:     description,
					Account:  account,
					VatType:  vatType,
					NetPrice: amountCents,
				},
			},
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointPurchaseDrafts, slug), req)
		if err != nil {
			return fmt.Errorf("creating purchase draft: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing draft ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Purchase draft created (ID: %d)", id))
		return nil
	},
}

var purchasesDraftsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a single purchase draft by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var draft api.PurchaseDraft
		_, err = client.Get(fmt.Sprintf(api.EndpointPurchaseDraft, slug, id), &draft)
		if err != nil {
			return fmt.Errorf("fetching purchase draft: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(draft)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("Draft ID", fmt.Sprintf("%d", draft.DraftId))
		table.AddRow("UUID", draft.Uuid)
		table.AddRow("Cash", BoolToYesNo(draft.Cash))
		table.AddRow("Supplier ID", fmt.Sprintf("%d", draft.SupplierId))
		table.AddRow("Due Date", draft.DueDate)
		table.AddRow("Invoice Number", draft.InvoiceNumber)
		table.AddRow("KID", draft.Kid)
		table.AddRow("Project ID", fmt.Sprintf("%d", draft.ProjectId))
		table.AddRow("Paid", BoolToYesNo(draft.Paid))
		table.Print()

		if len(draft.Lines) > 0 {
			fmt.Println()
			lineTable := output.NewTable("TEXT", "ACCOUNT", "VAT TYPE", "NET PRICE", "VAT", "GROSS")
			for _, l := range draft.Lines {
				lineTable.AddRow(
					l.Text,
					l.Account,
					l.VatType,
					output.FormatAmount(l.NetPrice),
					output.FormatAmount(l.Vat),
					output.FormatAmount(l.Gross),
				)
			}
			lineTable.Print()
		}

		return nil
	},
}

var purchasesDraftsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a purchase draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var existing api.PurchaseDraft
		endpoint := fmt.Sprintf(api.EndpointPurchaseDraft, slug, id)
		if _, err := client.Get(endpoint, &existing); err != nil {
			return fmt.Errorf("fetching draft for update: %w", err)
		}

		var existingLine api.DraftLineRequest
		if len(existing.Lines) > 0 {
			existingLine = existing.Lines[0]
		}

		req := api.PurchaseDraftRequest{
			Cash:          existing.Cash,
			SupplierId:    existing.SupplierId,
			DueDate:       existing.DueDate,
			InvoiceNumber: existing.InvoiceNumber,
			Kid:           existing.Kid,
			ProjectId:     existing.ProjectId,
			Paid:          existing.Paid,
			Lines: []api.DraftLineRequest{
				{
					Text:     existingLine.Text,
					Account:  existingLine.Account,
					VatType:  existingLine.VatType,
					NetPrice: existingLine.NetPrice,
				},
			},
		}

		if cmd.Flags().Changed("cash") {
			req.Cash, _ = cmd.Flags().GetBool("cash")
		}
		if cmd.Flags().Changed("supplier-id") {
			req.SupplierId, _ = cmd.Flags().GetInt64("supplier-id")
		}
		if cmd.Flags().Changed("due-date") {
			req.DueDate, _ = cmd.Flags().GetString("due-date")
		}
		if cmd.Flags().Changed("invoice-number") {
			req.InvoiceNumber, _ = cmd.Flags().GetString("invoice-number")
		}
		if cmd.Flags().Changed("kid") {
			req.Kid, _ = cmd.Flags().GetString("kid")
		}
		if cmd.Flags().Changed("project-id") {
			req.ProjectId, _ = cmd.Flags().GetInt64("project-id")
		}
		if cmd.Flags().Changed("invoice-issue-date") {
			req.InvoiceIssueDate, _ = cmd.Flags().GetString("invoice-issue-date")
		}
		if cmd.Flags().Changed("paid") {
			req.Paid, _ = cmd.Flags().GetBool("paid")
		}
		if cmd.Flags().Changed("currency") {
			req.Currency, _ = cmd.Flags().GetString("currency")
		}
		if cmd.Flags().Changed("description") {
			req.Lines[0].Text, _ = cmd.Flags().GetString("description")
		}
		if cmd.Flags().Changed("account") {
			req.Lines[0].Account, _ = cmd.Flags().GetString("account")
		}
		if cmd.Flags().Changed("vat-type") {
			req.Lines[0].VatType, _ = cmd.Flags().GetString("vat-type")
		}
		if cmd.Flags().Changed("amount") {
			amountStr, _ := cmd.Flags().GetString("amount")
			req.Lines[0].NetPrice, err = ParseAmountToCents(amountStr)
			if err != nil {
				return err
			}
		}

		_, err = client.Put(endpoint, req)
		if err != nil {
			return fmt.Errorf("updating purchase draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Purchase draft %d updated", id))
		return nil
	},
}

var purchasesDraftsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a purchase draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointPurchaseDraft, slug, id)
		if _, err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting purchase draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Purchase draft %d deleted", id))
		return nil
	},
}

var purchasesDraftsAttachCmd = &cobra.Command{
	Use:   "attach <id>",
	Short: "Attach a file to a purchase draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID %q: %w", args[0], err)
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

		endpoint := fmt.Sprintf(api.EndpointPurchaseDraftAttachments, slug, id)
		if err := UploadAttachment(client, endpoint, filePath, map[string]string{
			"filename": filepath.Base(filePath),
		}); err != nil {
			return fmt.Errorf("attaching to draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to draft %d", id))
		return nil
	},
}

var purchasesDraftsAttachmentsCmd = &cobra.Command{
	Use:   "attachments [id]",
	Short: "List attachments for a purchase draft",
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
		endpoint := fmt.Sprintf(api.EndpointPurchaseDraftAttachments, slug, id)
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

var purchasesDraftsFinalizeCmd = &cobra.Command{
	Use:   "finalize <id>",
	Short: "Finalize a purchase draft into a purchase",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		locationURL, err := client.PostEmpty(fmt.Sprintf(api.EndpointPurchaseDraftFinalize, slug, id))
		if err != nil {
			return fmt.Errorf("finalizing purchase draft: %w", err)
		}

		purchaseID, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing purchase ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Purchase draft %d finalized → Purchase %d created", id, purchaseID))
		return nil
	},
}

func init() {
	purchasesDraftsCreateCmd.Flags().Bool("cash", false, "Whether this is a cash purchase draft")
	purchasesDraftsCreateCmd.Flags().String("description", "", "Line description (required)")
	purchasesDraftsCreateCmd.Flags().String("account", "", "Expense account code (required)")
	purchasesDraftsCreateCmd.Flags().String("amount", "", "Amount in decimal format e.g. '1000.00' (required)")
	purchasesDraftsCreateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM' (required)")
	purchasesDraftsCreateCmd.Flags().Int64("supplier-id", 0, "Supplier contact ID (optional)")
	purchasesDraftsCreateCmd.Flags().String("due-date", "", "Due date YYYY-MM-DD (optional)")
	purchasesDraftsCreateCmd.Flags().String("invoice-number", "", "Invoice number (optional)")
	purchasesDraftsCreateCmd.Flags().String("kid", "", "KID number (optional)")
	purchasesDraftsCreateCmd.Flags().Int64("project-id", 0, "Project ID (optional)")
	purchasesDraftsCreateCmd.Flags().String("invoice-issue-date", "", "Invoice issue date YYYY-MM-DD (optional)")
	purchasesDraftsCreateCmd.Flags().Bool("paid", false, "Whether the draft is paid")
	purchasesDraftsCreateCmd.Flags().String("currency", "NOK", "Currency code")

	purchasesDraftsUpdateCmd.Flags().Bool("cash", false, "Whether this is a cash purchase draft")
	purchasesDraftsUpdateCmd.Flags().String("description", "", "Line description")
	purchasesDraftsUpdateCmd.Flags().String("account", "", "Expense account code")
	purchasesDraftsUpdateCmd.Flags().String("amount", "", "Amount in decimal format e.g. '1000.00'")
	purchasesDraftsUpdateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM'")
	purchasesDraftsUpdateCmd.Flags().Int64("supplier-id", 0, "Supplier contact ID")
	purchasesDraftsUpdateCmd.Flags().String("due-date", "", "Due date YYYY-MM-DD")
	purchasesDraftsUpdateCmd.Flags().String("invoice-number", "", "Invoice number")
	purchasesDraftsUpdateCmd.Flags().String("kid", "", "KID number")
	purchasesDraftsUpdateCmd.Flags().Int64("project-id", 0, "Project ID")
	purchasesDraftsUpdateCmd.Flags().String("invoice-issue-date", "", "Invoice issue date YYYY-MM-DD")
	purchasesDraftsUpdateCmd.Flags().Bool("paid", false, "Whether the draft is paid")
	purchasesDraftsUpdateCmd.Flags().String("currency", "", "Currency code")

	purchasesDraftsAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	purchasesDraftsAttachCmd.MarkFlagRequired("file")

	purchasesDraftsCmd.AddCommand(purchasesDraftsListCmd)
	purchasesDraftsCmd.AddCommand(purchasesDraftsCreateCmd)
	purchasesDraftsCmd.AddCommand(purchasesDraftsGetCmd)
	purchasesDraftsCmd.AddCommand(purchasesDraftsUpdateCmd)
	purchasesDraftsCmd.AddCommand(purchasesDraftsDeleteCmd)
	purchasesDraftsCmd.AddCommand(purchasesDraftsAttachCmd)
	purchasesDraftsCmd.AddCommand(purchasesDraftsFinalizeCmd)
	purchasesDraftsCmd.AddCommand(purchasesDraftsAttachmentsCmd)

	purchasesCmd.AddCommand(purchasesDraftsCmd)
}
