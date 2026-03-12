package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/jakoblind/fiken-cli/api"
	"github.com/jakoblind/fiken-cli/output"
	"github.com/spf13/cobra"
)

var salesCmd = &cobra.Command{
	Use:   "sales",
	Short: "Manage sales",
}

var salesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sales",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not yet implemented.")
		return nil
	},
}

var salesAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a receipt/document to a sale",
	RunE: func(cmd *cobra.Command, args []string) error {
		saleID, _ := cmd.Flags().GetInt64("id")
		filePath, _ := cmd.Flags().GetString("file")
		attachToPayment, _ := cmd.Flags().GetBool("attach-to-payment")
		attachToSale, _ := cmd.Flags().GetBool("attach-to-sale")

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
			"filename":        filepath.Base(filePath),
			"attachToPayment": fmt.Sprintf("%v", attachToPayment),
			"attachToSale":    fmt.Sprintf("%v", attachToSale),
		}
		endpoint := fmt.Sprintf(api.EndpointSaleAttachments, slug, saleID)
		if err := UploadAttachment(client, endpoint, filePath, fields); err != nil {
			return fmt.Errorf("attaching to sale: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Attachment added to sale %d", saleID))
		return nil
	},
}

func init() {
	salesAttachCmd.Flags().Int64("id", 0, "Sale ID to attach to (required)")
	salesAttachCmd.MarkFlagRequired("id")
	salesAttachCmd.Flags().String("file", "", "Path to the file to attach (required)")
	salesAttachCmd.MarkFlagRequired("file")
	salesAttachCmd.Flags().Bool("attach-to-payment", true, "Whether this documents the payment")
	salesAttachCmd.Flags().Bool("attach-to-sale", true, "Whether this documents the sale itself")

	salesCmd.AddCommand(salesListCmd)
	salesCmd.AddCommand(salesAttachCmd)
	rootCmd.AddCommand(salesCmd)
}
