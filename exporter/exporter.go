package exporter

import (
	"encoding/csv"
	"fmt"
	"importer/customerimporter"
	"io"
	"log/slog"
	"os"
	"strconv"
)

type CustomerExporter struct {
	outputPath string
}

// NewCustomerExporter returns a new CustomerExporter that writes customer domain data to specified file.
func NewCustomerExporter(outputPath string) *CustomerExporter {
	return &CustomerExporter{
		outputPath: outputPath,
	}
}

// ExportData writes sorted customer domain data to a CSV file.
func (ex CustomerExporter) ExportData(data []customerimporter.DomainData) error {
	if data == nil {
		slog.Error("export data validation failed", "reason", "data is nil")
		return fmt.Errorf("error provided data is empty (nil)")
	}

	outputFile, err := os.Create(ex.outputPath)
	if err != nil {
		slog.Error("failed to create output file", "path", ex.outputPath, "error", err)
		return fmt.Errorf("error creating new file for saving: %v", err)
	}
	defer outputFile.Close()

	err = exportCsv(data, outputFile)
	if err != nil {
		slog.Error("CSV export failed", "path", ex.outputPath, "error", err)
		return err
	}

	return nil
}

// exportCsv writes customer domain data to a CSV file.
func exportCsv(data []customerimporter.DomainData, output io.Writer) error {
	headers := []string{"domain", "number_of_customers"}
	csvWriter := csv.NewWriter(output)
	defer func() {
		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			slog.Error("CSV writer flush failed", "error", err)
		}
	}()

	if err := csvWriter.Write(headers); err != nil {
		slog.Error("failed to write CSV headers", "error", err)
		return err
	}

	for _, v := range data {
		pair := []string{v.Domain, strconv.FormatUint(v.CustomerQuantity, 10)}
		if err := csvWriter.Write(pair); err != nil {
			slog.Error("failed to write CSV record", "domain", v.Domain, "error", err)
			return err
		}
	}

	return nil
}
