package cmd

import (
	"fmt"
	"strconv"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var contactsPersonsCmd = &cobra.Command{
	Use:   "persons",
	Short: "Manage contact persons",
}

var contactsPersonsListCmd = &cobra.Command{
	Use:   "list <contact-id>",
	Short: "List contact persons",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contactID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid contact ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var persons []api.ContactPerson
		endpoint := fmt.Sprintf(api.EndpointContactPersons, slug, contactID)
		if _, err := client.Get(endpoint, &persons); err != nil {
			return fmt.Errorf("fetching contact persons: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(persons)
		}

		if len(persons) == 0 {
			output.PrintInfo("No contact persons found.")
			return nil
		}

		table := output.NewTable("PERSON ID", "NAME", "EMAIL", "PHONE NUMBER")
		for _, p := range persons {
			table.AddRow(
				fmt.Sprintf("%d", p.ContactPersonId),
				p.Name,
				p.Email,
				p.PhoneNumber,
			)
		}
		table.Print()
		return nil
	},
}

var contactsPersonsCreateCmd = &cobra.Command{
	Use:   "create <contact-id>",
	Short: "Create a contact person",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contactID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid contact ID %q: %w", args[0], err)
		}

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("missing required flag: --name")
		}
		email, _ := cmd.Flags().GetString("email")
		phoneNumber, _ := cmd.Flags().GetString("phone-number")

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		req := api.ContactPersonRequest{
			Name:        name,
			Email:       email,
			PhoneNumber: phoneNumber,
		}

		endpoint := fmt.Sprintf(api.EndpointContactPersons, slug, contactID)
		locationURL, err := client.PostCreate(endpoint, req)
		if err != nil {
			return fmt.Errorf("creating contact person: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing contact person ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Contact person created (ID: %d)", id))
		return nil
	},
}

var contactsPersonsGetCmd = &cobra.Command{
	Use:   "get <contact-id>",
	Short: "Get a contact person by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contactID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid contact ID %q: %w", args[0], err)
		}

		personID, _ := cmd.Flags().GetInt64("person-id")

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var person api.ContactPerson
		endpoint := fmt.Sprintf(api.EndpointContactPerson, slug, contactID, personID)
		if _, err := client.Get(endpoint, &person); err != nil {
			return fmt.Errorf("fetching contact person: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(person)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("Person ID", fmt.Sprintf("%d", person.ContactPersonId))
		table.AddRow("Name", person.Name)
		table.AddRow("Email", person.Email)
		table.AddRow("Phone Number", person.PhoneNumber)
		table.Print()
		return nil
	},
}

var contactsPersonsUpdateCmd = &cobra.Command{
	Use:   "update <contact-id>",
	Short: "Update a contact person",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contactID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid contact ID %q: %w", args[0], err)
		}

		personID, _ := cmd.Flags().GetInt64("person-id")

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var existing api.ContactPerson
		endpoint := fmt.Sprintf(api.EndpointContactPerson, slug, contactID, personID)
		if _, err := client.Get(endpoint, &existing); err != nil {
			return fmt.Errorf("fetching contact person for update: %w", err)
		}

		req := api.ContactPersonRequest{
			Name:        existing.Name,
			Email:       existing.Email,
			PhoneNumber: existing.PhoneNumber,
		}

		if cmd.Flags().Changed("name") {
			req.Name, _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("email") {
			req.Email, _ = cmd.Flags().GetString("email")
		}
		if cmd.Flags().Changed("phone-number") {
			req.PhoneNumber, _ = cmd.Flags().GetString("phone-number")
		}

		if _, err := client.Put(endpoint, req); err != nil {
			return fmt.Errorf("updating contact person: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Contact person %d updated", personID))
		return nil
	},
}

var contactsPersonsDeleteCmd = &cobra.Command{
	Use:   "delete <contact-id>",
	Short: "Delete a contact person",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contactID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid contact ID %q: %w", args[0], err)
		}

		personID, _ := cmd.Flags().GetInt64("person-id")

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		endpoint := fmt.Sprintf(api.EndpointContactPerson, slug, contactID, personID)
		if _, err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting contact person: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Contact person %d deleted", personID))
		return nil
	},
}

func init() {
	contactsPersonsCreateCmd.Flags().String("name", "", "Person name (required)")
	contactsPersonsCreateCmd.Flags().String("email", "", "Email address")
	contactsPersonsCreateCmd.Flags().String("phone-number", "", "Phone number")

	contactsPersonsGetCmd.Flags().Int64("person-id", 0, "Contact person ID (required)")
	contactsPersonsGetCmd.MarkFlagRequired("person-id")

	contactsPersonsUpdateCmd.Flags().Int64("person-id", 0, "Contact person ID (required)")
	contactsPersonsUpdateCmd.MarkFlagRequired("person-id")
	contactsPersonsUpdateCmd.Flags().String("name", "", "Person name")
	contactsPersonsUpdateCmd.Flags().String("email", "", "Email address")
	contactsPersonsUpdateCmd.Flags().String("phone-number", "", "Phone number")

	contactsPersonsDeleteCmd.Flags().Int64("person-id", 0, "Contact person ID (required)")
	contactsPersonsDeleteCmd.MarkFlagRequired("person-id")

	contactsPersonsCmd.AddCommand(contactsPersonsListCmd)
	contactsPersonsCmd.AddCommand(contactsPersonsCreateCmd)
	contactsPersonsCmd.AddCommand(contactsPersonsGetCmd)
	contactsPersonsCmd.AddCommand(contactsPersonsUpdateCmd)
	contactsPersonsCmd.AddCommand(contactsPersonsDeleteCmd)
	contactsCmd.AddCommand(contactsPersonsCmd)
}
