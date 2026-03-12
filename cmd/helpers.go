package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jakoblind/fiken-cli/api"
)

// ValidateFile checks that filePath exists and has an allowed extension
// (.pdf, .png, .jpg, .jpeg, .gif).
func ValidateFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	allowed := map[string]bool{".pdf": true, ".png": true, ".jpg": true, ".jpeg": true, ".gif": true}
	if !allowed[ext] {
		return fmt.Errorf("unsupported file extension %q: must be .pdf, .png, .jpg, .jpeg, or .gif", ext)
	}
	return nil
}

// UploadAttachment builds a multipart form-data request with optional string
// fields and a file, then POSTs it to the given endpoint.
func UploadAttachment(client *api.Client, endpoint string, filePath string, fields map[string]string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range fields {
		if err := writer.WriteField(key, val); err != nil {
			return fmt.Errorf("writing field %q: %w", key, err)
		}
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

	if err := writer.Close(); err != nil {
		return fmt.Errorf("closing multipart writer: %w", err)
	}

	_, err = client.PostMultipart(endpoint, body, writer.FormDataContentType())
	return err
}

// ParseAmountToCents converts a decimal string like "1000.50" to int64 cents (100050).
// Uses integer arithmetic only to avoid float64 precision issues.
func ParseAmountToCents(s string) (int64, error) {
	s = strings.TrimSpace(s)
	parts := strings.SplitN(s, ".", 2)
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || whole < 0 {
		return 0, fmt.Errorf("invalid amount %q: must be a non-negative number like 100.00", s)
	}
	cents := whole * 100
	if len(parts) == 2 {
		dec := parts[1]
		switch len(dec) {
		case 1:
			d, _ := strconv.ParseInt(dec, 10, 64)
			cents += d * 10
		default:
			d, _ := strconv.ParseInt(dec[:2], 10, 64)
			cents += d
		}
	}
	return cents, nil
}

// BoolToYesNo returns "Yes" for true, "No" for false.
func BoolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// FetchAllPages fetches up to maxPages pages from a paginated endpoint,
// appending results to a slice. pageSize is capped at 100.
func FetchAllPages[T any](client *api.Client, endpoint string, params url.Values, pageSize, maxPages int) ([]T, error) {
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
	}
	if maxPages <= 0 {
		maxPages = 4
	}
	if params == nil {
		params = url.Values{}
	}
	params.Set("pageSize", strconv.Itoa(pageSize))

	var all []T
	for page := 0; page < maxPages; page++ {
		params.Set("page", strconv.Itoa(page))
		var pageItems []T
		pagination, err := client.GetWithParams(endpoint, params, &pageItems)
		if err != nil {
			return nil, err
		}
		all = append(all, pageItems...)
		if pagination == nil || page+1 >= pagination.PageCount || len(pageItems) == 0 {
			break
		}
	}
	return all, nil
}
