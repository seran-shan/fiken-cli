package cmd

import (
	"fmt"

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
		fmt.Println("Not yet implemented.")
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

		debitAmount, err := parseAmountToCents(debitAmountStr)
		if err != nil {
			return fmt.Errorf("invalid debit-amount: %w", err)
		}

		_, err = parseAmountToCents(creditAmountStr)
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

	journalCmd.AddCommand(journalListCmd)
	journalCmd.AddCommand(journalCreateCmd)
	rootCmd.AddCommand(journalCmd)
}
