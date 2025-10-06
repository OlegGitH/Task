// Package customerimporter reads from a CSV file and returns a sorted (data
// structure of your choice) of email domains along with the number of customers
// with e-mail addresses for each domain. This should be able to be ran from the
// CLI and output the sorted domains to the terminal or to a file. Any errors
// should be logged (or handled). Performance matters (this is only ~3k lines,
// but could be 1m lines or run on a small machine).
package customerimporter

import (
	"bufio"
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type DomainData struct {
	Domain           string
	CustomerQuantity uint64
}

type CustomerImporter struct {
	path string
}

// NewCustomerImporter returns a new CustomerImporter that reads from file at specified path.
func NewCustomerImporter(filePath string) *CustomerImporter {
	return &CustomerImporter{
		path: filePath,
	}
}

// ImportDomainData reads and returns sorted customer domain data from CSV file.
// Uses optimized streaming algorithm with buffered IO for best performance across all file sizes
func (ci CustomerImporter) ImportDomainData() ([]DomainData, error) {
	// Validate input path
	if ci.path == "" {
		slog.Error("import validation failed", "reason", "file path is empty")
		return nil, fmt.Errorf("file path cannot be empty")
	}

	file, err := os.Open(ci.path)
	if err != nil {
		slog.Error("failed to open input file", "path", ci.path, "error", err)
		return nil, fmt.Errorf("failed to open file %s: %w", ci.path, err)
	}
	defer file.Close()

	// Use buffered scanner for optimal IO performance
	scanner := bufio.NewScanner(file)

	// reduce rehashing
	domainCounts := make(map[string]uint64, 1000)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			slog.Error("failed to read header", "path", ci.path, "error", err)
			return nil, fmt.Errorf("failed to read header from %s: %w", ci.path, err)
		}
		slog.Error("input file is empty", "path", ci.path)
		return nil, fmt.Errorf("file %s appears to be empty", ci.path)
	}

	// Process lines with optimized string operations
	linesProcessed := 0
	validEmails := 0
	invalidEmails := 0

	for scanner.Scan() {
		linesProcessed++
		line := scanner.Text()

		// Fast CSV parsing - split by comma
		fields := strings.Split(line, ",")
		if len(fields) < 3 {
			invalidEmails++
			continue
		}

		email := fields[2]

		// Fast domain extraction using IndexByte for optimal performance
		atIndex := strings.IndexByte(email, '@')
		if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
			invalidEmails++
			continue
		}

		domain := email[atIndex+1:]
		if domain == "" {
			invalidEmails++
			continue
		}

		domainCounts[domain]++
		validEmails++
	}

	if err := scanner.Err(); err != nil {
		slog.Error("scanner error during file processing", "path", ci.path, "error", err)
		return nil, fmt.Errorf("error reading file %s: %w", ci.path, err)
	}

	// Convert to slice with pre-allocated capacity
	domainData := make([]DomainData, 0, len(domainCounts))
	for domain, count := range domainCounts {
		domainData = append(domainData, DomainData{
			Domain:           domain,
			CustomerQuantity: count,
		})
	}

	// Optimized sorting using slices.SortFunc
	slices.SortFunc(domainData, func(l, r DomainData) int {
		return cmp.Compare(l.Domain, r.Domain)
	})
	return domainData, nil
}
