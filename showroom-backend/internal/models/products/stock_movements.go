package products

import (
	"database/sql/driver"
	"fmt"
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

// IsValid checks if the movement type is valid
func (t MovementType) IsValid() bool {
	switch t {
	case MovementTypeIn, MovementTypeOut, MovementTypeTransfer, MovementTypeAdjustment, MovementTypeDamage, MovementTypeExpired, MovementTypeReturn:
		return true
	default:
		return false
	}
}

// String returns the string representation of the movement type
func (t MovementType) String() string {
	return string(t)
}

// Value implements the driver.Valuer interface for MovementType
func (t MovementType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements the sql.Scanner interface for MovementType
func (t *MovementType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*t = MovementType(s)
	case []byte:
		*t = MovementType(s)
	default:
		return fmt.Errorf("cannot scan %T into MovementType", value)
	}
	return nil
}

// ReferenceType represents the type of reference for stock movement
type ReferenceType string

const (
	ReferenceTypePurchase   ReferenceType = "purchase"
	ReferenceTypeSales      ReferenceType = "sales"
	ReferenceTypeRepair     ReferenceType = "repair"
	ReferenceTypeAdjustment ReferenceType = "adjustment"
	ReferenceTypeTransfer   ReferenceType = "transfer"
	ReferenceTypeReturn     ReferenceType = "return"
)

// IsValid checks if the reference type is valid
func (r ReferenceType) IsValid() bool {
	switch r {
	case ReferenceTypePurchase, ReferenceTypeSales, ReferenceTypeRepair, ReferenceTypeAdjustment, ReferenceTypeTransfer, ReferenceTypeReturn:
		return true
	default:
		return false
	}
}

// String returns the string representation of the reference type
func (r ReferenceType) String() string {
	return string(r)
}

// Value implements the driver.Valuer interface for ReferenceType
func (r ReferenceType) Value() (driver.Value, error) {
	return string(r), nil
}

// Scan implements the sql.Scanner interface for ReferenceType
func (r *ReferenceType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*r = ReferenceType(s)
	case []byte:
		*r = ReferenceType(s)
	default:
		return fmt.Errorf("cannot scan %T into ReferenceType", value)
	}
	return nil
}

