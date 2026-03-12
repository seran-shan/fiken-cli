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

var salesCmd = &cobra.Command{
	Use:   "sales",
	Short: "Manage sales",
}

var salesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sales",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointSales, slug)
		sales, err := FetchAllPages[api.Sale](client, endpoint, nil, 25, 4)
		if err != nil {
			return fmt.Errorf("fetching sales: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(sales)
		}

		if len(sales) == 0 {
			output.PrintInfo("No sales found.")
			return nil
		}

		table := output.NewTable("ID", "DATE", "KIND", "CUSTOMER", "NET", "VAT", "GROSS", "PAID")
		for _, s := range sales {
			var totalNet, totalVat, totalGross int64
			for _, l := range s.Lines {
				totalNet += l.NetAmount
				totalVat += l.VatAmount
				totalGross += l.GrossAmount
			}
			table.AddRow(
				fmt.Sprintf("%d", s.SaleId),
				s.Date,
				s.Kind,
				s.Customer.Name,
				output.FormatAmount(totalNet),
				output.FormatAmount(totalVat),
				output.FormatAmount(totalGross),
				BoolToYesNo(s.Paid),
			)
		}
		table.Print()

		fmt.Printf("\n%d sales\n", len(sales))
		return nil
	},
}

var salesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a sale",
	Long:  "Create a new sale with a single order line.",
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")
		kind, _ := cmd.Flags().GetString("kind")
		description, _ := cmd.Flags().GetString("description")
		account, _ := cmd.Flags().GetString("account")
		amountStr, _ := cmd.Flags().GetString("amount")
		vatType, _ := cmd.Flags().GetString("vat-type")
		currency, _ := cmd.Flags().GetString("currency")
		customerID, _ := cmd.Flags().GetInt64("customer-id")
		paid, _ := cmd.Flags().GetBool("paid")
		paymentDate, _ := cmd.Flags().GetString("payment-date")
		paymentAccount, _ := cmd.Flags().GetString("payment-account")

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

		if kind != "cash_sale" && kind != "invoice" && kind != "external_invoice" {
			return fmt.Errorf("--kind must be 'cash_sale', 'invoice', or 'external_invoice', got %q", kind)
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

		saleReq := api.SaleRequest{
			Date:           date,
			Kind:           kind,
			Currency:       currency,
			Paid:           paid,
			CustomerId:     customerID,
			PaymentDate:    paymentDate,
			PaymentAccount: paymentAccount,
			Lines: []api.OrderLineRequest{
				{
					Description: description,
					NetPrice:    amountCents,
					Account:     account,
					VatType:     vatType,
				},
			},
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointSales, slug), saleReq)
		if err != nil {
			return fmt.Errorf("creating sale: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing sale ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Sale created (ID: %d)", id))
		return nil
	},
}

var salesGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a single sale by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
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

		var sale api.Sale
		_, err = client.Get(fmt.Sprintf(api.EndpointSale, slug, id), &sale)
		if err != nil {
			return fmt.Errorf("fetching sale: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(sale)
		}

		var totalNet, totalVat, totalGross int64
		for _, l := range sale.Lines {
			totalNet += l.NetAmount
			totalVat += l.VatAmount
			totalGross += l.GrossAmount
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", sale.SaleId))
		table.AddRow("Date", sale.Date)
		table.AddRow("Kind", sale.Kind)
		table.AddRow("Customer", sale.Customer.Name)
		table.AddRow("Currency", sale.Currency)
		table.AddRow("Due Date", sale.DueDate)
		table.AddRow("Net", output.FormatAmount(totalNet))
		table.AddRow("VAT", output.FormatAmount(totalVat))
		table.AddRow("Gross", output.FormatAmount(totalGross))
		table.AddRow("Paid", BoolToYesNo(sale.Paid))
		table.AddRow("Total Paid", output.FormatAmount(sale.TotalPaid))
		table.Print()

		if len(sale.Lines) > 0 {
			fmt.Println()
			lineTable := output.NewTable("DESCRIPTION", "ACCOUNT", "NET", "VAT", "GROSS", "VAT TYPE")
			for _, l := range sale.Lines {
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

var salesDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Soft-delete a sale",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid sale ID %q: %w", args[0], err)
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

		err = client.PatchWithParams(fmt.Sprintf(api.EndpointSaleDelete, slug, id), params)
		if err != nil {
			return fmt.Errorf("deleting sale: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Sale %d deleted", id))
		return nil
	},
}

var salesSettleCmd = &cobra.Command{
	Use:   "settle [id]",
	Short: "Mark a sale as settled",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid sale ID %q: %w", args[0], err)
		}

		settledDate, _ := cmd.Flags().GetString("settled-date")
		if settledDate == "" {
			return fmt.Errorf("--settled-date is required")
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		params := url.Values{}
		params.Set("settledDate", settledDate)

		err = client.PatchWithParams(fmt.Sprintf(api.EndpointSaleSettle, slug, id), params)
		if err != nil {
			return fmt.Errorf("settling sale: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Sale %d settled", id))
		return nil
	},
}

