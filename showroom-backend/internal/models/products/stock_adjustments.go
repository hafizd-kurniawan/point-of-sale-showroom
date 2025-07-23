package products

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// AdjustmentType represents the type of stock adjustment
type AdjustmentType string

const (
	AdjustmentTypePhysicalCount AdjustmentType = "physical_count"
	AdjustmentTypeDamage        AdjustmentType = "damage"
	AdjustmentTypeExpired       AdjustmentType = "expired"
	AdjustmentTypeTheft         AdjustmentType = "theft"
	AdjustmentTypeCorrection    AdjustmentType = "correction"
	AdjustmentTypeWriteOff      AdjustmentType = "write_off"
)

// IsValid checks if the adjustment type is valid
func (t AdjustmentType) IsValid() bool {
	switch t {
	case AdjustmentTypePhysicalCount, AdjustmentTypeDamage, AdjustmentTypeExpired, AdjustmentTypeTheft, AdjustmentTypeCorrection, AdjustmentTypeWriteOff:
		return true
	default:
		return false
	}
}

// String returns the string representation of the adjustment type
func (t AdjustmentType) String() string {
	return string(t)
}

// Value implements the driver.Valuer interface for AdjustmentType
func (t AdjustmentType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements the sql.Scanner interface for AdjustmentType
func (t *AdjustmentType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*t = AdjustmentType(s)
	case []byte:
		*t = AdjustmentType(s)
	default:
		return fmt.Errorf("cannot scan %T into AdjustmentType", value)
	}
	return nil
}

// StockAdjustment represents a stock adjustment record
type StockAdjustment struct {
	AdjustmentID             int            `json:"adjustment_id" db:"adjustment_id"`
	ProductID                int            `json:"product_id" db:"product_id"`
	AdjustmentType           AdjustmentType `json:"adjustment_type" db:"adjustment_type"`
	QuantitySystem           int            `json:"quantity_system" db:"quantity_system"`
	QuantityPhysical         int            `json:"quantity_physical" db:"quantity_physical"`
	QuantityVariance         int            `json:"quantity_variance" db:"quantity_variance"`
	CostImpact               float64        `json:"cost_impact" db:"cost_impact"`
	AdjustmentReason         string         `json:"adjustment_reason" db:"adjustment_reason"`
	Notes                    *string        `json:"notes,omitempty" db:"notes"`
	ApprovedBy               *int           `json:"approved_by,omitempty" db:"approved_by"`
	AdjustmentDate           time.Time      `json:"adjustment_date" db:"adjustment_date"`
	ApprovedAt               *time.Time     `json:"approved_at,omitempty" db:"approved_at"`
	SupportingDocumentsJSON  *string        `json:"supporting_documents_json,omitempty" db:"supporting_documents_json"`
	CreatedAt                time.Time      `json:"created_at" db:"created_at"`
	CreatedBy                int            `json:"created_by" db:"created_by"`
}

// StockAdjustmentListItem represents a simplified stock adjustment for list views
type StockAdjustmentListItem struct {
	AdjustmentID     int            `json:"adjustment_id" db:"adjustment_id"`
	ProductID        int            `json:"product_id" db:"product_id"`
	AdjustmentType   AdjustmentType `json:"adjustment_type" db:"adjustment_type"`
	QuantitySystem   int            `json:"quantity_system" db:"quantity_system"`
	QuantityPhysical int            `json:"quantity_physical" db:"quantity_physical"`
	QuantityVariance int            `json:"quantity_variance" db:"quantity_variance"`
	CostImpact       float64        `json:"cost_impact" db:"cost_impact"`
	AdjustmentReason string         `json:"adjustment_reason" db:"adjustment_reason"`
	AdjustmentDate   time.Time      `json:"adjustment_date" db:"adjustment_date"`
	ApprovedBy       *int           `json:"approved_by,omitempty" db:"approved_by"`
	CreatedBy        int            `json:"created_by" db:"created_by"`
}

// StockAdjustmentCreateRequest represents a request to create a stock adjustment
type StockAdjustmentCreateRequest struct {
	ProductID                int            `json:"product_id" binding:"required,min=1"`
	AdjustmentType           AdjustmentType `json:"adjustment_type" binding:"required"`
	QuantityPhysical         int            `json:"quantity_physical" binding:"required,min=0"`
	AdjustmentReason         string         `json:"adjustment_reason" binding:"required,max=255"`
	Notes                    *string        `json:"notes,omitempty"`
	AdjustmentDate           *time.Time     `json:"adjustment_date,omitempty"`
	SupportingDocumentsJSON  *string        `json:"supporting_documents_json,omitempty"`
}

// StockAdjustmentUpdateRequest represents a request to update a stock adjustment
type StockAdjustmentUpdateRequest struct {
	AdjustmentType           *AdjustmentType `json:"adjustment_type,omitempty"`
	QuantityPhysical         *int            `json:"quantity_physical,omitempty" binding:"omitempty,min=0"`
	AdjustmentReason         *string         `json:"adjustment_reason,omitempty" binding:"omitempty,max=255"`
	Notes                    *string         `json:"notes,omitempty"`
	AdjustmentDate           *time.Time      `json:"adjustment_date,omitempty"`
	SupportingDocumentsJSON  *string         `json:"supporting_documents_json,omitempty"`
}

