package cmd

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var journalCmd = &cobra.Command{
	Use:   "journal",
	Short: "Manage journal entries",
}

var journalListCmd = &cobra.Command{
	Use:   "list",
	Short: "List journal entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var params url.Values
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params = url.Values{}
			params.Set("date", date)
		}

		endpoint := fmt.Sprintf(api.EndpointJournalEntries, slug)
		entries, err := FetchAllPages[api.JournalEntry](client, endpoint, params, 25, 4)
		if err != nil {
			return fmt.Errorf("fetching journal entries: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(entries)
		}

		if len(entries) == 0 {
			output.PrintInfo("No journal entries found.")
			return nil
		}

		table := output.NewTable("ID", "DATE", "DESCRIPTION")
		for _, e := range entries {
			table.AddRow(
				fmt.Sprintf("%d", e.JournalEntryId),
				e.Date,
				e.Description,
			)
		}
		table.Print()

		fmt.Printf("\n%d journal entries\n", len(entries))
		return nil
	},
}

var journalGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a journal entry by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid journal entry ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var entry api.JournalEntry
		_, err = client.Get(fmt.Sprintf(api.EndpointJournalEntry, slug, id), &entry)
		if err != nil {
			return fmt.Errorf("fetching journal entry: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(entry)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", entry.JournalEntryId))
		table.AddRow("Date", entry.Date)
		table.AddRow("Description", entry.Description)
		table.Print()

		if len(entry.Lines) > 0 {
			fmt.Println()
			lineTable := output.NewTable("ACCOUNT", "DEBIT", "CREDIT")
			for _, l := range entry.Lines {
				lineTable.AddRow(
					l.Account,
					output.FormatAmount(l.DebitAmount),
					output.FormatAmount(l.CreditAmount),
				)
			}
			lineTable.Print()
		}

		return nil
	},
}

var journalAttachmentsCmd = &cobra.Command{
	Use:   "attachments [id]",
	Short: "List attachments for a journal entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid journal entry ID %q: %w", args[0], err)
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
		endpoint := fmt.Sprintf(api.EndpointJournalEntryAttachments, slug, id)
		_, err = client.Get(endpoint, &attachments)
		if err != nil {
			return fmt.Errorf("fetching attachments: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(attachments)
		}

		if len(attachments) == 0 {
			output.PrintInfo("No attachments found.")
			return nil
		}

		table := output.NewTable("ATTACHMENT ID", "FILENAME", "TYPE", "DATE")
		for _, a := range attachments {
			table.AddRow(
				fmt.Sprintf("%d", a.AttachmentId),
				a.Filename,
				a.Type,
				a.Date,
			)
		}
		table.Print()

		return nil
	},
}

var journalCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a general journal entry (fri postering)",
	Long: `Create a general journal entry (fri postering) with a single debit/credit pair.

Note: debit-amount is the net amount, credit-amount is the gross amount (including VAT).

Example:
  fiken journal create \
    --date 2025-01-15 \
    --description "Office supplies" \
    --debit-account 6800 \
    --debit-amount 1000.00 \
    --credit-account 1920 \
    --credit-amount 1250.00 \
    --debit-vat-code 3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")
		description, _ := cmd.Flags().GetString("description")
		debitAccount, _ := cmd.Flags().GetString("debit-account")
		debitAmountStr, _ := cmd.Flags().GetString("debit-amount")
		creditAccount, _ := cmd.Flags().GetString("credit-account")
		creditAmountStr, _ := cmd.Flags().GetString("credit-amount")
		debitVatCode, _ := cmd.Flags().GetInt64("debit-vat-code")
		creditVatCode, _ := cmd.Flags().GetInt64("credit-vat-code")
		open, _ := cmd.Flags().GetBool("open")

		debitAmount, err := ParseAmountToCents(debitAmountStr)
		if err != nil {
			return fmt.Errorf("invalid debit-amount: %w", err)
		}

		_, err = ParseAmountToCents(creditAmountStr)
		if err != nil {
			return fmt.Errorf("invalid credit-amount: %w", err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		req := api.GeneralJournalEntryRequest{
			Description: description,
			Open:        open,
			JournalEntries: []api.JournalEntryRequest{
				{
					Description: description,
					Date:        date,
					Lines: []api.JournalEntryLineRequest{
						{
							Amount:        debitAmount,
							DebitAccount:  debitAccount,
							DebitVatCode:  debitVatCode,
							CreditAccount: creditAccount,
							CreditVatCode: creditVatCode,
						},
					},
				},
			},
		}

		_, err = client.PostCreate(fmt.Sprintf(api.EndpointGeneralJournalEntries, slug), req)
		if err != nil {
			return err
		}

		output.PrintSuccess("Journal entry created")
		return nil
	},
}

var journalAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a document to a journal entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		entryID, _ := cmd.Flags().GetInt64("id")
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
		endpoint := fmt.Sprintf(api.EndpointJournalEntryAttachments, slug, entryID)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("attaching to journal entry: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to journal entry %d", entryID))
		return nil
	},
}

func init() {
	journalCreateCmd.Flags().String("date", "", "Entry date (yyyy-MM-dd)")
	journalCreateCmd.Flags().String("description", "", "Journal entry description (max 200 chars)")
	journalCreateCmd.Flags().String("debit-account", "", "Account to debit (e.g., 6800)")
	journalCreateCmd.Flags().String("debit-amount", "", "Net amount for debit in decimal (e.g., 1000.00)")
	journalCreateCmd.Flags().String("credit-account", "", "Account to credit (e.g., 1920)")
	journalCreateCmd.Flags().String("credit-amount", "", "Gross amount for credit in decimal, including VAT (e.g., 1250.00)")
	journalCreateCmd.Flags().Int64("debit-vat-code", 0, "VAT code integer for debit line")
	journalCreateCmd.Flags().Int64("credit-vat-code", 0, "VAT code integer for credit line")
	journalCreateCmd.Flags().Bool("open", false, "If true, entry can be deleted without reverse transaction")

	_ = journalCreateCmd.MarkFlagRequired("date")
	_ = journalCreateCmd.MarkFlagRequired("description")
	_ = journalCreateCmd.MarkFlagRequired("debit-account")
	_ = journalCreateCmd.MarkFlagRequired("debit-amount")
	_ = journalCreateCmd.MarkFlagRequired("credit-account")
	_ = journalCreateCmd.MarkFlagRequired("credit-amount")

	journalAttachCmd.Flags().Int64("id", 0, "Journal entry ID to attach to (required)")
	journalAttachCmd.MarkFlagRequired("id")
	journalAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	journalAttachCmd.MarkFlagRequired("file")

	journalListCmd.Flags().String("date", "", "Filter by date (YYYY-MM-DD)")

	journalCmd.AddCommand(journalListCmd)
	journalCmd.AddCommand(journalCreateCmd)
	journalCmd.AddCommand(journalGetCmd)
	journalCmd.AddCommand(journalAttachCmd)
	journalCmd.AddCommand(journalAttachmentsCmd)
	rootCmd.AddCommand(journalCmd)
}