// StockMovement represents a stock movement record
type StockMovement struct {
	MovementID      int           `json:"movement_id" db:"movement_id"`
	ProductID       int           `json:"product_id" db:"product_id"`
	MovementType    MovementType  `json:"movement_type" db:"movement_type"`
	ReferenceType   ReferenceType `json:"reference_type" db:"reference_type"`
	ReferenceID     int           `json:"reference_id" db:"reference_id"`
	QuantityBefore  int           `json:"quantity_before" db:"quantity_before"`
	QuantityMoved   int           `json:"quantity_moved" db:"quantity_moved"`
	QuantityAfter   int           `json:"quantity_after" db:"quantity_after"`
	UnitCost        float64       `json:"unit_cost" db:"unit_cost"`
	TotalValue      float64       `json:"total_value" db:"total_value"`
	LocationFrom    *string       `json:"location_from,omitempty" db:"location_from"`
	LocationTo      *string       `json:"location_to,omitempty" db:"location_to"`
	MovementDate    time.Time     `json:"movement_date" db:"movement_date"`
	ProcessedBy     int           `json:"processed_by" db:"processed_by"`
	MovementReason  *string       `json:"movement_reason,omitempty" db:"movement_reason"`
	Notes           *string       `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
}

// StockMovementListItem represents a simplified stock movement for list views
type StockMovementListItem struct {
	MovementID     int           `json:"movement_id" db:"movement_id"`
	ProductID      int           `json:"product_id" db:"product_id"`
	MovementType   MovementType  `json:"movement_type" db:"movement_type"`
	ReferenceType  ReferenceType `json:"reference_type" db:"reference_type"`
	ReferenceID    int           `json:"reference_id" db:"reference_id"`
	QuantityBefore int           `json:"quantity_before" db:"quantity_before"`
	QuantityMoved  int           `json:"quantity_moved" db:"quantity_moved"`
	QuantityAfter  int           `json:"quantity_after" db:"quantity_after"`
	UnitCost       float64       `json:"unit_cost" db:"unit_cost"`
	TotalValue     float64       `json:"total_value" db:"total_value"`
	MovementDate   time.Time     `json:"movement_date" db:"movement_date"`
	ProcessedBy    int           `json:"processed_by" db:"processed_by"`
}

// StockMovementCreateRequest represents a request to create a stock movement
type StockMovementCreateRequest struct {
	ProductID       int           `json:"product_id" binding:"required,min=1"`
	MovementType    MovementType  `json:"movement_type" binding:"required"`
	ReferenceType   ReferenceType `json:"reference_type" binding:"required"`
	ReferenceID     int           `json:"reference_id" binding:"required,min=1"`
	QuantityMoved   int           `json:"quantity_moved" binding:"required"`
	UnitCost        float64       `json:"unit_cost" binding:"required,min=0"`
	LocationFrom    *string       `json:"location_from,omitempty" binding:"omitempty,max=100"`
	LocationTo      *string       `json:"location_to,omitempty" binding:"omitempty,max=100"`
	MovementDate    *time.Time    `json:"movement_date,omitempty"`
	MovementReason  *string       `json:"movement_reason,omitempty" binding:"omitempty,max=255"`
	Notes           *string       `json:"notes,omitempty"`
}

// StockMovementFilterParams represents filtering parameters for stock movement queries
type StockMovementFilterParams struct {
	ProductID     *int           `json:"product_id,omitempty" form:"product_id"`
	MovementType  *MovementType  `json:"movement_type,omitempty" form:"movement_type"`
	ReferenceType *ReferenceType `json:"reference_type,omitempty" form:"reference_type"`
	ReferenceID   *int           `json:"reference_id,omitempty" form:"reference_id"`
	ProcessedBy   *int           `json:"processed_by,omitempty" form:"processed_by"`
	DateFrom      *time.Time     `json:"date_from,omitempty" form:"date_from"`
	DateTo        *time.Time     `json:"date_to,omitempty" form:"date_to"`
	LocationFrom  *string        `json:"location_from,omitempty" form:"location_from"`
	LocationTo    *string        `json:"location_to,omitempty" form:"location_to"`
	Search        string         `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// CalculateTotalValue calculates the total value of the movement
func (sm *StockMovement) CalculateTotalValue() {
	// For negative movements (out, damage, expired), quantity moved should be negative
	actualQuantity := sm.QuantityMoved
	if sm.MovementType == MovementTypeOut || sm.MovementType == MovementTypeDamage || sm.MovementType == MovementTypeExpired {
		actualQuantity = -actualQuantity
	}
	sm.TotalValue = float64(actualQuantity) * sm.UnitCost
}

// ValidateQuantities validates the quantity relationships
func (sm *StockMovement) ValidateQuantities() bool {
	switch sm.MovementType {
	case MovementTypeIn, MovementTypeReturn:
		return sm.QuantityAfter == sm.QuantityBefore + sm.QuantityMoved
	case MovementTypeOut, MovementTypeDamage, MovementTypeExpired:
		return sm.QuantityAfter == sm.QuantityBefore - sm.QuantityMoved
	case MovementTypeAdjustment:
		// For adjustments, quantity moved can be positive or negative
		return sm.QuantityAfter == sm.QuantityBefore + sm.QuantityMoved
	case MovementTypeTransfer:
		// For transfers, we need additional logic to handle source and destination
		return true // Will be validated in business logic
	}
	return false
}

// IsInbound checks if the movement increases stock
func (sm *StockMovement) IsInbound() bool {
	return sm.MovementType == MovementTypeIn || sm.MovementType == MovementTypeReturn ||
		(sm.MovementType == MovementTypeAdjustment && sm.QuantityMoved > 0)
}

// IsOutbound checks if the movement decreases stock
func (sm *StockMovement) IsOutbound() bool {
	return sm.MovementType == MovementTypeOut || sm.MovementType == MovementTypeDamage ||
		sm.MovementType == MovementTypeExpired ||
		(sm.MovementType == MovementTypeAdjustment && sm.QuantityMoved < 0)
}

// GetMovementDescription returns a human-readable description of the movement
func (sm *StockMovement) GetMovementDescription() string {
	switch sm.MovementType {
	case MovementTypeIn:
		return fmt.Sprintf("Stock in: +%d units", sm.QuantityMoved)
	case MovementTypeOut:
		return fmt.Sprintf("Stock out: -%d units", sm.QuantityMoved)
	case MovementTypeTransfer:
		return fmt.Sprintf("Transfer: %d units from %s to %s", sm.QuantityMoved, 
			getStringValue(sm.LocationFrom), getStringValue(sm.LocationTo))
	case MovementTypeAdjustment:
		if sm.QuantityMoved > 0 {
			return fmt.Sprintf("Adjustment: +%d units", sm.QuantityMoved)
		} else {
			return fmt.Sprintf("Adjustment: %d units", sm.QuantityMoved)
		}
	case MovementTypeDamage:
		return fmt.Sprintf("Damage: -%d units", sm.QuantityMoved)
	case MovementTypeExpired:
		return fmt.Sprintf("Expired: -%d units", sm.QuantityMoved)
	case MovementTypeReturn:
		return fmt.Sprintf("Return: +%d units", sm.QuantityMoved)
	default:
		return fmt.Sprintf("%s: %d units", sm.MovementType, sm.QuantityMoved)
	}
}

// helper function to get string value with default
func getStringValue(s *string) string {
	if s == nil {
		return "Unknown"
	}
	return *s
}