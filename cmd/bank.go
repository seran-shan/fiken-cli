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

var bankCmd = &cobra.Command{
	Use:   "bank",
	Short: "Manage bank accounts",
}

var bankListCmd = &cobra.Command{
	Use:   "list",
	Short: "List bank accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointBankAccounts, slug)

		var bankAccounts []api.BankAccount
		_, err = client.Get(endpoint, &bankAccounts)
		if err != nil {
			return fmt.Errorf("fetching bank accounts: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(bankAccounts)
		}

		if len(bankAccounts) == 0 {
			output.PrintInfo("No bank accounts found.")
			return nil
		}

		table := output.NewTable("ID", "NAME", "ACCOUNT", "BANK ACCOUNT", "TYPE", "ACTIVE")
		for _, ba := range bankAccounts {
			active := "Yes"
			if ba.Inactive {
				active = "No"
			}
			table.AddRow(
				fmt.Sprintf("%d", ba.BankAccountId),
				ba.Name,
				ba.AccountCode,
				ba.BankAccountNumber,
				ba.Type,
				active,
			)
		}
		table.Print()

		return nil
	},
}

var bankGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a bank account by ID",
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

		var bankAccount api.BankAccount
		endpoint := fmt.Sprintf(api.EndpointBankAccount, slug, id)
		if _, err := client.Get(endpoint, &bankAccount); err != nil {
			return fmt.Errorf("fetching bank account: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(bankAccount)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", bankAccount.BankAccountId))
		table.AddRow("Name", bankAccount.Name)
		table.AddRow("Account Code", bankAccount.AccountCode)
		table.AddRow("Bank Account Number", bankAccount.BankAccountNumber)
		table.AddRow("Type", bankAccount.Type)
		table.AddRow("IBAN", bankAccount.Iban)
		table.AddRow("BIC", bankAccount.Bic)
		table.AddRow("Foreign Service", bankAccount.ForeignService)
		table.AddRow("Inactive", BoolToYesNo(bankAccount.Inactive))
		table.Print()
		return nil
	},
}

var bankCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a bank account",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		bankAccountNumber, _ := cmd.Flags().GetString("bank-account-number")
		accountType, _ := cmd.Flags().GetString("type")
		bic, _ := cmd.Flags().GetString("bic")
		iban, _ := cmd.Flags().GetString("iban")

		var missing []string
		if name == "" {
			missing = append(missing, "--name")
		}
		if bankAccountNumber == "" {
			missing = append(missing, "--bank-account-number")
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

		req := api.BankAccountRequest{
			Name:              name,
			BankAccountNumber: bankAccountNumber,
			Type:              accountType,
			Bic:               bic,
			Iban:              iban,
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointBankAccounts, slug), req)
		if err != nil {
			return fmt.Errorf("creating bank account: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing bank account ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Bank account created (ID: %d)", id))
		return nil
	},
}

var bankBalancesCmd = &cobra.Command{
	Use:   "balances",
	Short: "List bank balances",
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		params := url.Values{}
		if date != "" {
			params.Set("date", date)
		}

		var balances []api.BankBalance
		_, err = client.GetWithParams(fmt.Sprintf(api.EndpointBankBalances, slug), params, &balances)
		if err != nil {
			return fmt.Errorf("fetching bank balances: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(balances)
		}

		if len(balances) == 0 {
			output.PrintInfo("No bank balances found.")
			return nil
		}

		table := output.NewTable("BANK ACCOUNT ID", "ACCOUNT CODE", "BALANCE", "DATE")
		for _, b := range balances {
			table.AddRow(
				fmt.Sprintf("%d", b.BankAccountId),
				b.AccountCode,
				output.FormatAmount(b.Balance),
				b.Date,
			)
		}
		table.Print()

		return nil
	},
}

func init() {
	bankCreateCmd.Flags().String("name", "", "Bank account name (required)")
	bankCreateCmd.Flags().String("bank-account-number", "", "Bank account number (required)")
	bankCreateCmd.Flags().String("type", "", "Account type: normal, tax, foreign, credit_card (optional)")
	bankCreateCmd.Flags().String("bic", "", "BIC/SWIFT code (optional)")
	bankCreateCmd.Flags().String("iban", "", "IBAN (optional)")

	bankBalancesCmd.Flags().String("date", "", "Balance date (YYYY-MM-DD)")

	bankCmd.AddCommand(bankListCmd)
	bankCmd.AddCommand(bankGetCmd)
	bankCmd.AddCommand(bankCreateCmd)
	bankCmd.AddCommand(bankBalancesCmd)
	rootCmd.AddCommand(bankCmd)
}
