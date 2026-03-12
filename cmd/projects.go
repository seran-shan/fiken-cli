package cmd

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  "List, create, get, update, and delete projects.",
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var params url.Values
		if cmd.Flags().Changed("completed") {
			completed, _ := cmd.Flags().GetBool("completed")
			params = url.Values{}
			params.Set("completed", strconv.FormatBool(completed))
		}

		endpoint := fmt.Sprintf(api.EndpointProjects, slug)
		projects, err := FetchAllPages[api.Project](client, endpoint, params, 100, 4)
		if err != nil {
			return fmt.Errorf("fetching projects: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(projects)
		}

		if len(projects) == 0 {
			output.PrintInfo("No projects found.")
			return nil
		}

		table := output.NewTable("ID", "NUMBER", "NAME", "START DATE", "END DATE", "CONTACT", "COMPLETED")
		for _, p := range projects {
			table.AddRow(
				fmt.Sprintf("%d", p.ProjectId),
				p.Number,
				p.Name,
				p.StartDate,
				p.EndDate,
				p.Contact.Name,
				BoolToYesNo(p.Completed),
			)
		}
		table.Print()
		fmt.Printf("\n%d projects\n", len(projects))
		return nil
	},
}

var projectsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("missing required flag: --name")
		}

		number, _ := cmd.Flags().GetString("number")
		description, _ := cmd.Flags().GetString("description")
		startDate, _ := cmd.Flags().GetString("start-date")
		endDate, _ := cmd.Flags().GetString("end-date")
		contactId, _ := cmd.Flags().GetInt64("contact-id")
		completed, _ := cmd.Flags().GetBool("completed")

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		req := api.ProjectRequest{
			Name:        name,
			Number:      number,
			Description: description,
			StartDate:   startDate,
			EndDate:     endDate,
			ContactId:   contactId,
			Completed:   completed,
		}

		locationURL, err := client.PostCreate(fmt.Sprintf(api.EndpointProjects, slug), req)
		if err != nil {
			return fmt.Errorf("creating project: %w", err)
		}

		id, err := api.ParseIDFromLocation(locationURL)
		if err != nil {
			return fmt.Errorf("parsing project ID: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Project created (ID: %d)", id))
		return nil
	},
}

var projectsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a project by ID",
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

		var project api.Project
		endpoint := fmt.Sprintf(api.EndpointProject, slug, id)
		if _, err := client.Get(endpoint, &project); err != nil {
			return fmt.Errorf("fetching project: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(project)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", project.ProjectId))
		table.AddRow("Number", project.Number)
		table.AddRow("Name", project.Name)
		table.AddRow("Description", project.Description)
		table.AddRow("Start Date", project.StartDate)
		table.AddRow("End Date", project.EndDate)
		table.AddRow("Contact", project.Contact.Name)
		table.AddRow("Completed", BoolToYesNo(project.Completed))
		table.Print()
		return nil
	},
}

var projectsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a project",
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

		var existing api.Project
		endpoint := fmt.Sprintf(api.EndpointProject, slug, id)
		if _, err := client.Get(endpoint, &existing); err != nil {
			return fmt.Errorf("fetching project for update: %w", err)
		}

		req := api.ProjectRequest{
			Name:        existing.Name,
			Number:      existing.Number,
			Description: existing.Description,
			StartDate:   existing.StartDate,
			EndDate:     existing.EndDate,
			ContactId:   existing.Contact.ContactId,
			Completed:   existing.Completed,
		}

		if cmd.Flags().Changed("name") {
			req.Name, _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("number") {
			req.Number, _ = cmd.Flags().GetString("number")
		}
		if cmd.Flags().Changed("description") {
			req.Description, _ = cmd.Flags().GetString("description")
		}
		if cmd.Flags().Changed("start-date") {
			req.StartDate, _ = cmd.Flags().GetString("start-date")
		}
		if cmd.Flags().Changed("end-date") {
			req.EndDate, _ = cmd.Flags().GetString("end-date")
		}
		if cmd.Flags().Changed("contact-id") {
			req.ContactId, _ = cmd.Flags().GetInt64("contact-id")
		}
		if cmd.Flags().Changed("completed") {
			req.Completed, _ = cmd.Flags().GetBool("completed")
		}

		if err := client.Patch(endpoint, req, nil); err != nil {
			return fmt.Errorf("updating project: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Project %d updated", id))
		return nil
	},
}

var projectsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a project",
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

		endpoint := fmt.Sprintf(api.EndpointProject, slug, id)
		if _, err := client.Delete(endpoint); err != nil {
			return fmt.Errorf("deleting project: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Project %d deleted", id))
		return nil
	},
}

func init() {
	projectsListCmd.Flags().Bool("completed", false, "Filter by completed status")

	projectsCreateCmd.Flags().String("name", "", "Project name (required)")
	projectsCreateCmd.Flags().String("number", "", "Project number")
	projectsCreateCmd.Flags().String("description", "", "Project description")
	projectsCreateCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	projectsCreateCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	projectsCreateCmd.Flags().Int64("contact-id", 0, "Associated contact ID")
	projectsCreateCmd.Flags().Bool("completed", false, "Mark project as completed")

	projectsUpdateCmd.Flags().String("name", "", "Project name")
	projectsUpdateCmd.Flags().String("number", "", "Project number")
	projectsUpdateCmd.Flags().String("description", "", "Project description")
	projectsUpdateCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	projectsUpdateCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	projectsUpdateCmd.Flags().Int64("contact-id", 0, "Associated contact ID")
	projectsUpdateCmd.Flags().Bool("completed", false, "Mark project as completed")

	projectsCmd.AddCommand(projectsListCmd)
	projectsCmd.AddCommand(projectsCreateCmd)
	projectsCmd.AddCommand(projectsGetCmd)
	projectsCmd.AddCommand(projectsUpdateCmd)
	projectsCmd.AddCommand(projectsDeleteCmd)
	rootCmd.AddCommand(projectsCmd)
}
