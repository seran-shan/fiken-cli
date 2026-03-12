package cmd

import (
	"fmt"
	"strconv"
	"strings"

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

var offersCounterCmd = &cobra.Command{
	Use:   "counter",
	Short: "Get or set the offer counter",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}
		var counter api.OfferCounter
		if _, err := client.Get(fmt.Sprintf(api.EndpointOfferCounter, slug), &counter); err != nil {
			return fmt.Errorf("fetching offer counter: %w", err)
		}
		if jsonOutput {
			return output.PrintJSON(counter)
		}
		fmt.Printf("Offer counter: %d\n", counter.Value)
		return nil
	},
}

var offersCounterSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the offer counter",
	RunE: func(cmd *cobra.Command, args []string) error {
		value, _ := cmd.Flags().GetInt64("value")
		client, err := getClient()
		if err != nil {
			return err
		}
		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}
		req := api.OfferCounter{Value: int32(value)}
		if err := client.Post(fmt.Sprintf(api.EndpointOfferCounter, slug), req, nil); err != nil {
			return fmt.Errorf("setting offer counter: %w", err)
		}
		output.PrintSuccess(fmt.Sprintf("Offer counter set to %d", value))
		return nil
	},
}

var offersSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an offer",
	RunE: func(cmd *cobra.Command, args []string) error {
		offerID, _ := cmd.Flags().GetInt64("offer-id")
		if offerID == 0 {
			return fmt.Errorf("missing required flag: --offer-id")
		}
		methodStr, _ := cmd.Flags().GetString("method")
		recipientName, _ := cmd.Flags().GetString("recipient-name")
		recipientEmail, _ := cmd.Flags().GetString("recipient-email")
		message, _ := cmd.Flags().GetString("message")
		includeAttachments, _ := cmd.Flags().GetBool("include-document-attachments")

		methods := strings.Split(methodStr, ",")

		client, err := getClient()
		if err != nil {
			return err
		}
		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		req := api.SendOfferRequest{
			OfferId:                    offerID,
			Method:                     methods,
			RecipientName:              recipientName,
			RecipientEmail:             recipientEmail,
			Message:                    message,
			IncludeDocumentAttachments: includeAttachments,
		}

		if err := client.Post(fmt.Sprintf(api.EndpointOfferSend, slug), req, nil); err != nil {
			return fmt.Errorf("sending offer: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Offer %d sent", offerID))
		return nil
	},
}

func init() {
	offersCounterSetCmd.Flags().Int64("value", 0, "Counter value to set (required)")
	offersCounterSetCmd.MarkFlagRequired("value")

	offersSendCmd.Flags().Int64("offer-id", 0, "Offer ID to send (required)")
	offersSendCmd.Flags().String("method", "auto", "Delivery method (comma-separated: auto,email,ehf,efaktura,vipps,sms,letter)")
	offersSendCmd.Flags().String("recipient-name", "", "Recipient name (optional)")
	offersSendCmd.Flags().String("recipient-email", "", "Recipient email (optional)")
	offersSendCmd.Flags().String("message", "", "Message to include (optional)")
	offersSendCmd.Flags().Bool("include-document-attachments", false, "Include document attachments")

	offersCmd.AddCommand(offersListCmd)
	offersCmd.AddCommand(offersGetCmd)
	offersCmd.AddCommand(offersCounterCmd)
	offersCounterCmd.AddCommand(offersCounterSetCmd)
	offersCmd.AddCommand(offersSendCmd)
	rootCmd.AddCommand(offersCmd)
}
