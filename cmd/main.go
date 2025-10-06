package main

import (
	"flag"
	"fmt"
	"importer/customerimporter"
	"importer/exporter"
	"log/slog"
)

// ProcessCustomerData handles the main application logic
func ProcessCustomerData(inputPath, outputPath string) error {
	// Import customer data
	importer := customerimporter.NewCustomerImporter(inputPath)
	data, err := importer.ImportDomainData()
	if err != nil {
		return fmt.Errorf("failed to import customer data from %s: %w", inputPath, err)
	}

	if outputPath == "" {
		fmt.Println("domain,number_of_customers")
		for _, v := range data {
			fmt.Printf("%s,%v\n", v.Domain, v.CustomerQuantity)
		}
	} else {
		exporter := exporter.NewCustomerExporter(outputPath)
		if err := exporter.ExportData(data); err != nil {
			return fmt.Errorf("failed to export data to %s: %w", outputPath, err)
		}
	}

	return nil
}

func main() {
	inputPath := flag.String("path", "./customers.csv", "Path to the input CSV file")
	outputPath := flag.String("out", "", "Output file path (if empty, prints to terminal)")
	flag.Parse()

	// Process the data
	if err := ProcessCustomerData(*inputPath, *outputPath); err != nil {
		slog.Error("application error", "error", err)
	}
}
