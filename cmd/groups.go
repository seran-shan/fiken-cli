package cmd

import (
	"fmt"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List contact groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var groups []api.Group
		_, err = client.Get(fmt.Sprintf(api.EndpointGroups, slug), &groups)
		if err != nil {
			return fmt.Errorf("fetching groups: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(groups)
		}

		if len(groups) == 0 {
			output.PrintInfo("No groups found.")
			return nil
		}

		table := output.NewTable("NAME")
		for _, g := range groups {
			table.AddRow(g.Name)
		}
		table.Print()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(groupsCmd)
}
