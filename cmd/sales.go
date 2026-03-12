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
		writer.WriteField("attachToPayment", fmt.Sprintf("%v", attachToPayment))
		writer.WriteField("attachToSale", fmt.Sprintf("%v", attachToSale))

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

		endpoint := fmt.Sprintf(api.EndpointSaleAttachments, slug, saleID)
		_, err = client.PostMultipart(endpoint, body, writer.FormDataContentType())
		if err != nil {
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
