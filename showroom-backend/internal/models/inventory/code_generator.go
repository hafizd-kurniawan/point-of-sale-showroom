package inventory

import (
	"fmt"
	"time"
)

// CodeGenerator provides auto-generation functionality for various entity codes
type CodeGenerator struct{}

// NewCodeGenerator creates a new code generator instance
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

// GenerateProductCode generates a product code in format PRD-001, PRD-002, etc.
func (cg *CodeGenerator) GenerateProductCode(lastID int) string {
	return fmt.Sprintf("PRD-%03d", lastID+1)
}

// GeneratePONumber generates a purchase order number in format PO-2025-001, PO-2025-002, etc.
func (cg *CodeGenerator) GeneratePONumber(lastID int) string {
	year := time.Now().Year()
	return fmt.Sprintf("PO-%d-%03d", year, lastID+1)
}

// GenerateReceiptNumber generates a goods receipt number in format GR-2025-001, GR-2025-002, etc.
func (cg *CodeGenerator) GenerateReceiptNumber(lastID int) string {
	year := time.Now().Year()
	return fmt.Sprintf("GR-%d-%03d", year, lastID+1)
}

// GeneratePaymentNumber generates a payment number in format PAY-2025-001, PAY-2025-002, etc.
func (cg *CodeGenerator) GeneratePaymentNumber(lastID int) string {
	year := time.Now().Year()
	return fmt.Sprintf("PAY-%d-%03d", year, lastID+1)
}

// GenerateAdjustmentNumber generates a stock adjustment number in format ADJ-2025-001, ADJ-2025-002, etc.
func (cg *CodeGenerator) GenerateAdjustmentNumber(lastID int) string {
	year := time.Now().Year()
	return fmt.Sprintf("ADJ-%d-%03d", year, lastID+1)
}

// GenerateMovementNumber generates a stock movement number in format MOV-2025-001, MOV-2025-002, etc.
func (cg *CodeGenerator) GenerateMovementNumber(lastID int) string {
	year := time.Now().Year()
	return fmt.Sprintf("MOV-%d-%03d", year, lastID+1)
}

// GetNextProductCode gets the next available product code based on database query
func (cg *CodeGenerator) GetNextProductCode(getLastIDFunc func() (int, error)) (string, error) {
	lastID, err := getLastIDFunc()
	if err != nil {
		return "", err
	}
	return cg.GenerateProductCode(lastID), nil
}

// GetNextPONumber gets the next available PO number based on database query
func (cg *CodeGenerator) GetNextPONumber(getLastIDFunc func() (int, error)) (string, error) {
	lastID, err := getLastIDFunc()
	if err != nil {
		return "", err
	}
	return cg.GeneratePONumber(lastID), nil
}

// GetNextReceiptNumber gets the next available receipt number based on database query
func (cg *CodeGenerator) GetNextReceiptNumber(getLastIDFunc func() (int, error)) (string, error) {
	lastID, err := getLastIDFunc()
	if err != nil {
		return "", err
	}
	return cg.GenerateReceiptNumber(lastID), nil
}

// GetNextPaymentNumber gets the next available payment number based on database query
func (cg *CodeGenerator) GetNextPaymentNumber(getLastIDFunc func() (int, error)) (string, error) {
	lastID, err := getLastIDFunc()
	if err != nil {
		return "", err
	}
	return cg.GeneratePaymentNumber(lastID), nil
}

// ValidateProductCode validates if a product code follows the correct format
func (cg *CodeGenerator) ValidateProductCode(code string) bool {
	// PRD-XXX format validation
	if len(code) != 7 {
		return false
	}
	if code[:4] != "PRD-" {
		return false
	}
	// Check if last 3 characters are digits
	for i := 4; i < 7; i++ {
		if code[i] < '0' || code[i] > '9' {
			return false
		}
	}
	return true
}

// ValidatePONumber validates if a PO number follows the correct format
func (cg *CodeGenerator) ValidatePONumber(number string) bool {
	// PO-YYYY-XXX format validation
	if len(number) < 10 {
		return false
	}
	if number[:3] != "PO-" {
		return false
	}
	// Additional validation can be added here
	return true
}

// ValidateReceiptNumber validates if a receipt number follows the correct format
func (cg *CodeGenerator) ValidateReceiptNumber(number string) bool {
	// GR-YYYY-XXX format validation
	if len(number) < 10 {
		return false
	}
	if number[:3] != "GR-" {
		return false
	}
	// Additional validation can be added here
	return true
}

// ValidatePaymentNumber validates if a payment number follows the correct format
func (cg *CodeGenerator) ValidatePaymentNumber(number string) bool {
	// PAY-YYYY-XXX format validation
	if len(number) < 11 {
		return false
	}
	if number[:4] != "PAY-" {
		return false
	}
	// Additional validation can be added here
	return true
}

// GetCodeInfo extracts information from a generated code
type CodeInfo struct {
	Prefix string
	Year   *int
	Number int
}

// ParseProductCode parses a product code and returns its components
func (cg *CodeGenerator) ParseProductCode(code string) (*CodeInfo, error) {
	if !cg.ValidateProductCode(code) {
		return nil, fmt.Errorf("invalid product code format: %s", code)
	}
	
	var number int
	_, err := fmt.Sscanf(code, "PRD-%03d", &number)
	if err != nil {
		return nil, err
	}
	
	return &CodeInfo{
		Prefix: "PRD",
		Number: number,
	}, nil
}

// ParsePONumber parses a PO number and returns its components
func (cg *CodeGenerator) ParsePONumber(number string) (*CodeInfo, error) {
	if !cg.ValidatePONumber(number) {
		return nil, fmt.Errorf("invalid PO number format: %s", number)
	}
	
	var year, num int
	_, err := fmt.Sscanf(number, "PO-%d-%03d", &year, &num)
	if err != nil {
		return nil, err
	}
	
	return &CodeInfo{
		Prefix: "PO",
		Year:   &year,
		Number: num,
	}, nil
}

// GetYearFromPONumber extracts year from PO number
func (cg *CodeGenerator) GetYearFromPONumber(number string) (int, error) {
	info, err := cg.ParsePONumber(number)
	if err != nil {
		return 0, err
	}
	if info.Year == nil {
		return 0, fmt.Errorf("no year found in PO number: %s", number)
	}
	return *info.Year, nil
}

// IsCurrentYearPONumber checks if PO number is from current year
func (cg *CodeGenerator) IsCurrentYearPONumber(number string) bool {
	year, err := cg.GetYearFromPONumber(number)
	if err != nil {
		return false
	}
	return year == time.Now().Year()
}