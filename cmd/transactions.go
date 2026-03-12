package cmd

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Manage transactions",
}

var transactionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List transactions",
	Long:  "List financial transactions for the selected company with optional date filters.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		params := url.Values{}
		params.Set("pageSize", "25")

		if v, _ := cmd.Flags().GetString("created-date"); v != "" {
			params.Set("createdDate", v)
		}
		if v, _ := cmd.Flags().GetString("last-modified"); v != "" {
			params.Set("lastModified", v)
		}

		endpoint := fmt.Sprintf(api.EndpointTransactions, slug)

		var transactions []api.Transaction
		page := 0
		for {
			params.Set("page", fmt.Sprintf("%d", page))
			var pageTransactions []api.Transaction
			pagination, err := client.GetWithParams(endpoint, params, &pageTransactions)
			if err != nil {
				return fmt.Errorf("fetching transactions: %w", err)
			}
			transactions = append(transactions, pageTransactions...)

			if pagination == nil || page+1 >= pagination.PageCount || len(pageTransactions) == 0 {
				break
			}
			page++
			if page >= 4 {
				break
			}
		}

		if jsonOutput {
			return output.PrintJSON(transactions)
		}

		if len(transactions) == 0 {
			output.PrintInfo("No transactions found.")
			return nil
		}

		table := output.NewTable("ID", "DATE", "TYPE", "DESCRIPTION", "DELETED")
		for _, t := range transactions {
			deleted := "No"
			if t.Deleted {
				deleted = "Yes"
			}
			table.AddRow(
				fmt.Sprintf("%d", t.TransactionId),
				t.CreatedDate,
				t.Type,
				t.Description,
				deleted,
			)
		}
		table.Print()

		fmt.Printf("\n%d transactions\n", len(transactions))
		return nil
	},
}

var transactionsDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Soft-delete a transaction",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid transaction ID %q: %w", args[0], err)
		}

		description, _ := cmd.Flags().GetString("description")

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		params := url.Values{}
		params.Set("description", description)

		err = client.PatchWithParams(fmt.Sprintf(api.EndpointTransactionDelete, slug, id), params)
		if err != nil {
			return fmt.Errorf("deleting transaction: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Transaction %d deleted", id))
		return nil
	},
}

var transactionsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a single transaction by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("transaction ID required")
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid transaction ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var transaction api.Transaction
		_, err = client.Get(fmt.Sprintf(api.EndpointTransaction, slug, id), &transaction)
		if err != nil {
			return fmt.Errorf("fetching transaction: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(transaction)
		}

		deleted := "No"
		if transaction.Deleted {
			deleted = "Yes"
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", transaction.TransactionId))
		table.AddRow("Date", transaction.CreatedDate)
		table.AddRow("Last Modified", transaction.LastModifiedDate)
		table.AddRow("Type", transaction.Type)
		table.AddRow("Description", transaction.Description)
		table.AddRow("Deleted", deleted)
		table.Print()

		return nil
	},
}

func init() {
	transactionsListCmd.Flags().String("created-date", "", "Filter by created date (YYYY-MM-DD)")
	transactionsListCmd.Flags().String("last-modified", "", "Filter by last modified date (YYYY-MM-DD)")

	transactionsDeleteCmd.Flags().String("description", "", "Deletion description (optional)")

	transactionsCmd.AddCommand(transactionsListCmd)
	transactionsCmd.AddCommand(transactionsGetCmd)
	transactionsCmd.AddCommand(transactionsDeleteCmd)
	rootCmd.AddCommand(transactionsCmd)
}
