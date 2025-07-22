package inventory

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// MovementType represents the type of stock movement
type MovementType string

const (
	MovementTypeIn         MovementType = "in"
	MovementTypeOut        MovementType = "out"
	MovementTypeTransfer   MovementType = "transfer"
	MovementTypeAdjustment MovementType = "adjustment"
	MovementTypeDamage     MovementType = "damage"
	MovementTypeExpired    MovementType = "expired"
	MovementTypeReturn     MovementType = "return"
)

// ReferenceType represents the reference type for stock movement
type ReferenceType string

const (
	ReferenceTypePurchase   ReferenceType = "purchase"
	ReferenceTypeSales      ReferenceType = "sales"
	ReferenceTypeRepair     ReferenceType = "repair"
	ReferenceTypeAdjustment ReferenceType = "adjustment"
	ReferenceTypeTransfer   ReferenceType = "transfer"
	ReferenceTypeReturn     ReferenceType = "return"
)

// StockMovement represents a stock movement record in the system
type StockMovement struct {
	MovementID     int           `json:"movement_id" db:"movement_id"`
	ProductID      int           `json:"product_id" db:"product_id"`
	MovementType   MovementType  `json:"movement_type" db:"movement_type"`
	ReferenceType  ReferenceType `json:"reference_type" db:"reference_type"`
	ReferenceID    *int          `json:"reference_id,omitempty" db:"reference_id"`
	QuantityBefore int           `json:"quantity_before" db:"quantity_before"`
	QuantityMoved  int           `json:"quantity_moved" db:"quantity_moved"`
	QuantityAfter  int           `json:"quantity_after" db:"quantity_after"`
	UnitCost       float64       `json:"unit_cost" db:"unit_cost"`
	TotalValue     float64       `json:"total_value" db:"total_value"`
	LocationFrom   *string       `json:"location_from,omitempty" db:"location_from"`
	LocationTo     *string       `json:"location_to,omitempty" db:"location_to"`
	MovementDate   time.Time     `json:"movement_date" db:"movement_date"`
	ProcessedBy    int           `json:"processed_by" db:"processed_by"`
	MovementReason *string       `json:"movement_reason,omitempty" db:"movement_reason"`
	Notes          *string       `json:"notes,omitempty" db:"notes"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`

	// Related data
	ProductCode    string  `json:"product_code,omitempty" db:"product_code"`
	ProductName    string  `json:"product_name,omitempty" db:"product_name"`
	UnitMeasure    string  `json:"unit_measure,omitempty" db:"unit_measure"`
	ProcessedByName string `json:"processed_by_name,omitempty" db:"processed_by_name"`
}

// StockMovementListItem represents a simplified stock movement for list views
type StockMovementListItem struct {
	MovementID      int           `json:"movement_id" db:"movement_id"`
	ProductCode     string        `json:"product_code" db:"product_code"`
	ProductName     string        `json:"product_name" db:"product_name"`
	MovementType    MovementType  `json:"movement_type" db:"movement_type"`
	ReferenceType   ReferenceType `json:"reference_type" db:"reference_type"`
	ReferenceID     *int          `json:"reference_id,omitempty" db:"reference_id"`
	QuantityMoved   int           `json:"quantity_moved" db:"quantity_moved"`
	UnitCost        float64       `json:"unit_cost" db:"unit_cost"`
	TotalValue      float64       `json:"total_value" db:"total_value"`
	MovementDate    time.Time     `json:"movement_date" db:"movement_date"`
	ProcessedByName string        `json:"processed_by_name" db:"processed_by_name"`
	MovementReason  *string       `json:"movement_reason,omitempty" db:"movement_reason"`
}

// StockMovementCreateRequest represents a request to create a stock movement
type StockMovementCreateRequest struct {
	ProductID      int           `json:"product_id" binding:"required"`
	MovementType   MovementType  `json:"movement_type" binding:"required,oneof=in out transfer adjustment damage expired return"`
	ReferenceType  ReferenceType `json:"reference_type" binding:"required,oneof=purchase sales repair adjustment transfer return"`
	ReferenceID    *int          `json:"reference_id,omitempty"`
	QuantityMoved  int           `json:"quantity_moved" binding:"required,ne=0"`
	UnitCost       float64       `json:"unit_cost" binding:"required,min=0"`
	LocationFrom   *string       `json:"location_from,omitempty" binding:"omitempty,max=100"`
	LocationTo     *string       `json:"location_to,omitempty" binding:"omitempty,max=100"`
	MovementReason *string       `json:"movement_reason,omitempty" binding:"omitempty,max=255"`
	Notes          *string       `json:"notes,omitempty"`
}

