package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var balancesCmd = &cobra.Command{
	Use:   "balances",
	Short: "List account balances",
	Long:  "List account balances for the selected company.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointAccountBalances, slug)

		params := url.Values{}
		dateStr, _ := cmd.Flags().GetString("date")
		if dateStr == "" {
			dateStr = time.Now().Format("2006-01-02")
		}
		params.Set("date", dateStr)

		var balances []api.AccountBalance
		_, err = client.GetWithParams(endpoint, params, &balances)
		if err != nil {
			return fmt.Errorf("fetching balances: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(balances)
		}

		if len(balances) == 0 {
			output.PrintInfo("No account balances found.")
			return nil
		}

		table := output.NewTable("CODE", "NAME", "BALANCE")
		for _, b := range balances {
			table.AddRow(b.Account.Code, b.Account.Name, output.FormatAmount(b.Balance))
		}
		table.Print()

		return nil
	},
}

var balancesGetCmd = &cobra.Command{
	Use:   "get [account-code]",
	Short: "Get balance for a specific account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointAccountBalance, slug, args[0])

		var balance api.AccountBalance
		_, err = client.Get(endpoint, &balance)
		if err != nil {
			return fmt.Errorf("fetching balance: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(balance)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("Account Code", balance.Account.Code)
		table.AddRow("Account Name", balance.Account.Name)
		table.AddRow("Balance", output.FormatAmount(balance.Balance))
		table.Print()

		return nil
	},
}

func init() {
	balancesCmd.Flags().String("date", "", "Balance date (YYYY-MM-DD, default: today)")
	balancesCmd.AddCommand(balancesGetCmd)
	rootCmd.AddCommand(balancesCmd)
}
