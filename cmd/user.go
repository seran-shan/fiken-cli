package cmd

import (
	"fmt"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Show authenticated user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		var user api.UserInfo
		if _, err := client.Get(api.EndpointUser, &user); err != nil {
			return fmt.Errorf("fetching user info: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(user)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("Name", user.Name)
		table.AddRow("Email", user.Email)
		table.Print()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
}
