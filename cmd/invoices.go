package cmd

import (
	"fmt"
	"path/filepath"

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

		if err := ValidateFile(filePath); err != nil {
			return err
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
		}
		endpoint := fmt.Sprintf(api.EndpointInvoiceAttachments, slug, invoiceID)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
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
