package cmd

import (
	"fmt"
	"net/url"
	"path/filepath"
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
			paid := BoolToYesNo(p.Paid)
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

		if filePath != "" {
			if err := ValidateFile(filePath); err != nil {
				return err
			}
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
			fields := map[string]string{
				"filename":        filepath.Base(filePath),
				"attachToPayment": "true",
				"attachToSale":    "true",
			}
			endpoint := fmt.Sprintf(api.EndpointPurchaseAttachments, slug, id)
			if attachErr := UploadAttachment(client, endpoint, filePath, fields); attachErr != nil {
				output.PrintError(fmt.Sprintf("Purchase created (ID: %d) but attachment failed: %v. Use 'fiken purchases attach --id %d --file %s' to retry.", id, attachErr, id, filePath))
				return nil
			}
			output.PrintSuccess(fmt.Sprintf("Receipt attached to purchase %d", id))
		}

		return nil
	},
}

var purchasesAttachmentsCmd = &cobra.Command{
	Use:   "attachments [id]",
	Short: "List attachments for a purchase",
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
		var attachments []api.Attachment
		endpoint := fmt.Sprintf(api.EndpointPurchaseAttachments, slug, id)
		if _, err := client.Get(endpoint, &attachments); err != nil {
			return fmt.Errorf("fetching attachments: %w", err)
		}
		if jsonOutput {
			return output.PrintJSON(attachments)
		}
		if len(attachments) == 0 {
			output.PrintInfo("No attachments found.")
			return nil
		}
		table := output.NewTable("IDENTIFIER", "TYPE", "COMMENT", "DOWNLOAD URL")
		for _, a := range attachments {
			table.AddRow(
				a.Identifier,
				a.Type,
				a.Comment,
				a.DownloadUrl,
			)
		}
		table.Print()
		return nil
	},
}

var purchasesAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a receipt/document to a purchase",
	RunE: func(cmd *cobra.Command, args []string) error {
		purchaseID, _ := cmd.Flags().GetInt64("id")
		filePath, _ := cmd.Flags().GetString("file")
		attachToPayment, _ := cmd.Flags().GetBool("attach-to-payment")
		attachToSale, _ := cmd.Flags().GetBool("attach-to-sale")

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
			"filename":        filepath.Base(filePath),
			"attachToPayment": fmt.Sprintf("%v", attachToPayment),
			"attachToSale":    fmt.Sprintf("%v", attachToSale),
		}
		endpoint := fmt.Sprintf(api.EndpointPurchaseAttachments, slug, purchaseID)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("attaching to purchase: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to purchase %d", purchaseID))
		return nil
	},
}

var purchasesGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a single purchase by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid purchase ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var purchase api.Purchase
		_, err = client.Get(fmt.Sprintf(api.EndpointPurchase, slug, id), &purchase)
		if err != nil {
			return fmt.Errorf("fetching purchase: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(purchase)
		}

		var totalNet, totalVat, totalGross int64
		for _, l := range purchase.Lines {
			totalNet += l.NetAmount
			totalVat += l.VatAmount
			totalGross += l.GrossAmount
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", purchase.PurchaseId))
		table.AddRow("Date", purchase.Date)
		table.AddRow("Due Date", purchase.DueDate)
		table.AddRow("Kind", purchase.Kind)
		table.AddRow("Supplier", purchase.Supplier.Name)
		table.AddRow("Currency", purchase.Currency)
		table.AddRow("Paid", BoolToYesNo(purchase.Paid))
		table.AddRow("Total Paid", output.FormatAmount(purchase.TotalPaid))
		table.AddRow("Identifier", purchase.Identifier)
		table.Print()

		if len(purchase.Lines) > 0 {
			fmt.Println()
			lineTable := output.NewTable("DESCRIPTION", "ACCOUNT", "NET", "VAT", "GROSS", "VAT TYPE")
			for _, l := range purchase.Lines {
				lineTable.AddRow(
					l.Description,
					l.Account,
					output.FormatAmount(l.NetAmount),
					output.FormatAmount(l.VatAmount),
					output.FormatAmount(l.GrossAmount),
					l.VatType,
				)
			}
			lineTable.Print()
		}

		return nil
	},
}

var purchasesDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Soft-delete a purchase",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid purchase ID %q: %w", args[0], err)
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

		err = client.PatchWithParams(fmt.Sprintf(api.EndpointPurchaseDelete, slug, id), params)
		if err != nil {
			return fmt.Errorf("deleting purchase: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Purchase %d deleted", id))
		return nil
	},
}

func init() {
	purchasesAttachCmd.Flags().Int64("id", 0, "Purchase ID to attach to (required)")
	purchasesAttachCmd.MarkFlagRequired("id")
	purchasesAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	purchasesAttachCmd.MarkFlagRequired("file")
	purchasesAttachCmd.Flags().Bool("attach-to-payment", true, "Whether this documents the payment")
	purchasesAttachCmd.Flags().Bool("attach-to-sale", true, "Whether this documents the purchase itself")

	purchasesDeleteCmd.Flags().String("description", "", "Deletion description (optional)")

	purchasesCmd.AddCommand(purchasesListCmd)
	purchasesCmd.AddCommand(purchasesCreateCmd)
	purchasesCmd.AddCommand(purchasesGetCmd)
	purchasesCmd.AddCommand(purchasesDeleteCmd)
	purchasesCmd.AddCommand(purchasesAttachCmd)
	purchasesCmd.AddCommand(purchasesAttachmentsCmd)
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
