package cmd

import (
	"fmt"
	"strconv"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage contacts",
	Long:  "List, create, get, update, and delete contacts.",
}

var contactsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List contacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointContacts, slug)
		contacts, err := FetchAllPages[api.Contact](client, endpoint, nil, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching contacts: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(contacts)
		}

		if len(contacts) == 0 {
			output.PrintInfo("No contacts found.")
			return nil
		}

		table := output.NewTable("ID", "NAME", "EMAIL", "PHONE", "CUSTOMER", "SUPPLIER", "ORG NUMBER")
		for _, c := range contacts {
			table.AddRow(
				fmt.Sprintf("%d", c.ContactId),
				c.Name,
				c.Email,
				c.PhoneNumber,
				BoolToYesNo(c.Customer),
				BoolToYesNo(c.Supplier),
				c.OrganizationNumber,
			)
		}
		table.Print()
		fmt.Printf("\n%d contacts\n", len(contacts))
		return nil
	},
}

var contactsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new contact",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		phone, _ := cmd.Flags().GetString("phone")
		orgNumber, _ := cmd.Flags().GetString("org-number")
		customer, _ := cmd.Flags().GetBool("customer")
		supplier, _ := cmd.Flags().GetBool("supplier")
		language, _ := cmd.Flags().GetString("language")
		memberNumber, _ := cmd.Flags().GetInt64("member-number")
		street, _ := cmd.Flags().GetString("street")
		city, _ := cmd.Flags().GetString("city")
		postcode, _ := cmd.Flags().GetString("postcode")
		country, _ := cmd.Flags().GetString("country")

		if name == "" {
			return fmt.Errorf("missing required flag: --name")
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		req := api.ContactRequest{
			Name:               name,
			Email:              email,
			PhoneNumber:        phone,
			OrganizationNumber: orgNumber,
			Customer:           customer,
			Supplier:           supplier,
			Language:           language,
			MemberNumber:       memberNumber,
		}
		if street != "" || city != "" || postcode != "" || country != "" {
			req.Address = api.Address{
				StreetAddress: street,
				City:          city,
				PostCode:      postcode,
				Country:       country,
			}
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointContacts, slug), req)
		if err != nil {
			return fmt.Errorf("creating contact: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing contact ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Contact created (ID: %d)", id))
		return nil
	},
}

var contactsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a contact by ID",
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

		var contact api.Contact
		endpoint := fmt.Sprintf(api.EndpointContact, slug, id)
		if _, err := client.Get(endpoint, &contact); err != nil {
			return fmt.Errorf("fetching contact: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(contact)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", contact.ContactId))
		table.AddRow("Name", contact.Name)
		table.AddRow("Email", contact.Email)
		table.AddRow("Phone", contact.PhoneNumber)
		table.AddRow("Org Number", contact.OrganizationNumber)
		table.AddRow("Customer", BoolToYesNo(contact.Customer))
		table.AddRow("Supplier", BoolToYesNo(contact.Supplier))
		table.AddRow("Language", contact.Language)
		table.AddRow("Inactive", BoolToYesNo(contact.Inactive))
		table.Print()
		return nil
	},
}

var contactsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a contact",
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

		var existing api.Contact
		endpoint := fmt.Sprintf(api.EndpointContact, slug, id)
		if _, err := client.Get(endpoint, &existing); err != nil {
			return fmt.Errorf("fetching contact for update: %w", err)
		}

		req := api.ContactRequest{
			Name:               existing.Name,
			Email:              existing.Email,
			PhoneNumber:        existing.PhoneNumber,
			OrganizationNumber: existing.OrganizationNumber,
			Customer:           existing.Customer,
			Supplier:           existing.Supplier,
			Language:           existing.Language,
			MemberNumber:       existing.MemberNumber,
			Address:            existing.Address,
			Inactive:           existing.Inactive,
		}

		if cmd.Flags().Changed("name") {
			req.Name, _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("email") {
			req.Email, _ = cmd.Flags().GetString("email")
		}
		if cmd.Flags().Changed("phone") {
			req.PhoneNumber, _ = cmd.Flags().GetString("phone")
		}
		if cmd.Flags().Changed("org-number") {
			req.OrganizationNumber, _ = cmd.Flags().GetString("org-number")
		}
		if cmd.Flags().Changed("customer") {
			req.Customer, _ = cmd.Flags().GetBool("customer")
		}
		if cmd.Flags().Changed("supplier") {
			req.Supplier, _ = cmd.Flags().GetBool("supplier")
		}
		if cmd.Flags().Changed("language") {
			req.Language, _ = cmd.Flags().GetString("language")
		}
		if cmd.Flags().Changed("member-number") {
			req.MemberNumber, _ = cmd.Flags().GetInt64("member-number")
		}

		_, err = client.Put(endpoint, req)
		if err != nil {
			return fmt.Errorf("updating contact: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Contact %d updated", id))
		return nil
	},
}

var contactsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a contact",
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

		endpoint := fmt.Sprintf(api.EndpointContact, slug, id)
		if err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting contact: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Contact %d deleted (or deactivated if it had associated transactions)", id))
		return nil
	},
}

func init() {
	contactsCreateCmd.Flags().String("name", "", "Contact name (required)")
	contactsCreateCmd.Flags().String("email", "", "Email address")
	contactsCreateCmd.Flags().String("phone", "", "Phone number")
	contactsCreateCmd.Flags().String("org-number", "", "Organization number")
	contactsCreateCmd.Flags().Bool("customer", false, "Mark as customer")
	contactsCreateCmd.Flags().Bool("supplier", false, "Mark as supplier")
	contactsCreateCmd.Flags().String("language", "", "Language code (e.g. 'norwegian')")
	contactsCreateCmd.Flags().Int64("member-number", 0, "Member number")
	contactsCreateCmd.Flags().String("street", "", "Street address")
	contactsCreateCmd.Flags().String("city", "", "City")
	contactsCreateCmd.Flags().String("postcode", "", "Post code")
	contactsCreateCmd.Flags().String("country", "", "Country code (e.g. 'NOR')")

	contactsUpdateCmd.Flags().String("name", "", "Contact name")
	contactsUpdateCmd.Flags().String("email", "", "Email address")
	contactsUpdateCmd.Flags().String("phone", "", "Phone number")
	contactsUpdateCmd.Flags().String("org-number", "", "Organization number")
	contactsUpdateCmd.Flags().Bool("customer", false, "Mark as customer")
	contactsUpdateCmd.Flags().Bool("supplier", false, "Mark as supplier")
	contactsUpdateCmd.Flags().String("language", "", "Language code")
	contactsUpdateCmd.Flags().Int64("member-number", 0, "Member number")

	contactsCmd.AddCommand(contactsListCmd)
	contactsCmd.AddCommand(contactsCreateCmd)
	contactsCmd.AddCommand(contactsGetCmd)
	contactsCmd.AddCommand(contactsUpdateCmd)
	contactsCmd.AddCommand(contactsDeleteCmd)
	rootCmd.AddCommand(contactsCmd)
}
