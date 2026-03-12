package cmd

import (
	"fmt"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/auth"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var companiesCmd = &cobra.Command{
	Use:   "companies",
	Short: "List companies",
	Long:  "List all companies you have access to on Fiken.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		var companies []api.Company
		_, err = client.Get(api.EndpointCompanies, &companies)
		if err != nil {
			return fmt.Errorf("fetching companies: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(companies)
		}

		if len(companies) == 0 {
			output.PrintInfo("No companies found.")
			return nil
		}

		table := output.NewTable("NAME", "SLUG", "ORG.NR", "VAT TYPE")
		for _, c := range companies {
			table.AddRow(c.Name, c.Slug, c.OrganizationNumber, c.VatType)
		}
		table.Print()

		// Show default company hint
		cfg, _ := auth.LoadConfig()
		if cfg != nil && cfg.DefaultCompany == "" && len(companies) > 1 {
			fmt.Printf("\nTip: Set default company with 'fiken companies default <slug>'\n")
		}

		return nil
	},
}

var companiesDefaultCmd = &cobra.Command{
	Use:   "default [slug]",
	Short: "Set or show default company",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := auth.LoadConfig()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			if cfg.DefaultCompany != "" {
				fmt.Printf("Default company: %s\n", cfg.DefaultCompany)
			} else {
				output.PrintInfo("No default company set.")
			}
			return nil
		}

		cfg.DefaultCompany = args[0]
		if err := auth.SaveConfig(cfg); err != nil {
			return err
		}
		output.PrintSuccess(fmt.Sprintf("Default company set to '%s'", args[0]))
		return nil
	},
}

var companiesGetCmd = &cobra.Command{
	Use:   "get [slug]",
	Short: "Get a company by slug",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		var company api.Company
		if _, err := client.Get(fmt.Sprintf(api.EndpointCompany, args[0]), &company); err != nil {
			return fmt.Errorf("fetching company: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(company)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("Name", company.Name)
		table.AddRow("Slug", company.Slug)
		table.AddRow("Org Number", company.OrganizationNumber)
		table.AddRow("VAT Type", company.VatType)
		table.AddRow("Email", company.Email)
		table.AddRow("Phone", company.PhoneNumber)
		table.AddRow("Creation Date", company.CreationDate)
		table.AddRow("API Access", BoolToYesNo(company.HasApiAccess))
		table.AddRow("Test Company", BoolToYesNo(company.TestCompany))
		table.AddRow("Accounting Start Date", company.AccountingStartDate)
		table.AddRow("Street", company.Address.StreetAddress)
		table.AddRow("City", company.Address.City)
		table.AddRow("Post Code", company.Address.PostCode)
		table.AddRow("Country", company.Address.Country)
		table.Print()
		return nil
	},
}

func init() {
	companiesCmd.AddCommand(companiesDefaultCmd)
	companiesCmd.AddCommand(companiesGetCmd)
	rootCmd.AddCommand(companiesCmd)
}
