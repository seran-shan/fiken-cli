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

var purchasesCmd = &cobra.Command{
	Use:   "purchases",
	Short: "Manage purchases",
	Long:  "List and manage purchases/expenses.",
}

var purchasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List purchases",
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

		endpoint := fmt.Sprintf(api.EndpointPurchases, slug)

		var purchases []api.Purchase
		page := 0
		for {
			params.Set("page", fmt.Sprintf("%d", page))
			var pagePurchases []api.Purchase
			pagination, err := client.GetWithParams(endpoint, params, &pagePurchases)
			if err != nil {
				return fmt.Errorf("fetching purchases: %w", err)
			}
			purchases = append(purchases, pagePurchases...)

			if pagination == nil || page+1 >= pagination.PageCount || len(pagePurchases) == 0 {
				break
			}
			page++
			// Only fetch first few pages by default
			if page >= 4 {
				break
			}
		}

		if jsonOutput {
			return output.PrintJSON(purchases)
		}

		if len(purchases) == 0 {
			output.PrintInfo("No purchases found.")
			return nil
		}

		table := output.NewTable("ID", "DATE", "KIND", "PAID", "AMOUNT", "IDENTIFIER")
		for _, p := range purchases {
			paid := "No"
			if p.Paid {
				paid = "Yes"
			}
			// Sum net amounts from lines
			var totalNet int64
			for _, l := range p.Lines {
				totalNet += l.NetAmount
			}
			table.AddRow(
				fmt.Sprintf("%d", p.PurchaseId),
				p.Date,
				p.Kind,
				paid,
				output.FormatAmount(totalNet),
				p.Identifier,
			)
		}
		table.Print()

		fmt.Printf("\n%d purchases\n", len(purchases))
		return nil
	},
}

var purchasesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a purchase",
	Long:  "Create a new purchase/expense with a single order line.",
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")
		kind, _ := cmd.Flags().GetString("kind")
		paid, _ := cmd.Flags().GetBool("paid")
		description, _ := cmd.Flags().GetString("description")
		account, _ := cmd.Flags().GetString("account")
		amountStr, _ := cmd.Flags().GetString("amount")
		vatType, _ := cmd.Flags().GetString("vat-type")
		currency, _ := cmd.Flags().GetString("currency")
		identifier, _ := cmd.Flags().GetString("identifier")
		supplierID, _ := cmd.Flags().GetInt64("supplier-id")
		paymentAccount, _ := cmd.Flags().GetString("payment-account")
		paymentDate, _ := cmd.Flags().GetString("payment-date")
		filePath, _ := cmd.Flags().GetString("file")

		var missing []string
		if date == "" {
			missing = append(missing, "--date")
		}
		if kind == "" {
			missing = append(missing, "--kind")
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

		if kind != "cash_purchase" && kind != "supplier" {
			return fmt.Errorf("--kind must be 'cash_purchase' or 'supplier', got %q", kind)
		}

		amountCents, err := parseAmountToCents(amountStr)
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

		purchaseReq := api.PurchaseRequest{
			Date:           date,
			Kind:           kind,
			Paid:           paid,
			Currency:       currency,
			Identifier:     identifier,
			SupplierId:     supplierID,
			PaymentAccount: paymentAccount,
			PaymentDate:    paymentDate,
			Lines: []api.OrderLineRequest{
				{
					Description: description,
					NetPrice:    amountCents,
					Account:     account,
					VatType:     vatType,
				},
			},
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointPurchases, slug), purchaseReq)
		if err != nil {
			return fmt.Errorf("creating purchase: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing purchase ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Purchase created (ID: %d)", id))

		if filePath != "" {
			output.PrintInfo("Note: --file attachment will be added in a future update. Use 'fiken purchases attach' to attach manually.")
		}

		return nil
	},
}

// parseAmountToCents converts a decimal string like "1000.50" to int64 cents (100050).
// Uses integer arithmetic only to avoid float64 precision issues.
func parseAmountToCents(s string) (int64, error) {
	s = strings.TrimSpace(s)
	parts := strings.SplitN(s, ".", 2)
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || whole < 0 {
		return 0, fmt.Errorf("invalid amount %q: must be a non-negative number like 100.00", s)
	}
	cents := whole * 100
	if len(parts) == 2 {
		dec := parts[1]
		switch len(dec) {
		case 1:
			d, _ := strconv.ParseInt(dec, 10, 64)
			cents += d * 10
		default:
			d, _ := strconv.ParseInt(dec[:2], 10, 64)
			cents += d
		}
	}
	return cents, nil
}

func init() {
	purchasesCmd.AddCommand(purchasesListCmd)
	purchasesCmd.AddCommand(purchasesCreateCmd)
	rootCmd.AddCommand(purchasesCmd)

	purchasesCreateCmd.Flags().String("date", "", "Purchase date (YYYY-MM-DD, required)")
	purchasesCreateCmd.Flags().String("kind", "", "Purchase kind: 'cash_purchase' or 'supplier' (required)")
	purchasesCreateCmd.Flags().Bool("paid", false, "Whether the purchase is paid")
	purchasesCreateCmd.Flags().String("description", "", "Line description (required)")
	purchasesCreateCmd.Flags().String("account", "", "Expense account code (required)")
	purchasesCreateCmd.Flags().String("amount", "", "Amount in decimal format e.g. '1000.00' (required)")
	purchasesCreateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM' (required)")
	purchasesCreateCmd.Flags().String("currency", "NOK", "Currency code")
	purchasesCreateCmd.Flags().String("identifier", "", "External identifier (optional)")
	purchasesCreateCmd.Flags().Int64("supplier-id", 0, "Supplier contact ID (optional)")
	purchasesCreateCmd.Flags().String("payment-account", "", "Payment account code (optional)")
	purchasesCreateCmd.Flags().String("payment-date", "", "Payment date YYYY-MM-DD (optional)")
	purchasesCreateCmd.Flags().String("file", "", "Path to receipt file to attach after creation (optional)")
}