var salesAttachmentsCmd = &cobra.Command{
	Use:   "attachments [id]",
	Short: "List attachments for a sale",
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
		endpoint := fmt.Sprintf(api.EndpointSaleAttachments, slug, id)
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
		table := output.NewTable("ATTACHMENT ID", "FILENAME", "TYPE", "DATE")
		for _, a := range attachments {
			table.AddRow(
				fmt.Sprintf("%d", a.AttachmentId),
				a.Filename,
				a.Type,
				a.Date,
			)
		}
		table.Print()
		return nil
	},
}

var salesAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a receipt/document to a sale",
	RunE: func(cmd *cobra.Command, args []string) error {
		saleID, _ := cmd.Flags().GetInt64("id")
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
		endpoint := fmt.Sprintf(api.EndpointSaleAttachments, slug, saleID)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("attaching to sale: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to sale %d", saleID))
		return nil
	},
}

func init() {
	salesAttachCmd.Flags().Int64("id", 0, "Sale ID to attach to (required)")
	salesAttachCmd.MarkFlagRequired("id")
	salesAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	salesAttachCmd.MarkFlagRequired("file")
	salesAttachCmd.Flags().Bool("attach-to-payment", true, "Whether this documents the payment")
	salesAttachCmd.Flags().Bool("attach-to-sale", true, "Whether this documents the sale itself")

	salesCreateCmd.Flags().String("date", "", "Sale date (YYYY-MM-DD, required)")
	salesCreateCmd.Flags().String("kind", "", "Sale kind: 'cash_sale', 'invoice', or 'external_invoice' (required)")
	salesCreateCmd.Flags().String("description", "", "Line description (required)")
	salesCreateCmd.Flags().String("account", "", "Revenue account code (required)")
	salesCreateCmd.Flags().String("amount", "", "Amount in decimal format e.g. '1000.00' (required)")
	salesCreateCmd.Flags().String("vat-type", "", "VAT type e.g. 'HIGH', 'NONE', 'MEDIUM' (required)")
	salesCreateCmd.Flags().String("currency", "NOK", "Currency code")
	salesCreateCmd.Flags().Int64("customer-id", 0, "Customer contact ID (optional)")
	salesCreateCmd.Flags().Bool("paid", false, "Whether the sale is paid")
	salesCreateCmd.Flags().String("payment-date", "", "Payment date YYYY-MM-DD (optional)")
	salesCreateCmd.Flags().String("payment-account", "", "Payment account code (optional)")

	salesDeleteCmd.Flags().String("description", "", "Deletion description (optional)")

	salesSettleCmd.Flags().String("settled-date", "", "Settlement date YYYY-MM-DD (required)")
	salesSettleCmd.MarkFlagRequired("settled-date")

	salesCmd.AddCommand(salesListCmd)
	salesCmd.AddCommand(salesCreateCmd)
	salesCmd.AddCommand(salesGetCmd)
	salesCmd.AddCommand(salesDeleteCmd)
	salesCmd.AddCommand(salesSettleCmd)
	salesCmd.AddCommand(salesAttachCmd)
	salesCmd.AddCommand(salesAttachmentsCmd)
	rootCmd.AddCommand(salesCmd)
}
