package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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

var inboxUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a document to the inbox",
	Long:  "Upload a PDF or image file to the company's inbox.",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		allowed := map[string]bool{".pdf": true, ".png": true, ".jpg": true, ".jpeg": true, ".gif": true}
		if !allowed[ext] {
			return fmt.Errorf("unsupported file extension %q: must be .pdf, .png, .jpg, .jpeg, or .gif", ext)
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

		endpoint := fmt.Sprintf(api.EndpointInbox, slug)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		writer.WriteField("name", name)
		writer.WriteField("filename", filepath.Base(filePath))
		if description != "" {
			writer.WriteField("description", description)
		}

		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			return fmt.Errorf("creating form file: %w", err)
		}

		if _, err := io.Copy(part, f); err != nil {
			return fmt.Errorf("writing file to multipart: %w", err)
		}

		// CRITICAL: Close writer BEFORE reading body
		writer.Close()

		_, err = client.PostMultipart(endpoint, body, writer.FormDataContentType())
		if err != nil {
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

	inboxCmd.AddCommand(inboxUploadCmd)
	rootCmd.AddCommand(inboxCmd)
}
