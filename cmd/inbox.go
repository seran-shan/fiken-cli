package cmd

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var inboxStatus string

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "List EHF inbox documents",
	Long:  "List documents in the EHF (electronic invoice) inbox.",
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
		if inboxStatus != "" {
			params.Set("status", inboxStatus)
		}
		params.Set("pageSize", "25")

		endpoint := fmt.Sprintf(api.EndpointInbox, slug)

		var documents []api.InboxDocument
		page := 0
		for {
			params.Set("page", fmt.Sprintf("%d", page))
			var pageDocs []api.InboxDocument
			pagination, err := client.GetWithParams(endpoint, params, &pageDocs)
			if err != nil {
				return fmt.Errorf("fetching inbox: %w", err)
			}
			documents = append(documents, pageDocs...)

			if pagination == nil || page+1 >= pagination.PageCount {
				break
			}
			page++
		}

		if jsonOutput {
			return output.PrintJSON(documents)
		}

		if len(documents) == 0 {
			output.PrintInfo("Inbox is empty.")
			return nil
		}

		table := output.NewTable("ID", "NAME", "FILENAME", "STATUS", "DATE")
		for _, d := range documents {
			table.AddRow(
				fmt.Sprintf("%d", d.DocumentId),
				d.Name,
				d.Filename,
				d.Status,
				d.CreatedDate.Format("2006-01-02"),
			)
		}
		table.Print()

		fmt.Printf("\n%d documents\n", len(documents))
		return nil
	},
}

var inboxGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get an inbox document by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid inbox document ID %q: %w", args[0], err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		var doc api.InboxDocument
		_, err = client.Get(fmt.Sprintf(api.EndpointInboxDocument, slug, id), &doc)
		if err != nil {
			return fmt.Errorf("fetching inbox document: %w", err)
		}

		if jsonOutput {
			return output.PrintJSON(doc)
		}

		table := output.NewTable("FIELD", "VALUE")
		table.AddRow("ID", fmt.Sprintf("%d", doc.DocumentId))
		table.AddRow("Name", doc.Name)
		table.AddRow("Description", doc.Description)
		table.AddRow("Filename", doc.Filename)
		table.AddRow("Status", doc.Status)
		table.AddRow("Created Date", doc.CreatedDate.Format("2006-01-02"))
		table.Print()

		return nil
	},
}

var inboxUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a document to the inbox",
	Long:  "Upload a PDF or image file to the company's inbox.",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if err := ValidateFile(filePath); err != nil {
			return err
		}

		if name == "" {
			name = filepath.Base(filePath)
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
			"filename": filepath.Base(filePath),
			"name":     name,
		}
		if description != "" {
			fields["description"] = description
		}
		endpoint := fmt.Sprintf(api.EndpointInbox, slug)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("uploading to inbox: %w", err)
		}

		output.PrintSuccess("Document uploaded to inbox: " + name)
		return nil
	},
}

func init() {
	inboxCmd.Flags().StringVar(&inboxStatus, "status", "", "Filter by status (pending, processed)")

	inboxUploadCmd.Flags().String("file", "", "Path to the file to upload (required)")
	inboxUploadCmd.MarkFlagRequired("file")
	inboxUploadCmd.Flags().String("name", "", "Document name (defaults to filename)")
	inboxUploadCmd.Flags().String("description", "", "Document description")

	inboxCmd.AddCommand(inboxGetCmd)
	inboxCmd.AddCommand(inboxUploadCmd)
	rootCmd.AddCommand(inboxCmd)
}