// StockAdjustmentFilterParams represents filtering parameters for stock adjustment queries
type StockAdjustmentFilterParams struct {
	ProductID      *int            `json:"product_id,omitempty" form:"product_id"`
	AdjustmentType *AdjustmentType `json:"adjustment_type,omitempty" form:"adjustment_type"`
	ApprovedBy     *int            `json:"approved_by,omitempty" form:"approved_by"`
	CreatedBy      *int            `json:"created_by,omitempty" form:"created_by"`
	DateFrom       *time.Time      `json:"date_from,omitempty" form:"date_from"`
	DateTo         *time.Time      `json:"date_to,omitempty" form:"date_to"`
	IsApproved     *bool           `json:"is_approved,omitempty" form:"is_approved"`
	HasVariance    *bool           `json:"has_variance,omitempty" form:"has_variance"`
	Search         string          `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// CalculateVariance calculates the quantity variance and cost impact
func (sa *StockAdjustment) CalculateVariance(unitCost float64) {
	sa.QuantityVariance = sa.QuantityPhysical - sa.QuantitySystem
	sa.CostImpact = float64(sa.QuantityVariance) * unitCost
}

// IsApproved checks if the adjustment is approved
func (sa *StockAdjustment) IsApproved() bool {
	return sa.ApprovedBy != nil && sa.ApprovedAt != nil
}

// CanApprove checks if the adjustment can be approved
func (sa *StockAdjustment) CanApprove() bool {
	return !sa.IsApproved()
}

// Approve approves the stock adjustment
func (sa *StockAdjustment) Approve(approvedBy int) {
	if !sa.CanApprove() {
		return
	}
	sa.ApprovedBy = &approvedBy
	now := time.Now()
	sa.ApprovedAt = &now
}

// HasPositiveVariance checks if there is a positive variance (more physical than system)
func (sa *StockAdjustment) HasPositiveVariance() bool {
	return sa.QuantityVariance > 0
}

// HasNegativeVariance checks if there is a negative variance (less physical than system)
func (sa *StockAdjustment) HasNegativeVariance() bool {
	return sa.QuantityVariance < 0
}

// HasVariance checks if there is any variance
func (sa *StockAdjustment) HasVariance() bool {
	return sa.QuantityVariance != 0
}

// GetVarianceType returns a description of the variance type
func (sa *StockAdjustment) GetVarianceType() string {
	if sa.QuantityVariance > 0 {
		return "surplus"
	} else if sa.QuantityVariance < 0 {
		return "shortage"
	}
	return "none"
}

// GetAdjustmentDescription returns a human-readable description of the adjustment
func (sa *StockAdjustment) GetAdjustmentDescription() string {
	varianceType := sa.GetVarianceType()
	varianceDesc := ""
	
	if sa.QuantityVariance != 0 {
		varianceDesc = fmt.Sprintf(" (%s of %d units)", varianceType, abs(sa.QuantityVariance))
	}
	
	switch sa.AdjustmentType {
	case AdjustmentTypePhysicalCount:
		return fmt.Sprintf("Physical count adjustment%s", varianceDesc)
	case AdjustmentTypeDamage:
		return fmt.Sprintf("Damage adjustment: %d units damaged", abs(sa.QuantityVariance))
	case AdjustmentTypeExpired:
		return fmt.Sprintf("Expiry adjustment: %d units expired", abs(sa.QuantityVariance))
	case AdjustmentTypeTheft:
		return fmt.Sprintf("Theft adjustment: %d units stolen", abs(sa.QuantityVariance))
	case AdjustmentTypeCorrection:
		return fmt.Sprintf("System correction%s", varianceDesc)
	case AdjustmentTypeWriteOff:
		return fmt.Sprintf("Write-off adjustment: %d units written off", abs(sa.QuantityVariance))
	default:
		return fmt.Sprintf("%s adjustment%s", sa.AdjustmentType, varianceDesc)
	}
}

// RequiresApproval checks if the adjustment requires approval based on variance
func (sa *StockAdjustment) RequiresApproval(threshold int) bool {
	return abs(sa.QuantityVariance) >= threshold || sa.AdjustmentType == AdjustmentTypeWriteOff
}

// GetFinancialImpact returns the financial impact as positive or negative
func (sa *StockAdjustment) GetFinancialImpact() float64 {
	return sa.CostImpact
}

// IsSignificantVariance checks if the variance is significant based on percentage
func (sa *StockAdjustment) IsSignificantVariance(threshold float64) bool {
	if sa.QuantitySystem == 0 {
		return sa.QuantityVariance != 0
	}
	percentage := float64(abs(sa.QuantityVariance)) / float64(sa.QuantitySystem) * 100
	return percentage >= threshold
}

// helper function to get absolute value
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}