package customerimporter

import (
	"testing"
)

func TestImportDomainData(t *testing.T) {
	path := "./test_data.csv"
	importer := NewCustomerImporter(path)

	data, err := importer.ImportDomainData()
	if err != nil {
		t.Fatalf("ImportDomainData failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected domain data, got empty result")
	}

	for _, domain := range data {
		if domain.Domain == "" {
			t.Error("Found empty domain name")
		}
		if domain.CustomerQuantity == 0 {
			t.Error("Found domain with zero customer count")
		}
	}
}

func TestImportDomainDataSorting(t *testing.T) {
	expectedDomains := []string{
		"360.cn", "acquirethisname.com", "blogtalkradio.com", "chicagotribune.com",
		"cnet.com", "cyberchimps.com", "github.io", "hubpages.com", "rediff.com", "statcounter.com",
	}

	path := "./test_data.csv"
	importer := NewCustomerImporter(path)
	data, err := importer.ImportDomainData()
	if err != nil {
		t.Fatalf("ImportDomainData failed: %v", err)
	}

	if len(data) != len(expectedDomains) {
		t.Fatalf("Expected %d domains, got %d", len(expectedDomains), len(data))
	}

	for i, expected := range expectedDomains {
		if data[i].Domain != expected {
			t.Errorf("Domain mismatch at index %d: got %s, want %s", i, data[i].Domain, expected)
		}
	}
}

func TestImportDomainDataInvalidPath(t *testing.T) {
	path := "./nonexistent_file.csv"
	importer := NewCustomerImporter(path)

	_, err := importer.ImportDomainData()
	if err == nil {
		t.Error("Expected error for invalid file path, got nil")
	}
}

func TestImportDomainDataEmptyPath(t *testing.T) {
	path := ""
	importer := NewCustomerImporter(path)

	_, err := importer.ImportDomainData()
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
}

func TestImportDomainDataInvalidData(t *testing.T) {
	path := "./test_invalid_data.csv"
	importer := NewCustomerImporter(path)

	data, err := importer.ImportDomainData()
	if err != nil {
		t.Fatalf("Should handle invalid data gracefully: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected some valid domains even with invalid data")
	}

	for _, domain := range data {
		if domain.Domain == "" {
			t.Error("Found empty domain in result")
		}
		if domain.CustomerQuantity == 0 {
			t.Error("Found domain with zero customer count")
		}
	}
}

func TestImportDomainDataLargeFile(t *testing.T) {
	path := "./benchmark10k.csv"
	importer := NewCustomerImporter(path)

	data, err := importer.ImportDomainData()
	if err != nil {
		t.Fatalf("ImportDomainData failed with large file: %v", err)
	}

	if len(data) < 100 {
		t.Errorf("Expected many domains from large file, got %d", len(data))
	}

	for i := 1; i < len(data); i++ {
		if data[i-1].Domain > data[i].Domain {
			t.Errorf("Data not sorted: %s > %s", data[i-1].Domain, data[i].Domain)
		}
	}
}

func TestNewCustomerImporter(t *testing.T) {
	path := "./test_data.csv"
	importer := NewCustomerImporter(path)

	if importer == nil {
		t.Error("NewCustomerImporter returned nil")
	}

	if importer.path != path {
		t.Error("Path not set correctly")
	}
}

func BenchmarkImportDomainData(b *testing.B) {
	path := "./benchmark10k.csv"
	importer := NewCustomerImporter(path)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if _, err := importer.ImportDomainData(); err != nil {
			b.Fatalf("ImportDomainData failed: %v", err)
		}
	}
}

func BenchmarkImportDomainDataSmallFile(b *testing.B) {
	path := "./test_data.csv"
	importer := NewCustomerImporter(path)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if _, err := importer.ImportDomainData(); err != nil {
			b.Fatalf("ImportDomainData failed: %v", err)
		}
	}
}
