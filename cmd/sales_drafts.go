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

var salesDraftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "Manage sale drafts",
}

var salesDraftsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sale drafts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointSaleDrafts, slug)
		drafts, err := FetchAllPages[api.SaleDraft](client, endpoint, nil, 25, 4)
		if err != nil {
			return fmt.Errorf("fetching sale drafts: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(drafts)
		}

		if len(drafts) == 0 {
			output.PrintInfo("No sale drafts found.")
			return nil
		}

		table := output.NewTable("DRAFT ID", "CASH", "CUSTOMER ID", "PAID")
		for _, d := range drafts {
			table.AddRow(
				fmt.Sprintf("%d", d.DraftId),
				BoolToYesNo(d.Cash),
				fmt.Sprintf("%d", d.CustomerId),
				BoolToYesNo(d.Paid),
			)
		}
		table.Print()

		fmt.Printf("\n%d sale drafts\n", len(drafts))
		return nil
	},
}

var salesDraftsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a sale draft",
	Long:  "Create a new sale draft with a single draft line.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cash, _ := cmd.Flags().GetBool("cash")
		description, _ := cmd.Flags().GetString("description")
		account, _ := cmd.Flags().GetString("account")
		amountStr, _ := cmd.Flags().GetString("amount")
		vatType, _ := cmd.Flags().GetString("vat-type")
		customerID, _ := cmd.Flags().GetInt64("customer-id")
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

		req := api.SaleDraftRequest{
			Cash:       cash,
			CustomerId: customerID,
			Paid:       paid,
			Currency:   currency,
			Lines: []api.DraftLineRequest{
				{
					Text:     description,
					Account:  account,
					VatType:  vatType,
					NetPrice: amountCents,
				},
			},
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointSaleDrafts, slug), req)
		if err != nil {
			return fmt.Errorf("creating sale draft: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing draft ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Sale draft created (ID: %d)", id))
		return nil
	},
}

var salesDraftsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a single sale draft by ID",
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

		var draft api.SaleDraft
		_, err = client.Get(fmt.Sprintf(api.EndpointSaleDraft, slug, id), &draft)
		if err != nil {
			return fmt.Errorf("fetching sale draft: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(draft)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("Draft ID", fmt.Sprintf("%d", draft.DraftId))
		table.AddRow("UUID", draft.Uuid)
		table.AddRow("Cash", BoolToYesNo(draft.Cash))
		table.AddRow("Customer ID", fmt.Sprintf("%d", draft.CustomerId))
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

var salesDraftsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a sale draft",
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

		var existing api.SaleDraft
		endpoint := fmt.Sprintf(api.EndpointSaleDraft, slug, id)
		if _, err := client.Get(endpoint, &existing); err != nil {
			return fmt.Errorf("fetching draft for update: %w", err)
		}

		// Start from existing values
		var existingLine api.DraftLineRequest
		if len(existing.Lines) > 0 {
			existingLine = existing.Lines[0]
		}

		req := api.SaleDraftRequest{
			Cash:       existing.Cash,
			CustomerId: existing.CustomerId,
			Paid:       existing.Paid,
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
		if cmd.Flags().Changed("customer-id") {
			req.CustomerId, _ = cmd.Flags().GetInt64("customer-id")
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
			return fmt.Errorf("updating sale draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Sale draft %d updated", id))
		return nil
	},
}

var salesDraftsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a sale draft",
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

		endpoint := fmt.Sprintf(api.EndpointSaleDraft, slug, id)
		if _, err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting sale draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Sale draft %d deleted", id))
		return nil
	},
}

var salesDraftsAttachCmd = &cobra.Command{
	Use:   "attach <id>",
	Short: "Attach a file to a sale draft",
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

		endpoint := fmt.Sprintf(api.EndpointSaleDraftAttachments, slug, id)
		if err := UploadAttachment(client, endpoint, filePath, map[string]string{
			"filename": filepath.Base(filePath),
		}); err != nil {
			return fmt.Errorf("attaching to draft: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to draft %d", id))
		return nil
	},
}

var salesDraftsFinalizeCmd = &cobra.Command{
	Use:   "finalize <id>",
	Short: "Finalize a sale draft into a sale",
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

		locationURL, err := client.PostEmpty(fmt.Sprintf(api.EndpointSaleDraftFinalize, slug, id))
		if err != nil {
			return fmt.Errorf("finalizing sale draft: %w", err)
		}

		saleID, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing sale ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Sale draft %d finalized → Sale %d created", id, saleID))
		return nil
	},
}

func init() {
	// Create flags
	salesDraftsCreateCmd.Flags().Bool("cash", false, "Whether this is a cash sale draft")
	salesDraftsCreateCmd.Flags().String("description", "", "Line description (required)")
	salesDraftsCreateCmd.Flags().String("account", "", "Revenue account code (required)")
	salesDraftsCreateCmd.Flags().String("amount", "", "Amount in decimal format e.g. '1000.00' (required)")
	salesDraftsCreateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM' (required)")
	salesDraftsCreateCmd.Flags().Int64("customer-id", 0, "Customer contact ID (optional)")
	salesDraftsCreateCmd.Flags().Bool("paid", false, "Whether the draft is paid")
	salesDraftsCreateCmd.Flags().String("currency", "NOK", "Currency code")

	// Update flags (same as create)
	salesDraftsUpdateCmd.Flags().Bool("cash", false, "Whether this is a cash sale draft")
	salesDraftsUpdateCmd.Flags().String("description", "", "Line description")
	salesDraftsUpdateCmd.Flags().String("account", "", "Revenue account code")
	salesDraftsUpdateCmd.Flags().String("amount", "", "Amount in decimal format e.g. '1000.00'")
	salesDraftsUpdateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM'")
	salesDraftsUpdateCmd.Flags().Int64("customer-id", 0, "Customer contact ID")
	salesDraftsUpdateCmd.Flags().Bool("paid", false, "Whether the draft is paid")
	salesDraftsUpdateCmd.Flags().String("currency", "", "Currency code")

	// Attach flags
	salesDraftsAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	salesDraftsAttachCmd.MarkFlagRequired("file")

	// Register subcommands on salesDraftsCmd
	salesDraftsCmd.AddCommand(salesDraftsListCmd)
	salesDraftsCmd.AddCommand(salesDraftsCreateCmd)
	salesDraftsCmd.AddCommand(salesDraftsGetCmd)
	salesDraftsCmd.AddCommand(salesDraftsUpdateCmd)
	salesDraftsCmd.AddCommand(salesDraftsDeleteCmd)
	salesDraftsCmd.AddCommand(salesDraftsAttachCmd)
	salesDraftsCmd.AddCommand(salesDraftsFinalizeCmd)

	// Register drafts on salesCmd
	salesCmd.AddCommand(salesDraftsCmd)
}
