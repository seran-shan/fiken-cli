package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var invoicesCmd = &cobra.Command{
	Use:   "invoices",
	Short: "Manage invoices",
}

var invoicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List invoices",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not yet implemented.")
		return nil
	},
}

var invoicesAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a document to an invoice",
	RunE: func(cmd *cobra.Command, args []string) error {
		invoiceID, _ := cmd.Flags().GetInt64("id")
		filePath, _ := cmd.Flags().GetString("file")

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		allowed := map[string]bool{".pdf": true, ".png": true, ".jpg": true, ".jpeg": true, ".gif": true}
		if !allowed[ext] {
			return fmt.Errorf("unsupported file extension %q: must be .pdf, .png, .jpg, .jpeg, or .gif", ext)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		slug, err := resolveCompany(client)
		if err != nil {
			return err
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		writer.WriteField("filename", filepath.Base(filePath))

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

		endpoint := fmt.Sprintf(api.EndpointInvoiceAttachments, slug, invoiceID)
		_, err = client.PostMultipart(endpoint, body, writer.FormDataContentType())
		if err != nil {
			return fmt.Errorf("attaching to invoice: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to invoice %d", invoiceID))
		return nil
	},
}

func init() {
	invoicesAttachCmd.Flags().Int64("id", 0, "Invoice ID to attach to (required)")
	invoicesAttachCmd.MarkFlagRequired("id")
	invoicesAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	invoicesAttachCmd.MarkFlagRequired("file")

	invoicesCmd.AddCommand(invoicesListCmd)
	invoicesCmd.AddCommand(invoicesAttachCmd)
	rootCmd.AddCommand(invoicesCmd)
}
