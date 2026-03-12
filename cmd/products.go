package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Manage products",
	Long:  "List, create, get, update, and delete products.",
}

var productsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List products",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointProducts, slug)
		products, err := FetchAllPages[api.Product](client, endpoint, nil, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching products: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(products)
		}

		if len(products) == 0 {
			output.PrintInfo("No products found.")
			return nil
		}

		table := output.NewTable("ID", "NAME", "UNIT PRICE", "INCOME ACCOUNT", "VAT TYPE", "ACTIVE", "PRODUCT #")
		for _, p := range products {
			table.AddRow(
				fmt.Sprintf("%d", p.ProductId),
				p.Name,
				output.FormatAmount(p.UnitPrice),
				p.IncomeAccount,
				p.VatType,
				BoolToYesNo(p.Active),
				p.ProductNumber,
			)
		}
		table.Print()
		fmt.Printf("\n%d products\n", len(products))
		return nil
	},
}

var productsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new product",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		incomeAccount, _ := cmd.Flags().GetString("income-account")
		vatType, _ := cmd.Flags().GetString("vat-type")
		unitPriceStr, _ := cmd.Flags().GetString("unit-price")
		active, _ := cmd.Flags().GetBool("active")
		productNumber, _ := cmd.Flags().GetString("product-number")
		stock, _ := cmd.Flags().GetFloat64("stock")
		note, _ := cmd.Flags().GetString("note")

		var missing []string
		if name == "" {
			missing = append(missing, "--name")
		}
		if incomeAccount == "" {
			missing = append(missing, "--income-account")
		}
		if vatType == "" {
			missing = append(missing, "--vat-type")
		}
		if len(missing) > 0 {
			return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
		}

		var unitPrice int64
		if unitPriceStr != "" {
			var err error
			unitPrice, err = ParseAmountToCents(unitPriceStr)
			if err != nil {
				return err
			}
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		req := api.ProductRequest{
			Name:          name,
			IncomeAccount: incomeAccount,
			VatType:       vatType,
			UnitPrice:     unitPrice,
			Active:        active,
			ProductNumber: productNumber,
			Stock:         stock,
			Note:          note,
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointProducts, slug), req)
		if err != nil {
			return fmt.Errorf("creating product: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing product ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Product created (ID: %d)", id))
		return nil
	},
}

var productsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a product by ID",
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

		var product api.Product
		endpoint := fmt.Sprintf(api.EndpointProduct, slug, id)
		if _, err := client.Get(endpoint, &product); err != nil {
			return fmt.Errorf("fetching product: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(product)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", product.ProductId))
		table.AddRow("Name", product.Name)
		table.AddRow("Unit Price", output.FormatAmount(product.UnitPrice))
		table.AddRow("Income Account", product.IncomeAccount)
		table.AddRow("VAT Type", product.VatType)
		table.AddRow("Active", BoolToYesNo(product.Active))
		table.AddRow("Product Number", product.ProductNumber)
		table.AddRow("Stock", fmt.Sprintf("%.2f", product.Stock))
		table.AddRow("Note", product.Note)
		table.Print()
		return nil
	},
}

var productsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a product",
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

		var existing api.Product
		endpoint := fmt.Sprintf(api.EndpointProduct, slug, id)
		if _, err := client.Get(endpoint, &existing); err != nil {
			return fmt.Errorf("fetching product for update: %w", err)
		}

		req := api.ProductRequest{
			Name:          existing.Name,
			UnitPrice:     existing.UnitPrice,
			IncomeAccount: existing.IncomeAccount,
			VatType:       existing.VatType,
			Active:        existing.Active,
			ProductNumber: existing.ProductNumber,
			Stock:         existing.Stock,
			Note:          existing.Note,
		}

		if cmd.Flags().Changed("name") {
			req.Name, _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("income-account") {
			req.IncomeAccount, _ = cmd.Flags().GetString("income-account")
		}
		if cmd.Flags().Changed("vat-type") {
			req.VatType, _ = cmd.Flags().GetString("vat-type")
		}
		if cmd.Flags().Changed("unit-price") {
			unitPriceStr, _ := cmd.Flags().GetString("unit-price")
			req.UnitPrice, err = ParseAmountToCents(unitPriceStr)
			if err != nil {
				return err
			}
		}
		if cmd.Flags().Changed("active") {
			req.Active, _ = cmd.Flags().GetBool("active")
		}
		if cmd.Flags().Changed("product-number") {
			req.ProductNumber, _ = cmd.Flags().GetString("product-number")
		}
		if cmd.Flags().Changed("stock") {
			req.Stock, _ = cmd.Flags().GetFloat64("stock")
		}
		if cmd.Flags().Changed("note") {
			req.Note, _ = cmd.Flags().GetString("note")
		}

		_, err = client.Put(endpoint, req)
		if err != nil {
			return fmt.Errorf("updating product: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Product %d updated", id))
		return nil
	},
}

var productsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a product",
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

		endpoint := fmt.Sprintf(api.EndpointProduct, slug, id)
		if _, err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting product: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Product %d deleted", id))
		return nil
	},
}

func init() {
	productsCreateCmd.Flags().String("name", "", "Product name (required)")
	productsCreateCmd.Flags().String("income-account", "", "Income account code, e.g. 3000 (required)")
	productsCreateCmd.Flags().String("vat-type", "", "VAT type, e.g. HIGH, NONE, MEDIUM (required)")
	productsCreateCmd.Flags().String("unit-price", "", "Unit price in decimal format, e.g. 100.00")
	productsCreateCmd.Flags().Bool("active", true, "Whether the product is active")
	productsCreateCmd.Flags().String("product-number", "", "Custom product number")
	productsCreateCmd.Flags().Float64("stock", 0, "Stock quantity")
	productsCreateCmd.Flags().String("note", "", "Internal note")

	productsUpdateCmd.Flags().String("name", "", "Product name")
	productsUpdateCmd.Flags().String("income-account", "", "Income account code")
	productsUpdateCmd.Flags().String("vat-type", "", "VAT type")
	productsUpdateCmd.Flags().String("unit-price", "", "Unit price in decimal format")
	productsUpdateCmd.Flags().Bool("active", true, "Whether the product is active")
	productsUpdateCmd.Flags().String("product-number", "", "Custom product number")
	productsUpdateCmd.Flags().Float64("stock", 0, "Stock quantity")
	productsUpdateCmd.Flags().String("note", "", "Internal note")

	productsCmd.AddCommand(productsListCmd)
	productsCmd.AddCommand(productsCreateCmd)
	productsCmd.AddCommand(productsGetCmd)
	productsCmd.AddCommand(productsUpdateCmd)
	productsCmd.AddCommand(productsDeleteCmd)
	rootCmd.AddCommand(productsCmd)
}
