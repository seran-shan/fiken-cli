package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var salesPaymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "Manage sale payments",
}

var salesPaymentsListCmd = &cobra.Command{
	Use:   "list <sale-id>",
	Short: "List payments for a sale",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		saleID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid sale ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var payments []api.Payment
		if _, err := client.Get(fmt.Sprintf(api.EndpointSalePayments, slug, saleID), &payments); err != nil {
			return fmt.Errorf("fetching sale payments: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(payments)
		}

		if len(payments) == 0 {
			output.PrintInfo("No payments found.")
			return nil
		}

		table := output.NewTable("ID", "DATE", "ACCOUNT", "AMOUNT", "CURRENCY")
		for _, p := range payments {
			table.AddRow(
				fmt.Sprintf("%d", p.PaymentId),
				p.Date,
				p.Account,
				output.FormatAmount(p.Amount),
				p.Currency,
			)
		}
		table.Print()

		fmt.Printf("\n%d payments\n", len(payments))
		return nil
	},
}

var salesPaymentsCreateCmd = &cobra.Command{
	Use:   "create <sale-id>",
	Short: "Create a payment for a sale",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		saleID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid sale ID %q: %w", args[0], err)
		}

		date, _ := cmd.Flags().GetString("date")
		account, _ := cmd.Flags().GetString("account")
		amountStr, _ := cmd.Flags().GetString("amount")
		currency, _ := cmd.Flags().GetString("currency")

		var missing []string
		if date == "" {
			missing = append(missing, "--date")
		}
		if account == "" {
			missing = append(missing, "--account")
		}
		if amountStr == "" {
			missing = append(missing, "--amount")
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

		req := api.PaymentRequest{
			Date:     date,
			Account:  account,
			Amount:   amountCents,
			Currency: currency,
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointSalePayments, slug, saleID), req)
		if err != nil {
			return fmt.Errorf("creating payment: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing payment ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Payment created (ID: %d)", id))
		return nil
	},
}

var salesPaymentsGetCmd = &cobra.Command{
	Use:   "get <sale-id> <payment-id>",
	Short: "Get a specific payment for a sale",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		saleID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid sale ID %q: %w", args[0], err)
		}

		paymentID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid payment ID %q: %w", args[1], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var payment api.Payment
		if _, err := client.Get(fmt.Sprintf(api.EndpointSalePayment, slug, saleID, paymentID), &payment); err != nil {
			return fmt.Errorf("fetching payment: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(payment)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", payment.PaymentId))
		table.AddRow("Date", payment.Date)
		table.AddRow("Account", payment.Account)
		table.AddRow("Amount", output.FormatAmount(payment.Amount))
		table.AddRow("Amount (NOK)", output.FormatAmount(payment.AmountInNok))
		table.AddRow("Currency", payment.Currency)
		table.AddRow("Fee", output.FormatAmount(payment.Fee))
		table.Print()
		return nil
	},
}

func init() {
	salesPaymentsCreateCmd.Flags().String("date", "", "Payment date (YYYY-MM-DD, required)")
	salesPaymentsCreateCmd.Flags().String("account", "", "Payment account code (required)")
	salesPaymentsCreateCmd.Flags().String("amount", "", "Amount in decimal format e.g. '1000.00' (required)")
	salesPaymentsCreateCmd.Flags().String("currency", "NOK", "Currency code")

	salesPaymentsCmd.AddCommand(salesPaymentsListCmd)
	salesPaymentsCmd.AddCommand(salesPaymentsCreateCmd)
	salesPaymentsCmd.AddCommand(salesPaymentsGetCmd)
	salesCmd.AddCommand(salesPaymentsCmd)
}
