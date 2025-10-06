# Provides convenient commands for building, testing, and running the application

# Variables
APP_NAME = customer-importer
MAIN_FILE = cmd/main.go
BUILD_DIR = build
COVERAGE_DIR = coverage

# Default target
.PHONY: all
all: build test

# Build the application
.PHONY: build
build:
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME).exe $(MAIN_FILE)

# Run the application with default settings
.PHONY: run
run: build
	@echo "Running $(APP_NAME) with default settings..."
	@$(BUILD_DIR)/$(APP_NAME).exe -path customers.csv

# Run with output to file
.PHONY: run-output
run-output: build
	@echo "Running $(APP_NAME) with file output..."
	@$(BUILD_DIR)/$(APP_NAME).exe -path customerimporter/test_data.csv -out output.csv

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with verbose output
.PHONY: test-verbose
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@if not exist $(COVERAGE_DIR) mkdir $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

# Run tests with coverage and show summary
.PHONY: test-coverage-summary
test-coverage-summary:
	@echo "Running tests with coverage summary..."
	@go test -cover ./...

# Run benchmarks
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. ./...

# Run benchmarks with memory profiling
.PHONY: benchmark-mem
benchmark-mem:
	@echo "Running benchmarks with memory profiling..."
	@go test -bench=. -benchmem ./...

# Run specific benchmark
.PHONY: benchmark-import
benchmark-import:
	@echo "Running import benchmarks..."
	@go test -bench=BenchmarkImportDomainData -benchmem ./customerimporter

# Run specific benchmark
.PHONY: benchmark-export
benchmark-export:
	@echo "Running export benchmarks..."
	@go test -bench=BenchmarkImportDomainData -benchmem ./exporter

# Run performance tests
.PHONY: perf-test
perf-test:
	@echo "Running performance tests..."
	@go test -v -run=TestPerformanceMetrics

# Run scalability tests
.PHONY: scalability-test
scalability-test:
	@echo "Running scalability tests..."
	@go test -v -run=TestScalability