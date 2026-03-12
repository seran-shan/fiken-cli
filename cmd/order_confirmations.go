package cmd

import (
	"fmt"
	"strconv"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var orderConfirmationsCmd = &cobra.Command{
	Use:   "order-confirmations",
	Short: "Manage order confirmations",
}

var orderConfirmationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List order confirmations",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointOrderConfirmations, slug)
		confirmations, err := FetchAllPages[api.OrderConfirmation](client, endpoint, nil, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching order confirmations: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(confirmations)
		}

		if len(confirmations) == 0 {
			output.PrintInfo("No order confirmations found.")
			return nil
		}

		table := output.NewTable("ID", "DATE", "CUSTOMER", "NET", "VAT", "GROSS", "CURRENCY")
		for _, c := range confirmations {
			table.AddRow(
				fmt.Sprintf("%d", c.ConfirmationId),
				c.Date,
				c.Customer.Name,
				output.FormatAmount(c.Net),
				output.FormatAmount(c.Vat),
				output.FormatAmount(c.Gross),
				c.Currency,
			)
		}
		table.Print()
		fmt.Printf("\n%d order confirmations\n", len(confirmations))
		return nil
	},
}

var orderConfirmationsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an order confirmation by ID",
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

		var confirmation api.OrderConfirmation
		endpoint := fmt.Sprintf(api.EndpointOrderConfirmation, slug, id)
		if _, err := client.Get(endpoint, &confirmation); err != nil {
			return fmt.Errorf("fetching order confirmation: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(confirmation)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", confirmation.ConfirmationId))
		table.AddRow("Date", confirmation.Date)
		table.AddRow("Customer", confirmation.Customer.Name)
		table.AddRow("Net", output.FormatAmount(confirmation.Net))
		table.AddRow("VAT", output.FormatAmount(confirmation.Vat))
		table.AddRow("Gross", output.FormatAmount(confirmation.Gross))
		table.AddRow("Currency", confirmation.Currency)
		table.Print()

		if len(confirmation.Lines) > 0 {
			fmt.Println()
			lines := output.NewTable("DESCRIPTION", "ACCOUNT", "NET", "VAT", "GROSS", "VAT TYPE")
			for _, l := range confirmation.Lines {
				lines.AddRow(
					l.Description,
					l.Account,
					output.FormatAmount(l.NetAmount),
					output.FormatAmount(l.VatAmount),
					output.FormatAmount(l.GrossAmount),
					l.VatType,
				)
			}
			lines.Print()
		}
		return nil
	},
}

func init() {
	orderConfirmationsCmd.AddCommand(orderConfirmationsListCmd)
	orderConfirmationsCmd.AddCommand(orderConfirmationsGetCmd)
	rootCmd.AddCommand(orderConfirmationsCmd)
}
