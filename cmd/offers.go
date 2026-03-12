package cmd

import (
	"fmt"
	"strconv"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var offersCmd = &cobra.Command{
	Use:   "offers",
	Short: "Manage offers",
}

var offersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List offers",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointOffers, slug)
		offers, err := FetchAllPages[api.Offer](client, endpoint, nil, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching offers: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(offers)
		}

		if len(offers) == 0 {
			output.PrintInfo("No offers found.")
			return nil
		}

		table := output.NewTable("ID", "OFFER NUMBER", "DATE", "CUSTOMER", "NET", "VAT", "GROSS", "CURRENCY")
		for _, o := range offers {
			table.AddRow(
				fmt.Sprintf("%d", o.OfferId),
				fmt.Sprintf("%d", o.OfferNumber),
				o.Date,
				o.Customer.Name,
				output.FormatAmount(o.Net),
				output.FormatAmount(o.Vat),
				output.FormatAmount(o.Gross),
				o.Currency,
			)
		}
		table.Print()
		fmt.Printf("\n%d offers\n", len(offers))
		return nil
	},
}

var offersGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an offer by ID",
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

		var offer api.Offer
		endpoint := fmt.Sprintf(api.EndpointOffer, slug, id)
		if _, err := client.Get(endpoint, &offer); err != nil {
			return fmt.Errorf("fetching offer: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(offer)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", offer.OfferId))
		table.AddRow("Offer Number", fmt.Sprintf("%d", offer.OfferNumber))
		table.AddRow("Date", offer.Date)
		table.AddRow("Customer", offer.Customer.Name)
		table.AddRow("Net", output.FormatAmount(offer.Net))
		table.AddRow("VAT", output.FormatAmount(offer.Vat))
		table.AddRow("Gross", output.FormatAmount(offer.Gross))
		table.AddRow("Currency", offer.Currency)
		table.Print()

		if len(offer.Lines) > 0 {
			fmt.Println()
			lines := output.NewTable("DESCRIPTION", "ACCOUNT", "NET", "VAT", "GROSS", "VAT TYPE")
			for _, l := range offer.Lines {
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
	offersCmd.AddCommand(offersListCmd)
	offersCmd.AddCommand(offersGetCmd)
	rootCmd.AddCommand(offersCmd)
}