// StockMovementFilterParams represents filtering parameters for stock movement queries
type StockMovementFilterParams struct {
	ProductID     *int          `json:"product_id,omitempty" form:"product_id"`
	MovementType  *MovementType `json:"movement_type,omitempty" form:"movement_type"`
	ReferenceType *ReferenceType `json:"reference_type,omitempty" form:"reference_type"`
	ReferenceID   *int          `json:"reference_id,omitempty" form:"reference_id"`
	ProcessedBy   *int          `json:"processed_by,omitempty" form:"processed_by"`
	LocationFrom  string        `json:"location_from,omitempty" form:"location_from"`
	LocationTo    string        `json:"location_to,omitempty" form:"location_to"`
	DateFrom      *time.Time    `json:"date_from,omitempty" form:"date_from"`
	DateTo        *time.Time    `json:"date_to,omitempty" form:"date_to"`
	Search        string        `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

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

// StockAdjustment represents a stock adjustment record in the system
type StockAdjustment struct {
	AdjustmentID               int            `json:"adjustment_id" db:"adjustment_id"`
	ProductID                  int            `json:"product_id" db:"product_id"`
	AdjustmentType             AdjustmentType `json:"adjustment_type" db:"adjustment_type"`
	QuantitySystem             int            `json:"quantity_system" db:"quantity_system"`
	QuantityPhysical           int            `json:"quantity_physical" db:"quantity_physical"`
	QuantityVariance           int            `json:"quantity_variance" db:"quantity_variance"`
	CostImpact                 float64        `json:"cost_impact" db:"cost_impact"`
	AdjustmentReason           string         `json:"adjustment_reason" db:"adjustment_reason"`
	Notes                      *string        `json:"notes,omitempty" db:"notes"`
	ApprovedBy                 *int           `json:"approved_by,omitempty" db:"approved_by"`
	AdjustmentDate             time.Time      `json:"adjustment_date" db:"adjustment_date"`
	ApprovedAt                 *time.Time     `json:"approved_at,omitempty" db:"approved_at"`
	SupportingDocumentsJSON    *string        `json:"supporting_documents_json,omitempty" db:"supporting_documents_json"`
	CreatedAt                  time.Time      `json:"created_at" db:"created_at"`
	CreatedBy                  int            `json:"created_by" db:"created_by"`

	// Related data
	ProductCode     string  `json:"product_code,omitempty" db:"product_code"`
	ProductName     string  `json:"product_name,omitempty" db:"product_name"`
	UnitMeasure     string  `json:"unit_measure,omitempty" db:"unit_measure"`
	CreatedByName   string  `json:"created_by_name,omitempty" db:"created_by_name"`
	ApprovedByName  *string `json:"approved_by_name,omitempty" db:"approved_by_name"`
}

// StockAdjustmentListItem represents a simplified stock adjustment for list views
type StockAdjustmentListItem struct {
	AdjustmentID     int            `json:"adjustment_id" db:"adjustment_id"`
	ProductCode      string         `json:"product_code" db:"product_code"`
	ProductName      string         `json:"product_name" db:"product_name"`
	AdjustmentType   AdjustmentType `json:"adjustment_type" db:"adjustment_type"`
	QuantityVariance int            `json:"quantity_variance" db:"quantity_variance"`
	CostImpact       float64        `json:"cost_impact" db:"cost_impact"`
	AdjustmentReason string         `json:"adjustment_reason" db:"adjustment_reason"`
	AdjustmentDate   time.Time      `json:"adjustment_date" db:"adjustment_date"`
	CreatedByName    string         `json:"created_by_name" db:"created_by_name"`
	ApprovedByName   *string        `json:"approved_by_name,omitempty" db:"approved_by_name"`
	ApprovedAt       *time.Time     `json:"approved_at,omitempty" db:"approved_at"`
}

// StockAdjustmentCreateRequest represents a request to create a stock adjustment
type StockAdjustmentCreateRequest struct {
	ProductID                int            `json:"product_id" binding:"required"`
	AdjustmentType           AdjustmentType `json:"adjustment_type" binding:"required,oneof=physical_count damage expired theft correction write_off"`
	QuantityPhysical         int            `json:"quantity_physical" binding:"required,min=0"`
	AdjustmentReason         string         `json:"adjustment_reason" binding:"required,max=255"`
	Notes                    *string        `json:"notes,omitempty"`
	SupportingDocumentsJSON  *string        `json:"supporting_documents_json,omitempty"`
}

// StockAdjustmentUpdateRequest represents a request to update a stock adjustment
type StockAdjustmentUpdateRequest struct {
	AdjustmentType          *AdjustmentType `json:"adjustment_type,omitempty" binding:"omitempty,oneof=physical_count damage expired theft correction write_off"`
	QuantityPhysical        *int            `json:"quantity_physical,omitempty" binding:"omitempty,min=0"`
	AdjustmentReason        *string         `json:"adjustment_reason,omitempty" binding:"omitempty,max=255"`
	Notes                   *string         `json:"notes,omitempty"`
	SupportingDocumentsJSON *string         `json:"supporting_documents_json,omitempty"`
}

// StockAdjustmentFilterParams represents filtering parameters for stock adjustment queries
type StockAdjustmentFilterParams struct {
	ProductID      *int           `json:"product_id,omitempty" form:"product_id"`
	AdjustmentType *AdjustmentType `json:"adjustment_type,omitempty" form:"adjustment_type"`
	CreatedBy      *int           `json:"created_by,omitempty" form:"created_by"`
	ApprovedBy     *int           `json:"approved_by,omitempty" form:"approved_by"`
	DateFrom       *time.Time     `json:"date_from,omitempty" form:"date_from"`
	DateTo         *time.Time     `json:"date_to,omitempty" form:"date_to"`
	IsApproved     *bool          `json:"is_approved,omitempty" form:"is_approved"`
	Search         string         `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// Methods for StockMovement

// CalculateTotalValue calculates and updates the total value of the movement
func (sm *StockMovement) CalculateTotalValue() {
	// For outgoing movements, use negative quantity to represent value reduction
	if sm.MovementType == MovementTypeOut || sm.MovementType == MovementTypeDamage || 
		sm.MovementType == MovementTypeExpired {
		sm.TotalValue = float64(-sm.QuantityMoved) * sm.UnitCost
	} else {
		sm.TotalValue = float64(sm.QuantityMoved) * sm.UnitCost
	}
}

// IsIncoming checks if the movement is incoming (adds stock)
func (sm *StockMovement) IsIncoming() bool {
	return sm.MovementType == MovementTypeIn || sm.MovementType == MovementTypeReturn
}

// IsOutgoing checks if the movement is outgoing (reduces stock)
func (sm *StockMovement) IsOutgoing() bool {
	return sm.MovementType == MovementTypeOut || sm.MovementType == MovementTypeDamage || 
		sm.MovementType == MovementTypeExpired
}

// IsTransfer checks if the movement is a transfer
func (sm *StockMovement) IsTransfer() bool {
	return sm.MovementType == MovementTypeTransfer
}

// IsAdjustment checks if the movement is an adjustment
func (sm *StockMovement) IsAdjustment() bool {
	return sm.MovementType == MovementTypeAdjustment
}

// RequiresLocation checks if the movement type requires location information
func (sm *StockMovement) RequiresLocation() bool {
	return sm.MovementType == MovementTypeTransfer
}

// Methods for StockAdjustment

// CalculateVariance calculates and updates the quantity variance
func (sa *StockAdjustment) CalculateVariance() {
	sa.QuantityVariance = sa.QuantityPhysical - sa.QuantitySystem
}

// CalculateCostImpact calculates and updates the cost impact of the adjustment
func (sa *StockAdjustment) CalculateCostImpact(unitCost float64) {
	sa.CostImpact = float64(sa.QuantityVariance) * unitCost
}

// IsPositiveAdjustment checks if the adjustment increases stock
func (sa *StockAdjustment) IsPositiveAdjustment() bool {
	return sa.QuantityVariance > 0
}

// IsNegativeAdjustment checks if the adjustment decreases stock
func (sa *StockAdjustment) IsNegativeAdjustment() bool {
	return sa.QuantityVariance < 0
}

// IsApproved checks if the adjustment is approved
func (sa *StockAdjustment) IsApproved() bool {
	return sa.ApprovedBy != nil && sa.ApprovedAt != nil
}

// CanBeApproved checks if the adjustment can be approved
func (sa *StockAdjustment) CanBeApproved() bool {
	return !sa.IsApproved()
}

// RequiresApproval checks if the adjustment requires approval based on variance
func (sa *StockAdjustment) RequiresApproval() bool {
	// Require approval for variances greater than certain threshold or negative adjustments
	return sa.QuantityVariance != 0 && (sa.QuantityVariance < 0 || sa.QuantityVariance > 10)
}