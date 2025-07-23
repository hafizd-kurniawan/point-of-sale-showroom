package repair

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// VehicleDamage represents damage found on a purchased vehicle
type VehicleDamage struct {
	DamageID         int       `json:"damage_id" db:"damage_id"`
	TransactionID    int       `json:"transaction_id" db:"transaction_id"`
	DamageCategory   string    `json:"damage_category" db:"damage_category"`
	DamageType       string    `json:"damage_type" db:"damage_type"`
	DamageSeverity   string    `json:"damage_severity" db:"damage_severity"`
	DamageLocation   string    `json:"damage_location" db:"damage_location"`
	DamageDescription string   `json:"damage_description" db:"damage_description"`
	EstimatedCost    float64   `json:"estimated_cost" db:"estimated_cost"`
	RepairPriority   int       `json:"repair_priority" db:"repair_priority"`
	RepairRequired   bool      `json:"repair_required" db:"repair_required"`
	DamagePhotosJSON *string   `json:"damage_photos_json,omitempty" db:"damage_photos_json"`
	AssessmentNotes  *string   `json:"assessment_notes,omitempty" db:"assessment_notes"`
	IdentifiedBy     int       `json:"identified_by" db:"identified_by"`
	IdentifiedAt     time.Time `json:"identified_at" db:"identified_at"`
	Status           string    `json:"status" db:"status"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`

	// Related data for joins
	TransactionNumber string `json:"transaction_number,omitempty" db:"transaction_number"`
	VehicleBrand      string `json:"vehicle_brand,omitempty" db:"vehicle_brand"`
	VehicleModel      string `json:"vehicle_model,omitempty" db:"vehicle_model"`
	IdentifiedByName  string `json:"identified_by_name,omitempty" db:"identified_by_name"`
}

// VehicleDamageListItem represents a simplified damage for list views
type VehicleDamageListItem struct {
	DamageID          int       `json:"damage_id" db:"damage_id"`
	TransactionNumber string    `json:"transaction_number" db:"transaction_number"`
	VehicleBrand      string    `json:"vehicle_brand" db:"vehicle_brand"`
	VehicleModel      string    `json:"vehicle_model" db:"vehicle_model"`
	DamageCategory    string    `json:"damage_category" db:"damage_category"`
	DamageSeverity    string    `json:"damage_severity" db:"damage_severity"`
	DamageLocation    string    `json:"damage_location" db:"damage_location"`
	EstimatedCost     float64   `json:"estimated_cost" db:"estimated_cost"`
	RepairPriority    int       `json:"repair_priority" db:"repair_priority"`
	Status            string    `json:"status" db:"status"`
	IdentifiedByName  string    `json:"identified_by_name" db:"identified_by_name"`
	IdentifiedAt      time.Time `json:"identified_at" db:"identified_at"`
}

// VehicleDamageCreateRequest represents a request to create a vehicle damage record
type VehicleDamageCreateRequest struct {
	TransactionID     int     `json:"transaction_id" binding:"required"`
	DamageCategory    string  `json:"damage_category" binding:"required,oneof=body engine interior electrical suspension brake transmission other"`
	DamageType        string  `json:"damage_type" binding:"required,max=100"`
	DamageSeverity    string  `json:"damage_severity" binding:"required,oneof=minor moderate major critical"`
	DamageLocation    string  `json:"damage_location" binding:"required,max=100"`
	DamageDescription string  `json:"damage_description" binding:"required"`
	EstimatedCost     float64 `json:"estimated_cost" binding:"min=0"`
	RepairPriority    int     `json:"repair_priority" binding:"required,min=1,max=5"`
	RepairRequired    bool    `json:"repair_required"`
	DamagePhotosJSON  *string `json:"damage_photos_json,omitempty"`
	AssessmentNotes   *string `json:"assessment_notes,omitempty"`
}

// VehicleDamageUpdateRequest represents a request to update a vehicle damage record
type VehicleDamageUpdateRequest struct {
	DamageType        *string  `json:"damage_type,omitempty" binding:"omitempty,max=100"`
	DamageSeverity    *string  `json:"damage_severity,omitempty" binding:"omitempty,oneof=minor moderate major critical"`
	DamageLocation    *string  `json:"damage_location,omitempty" binding:"omitempty,max=100"`
	DamageDescription *string  `json:"damage_description,omitempty"`
	EstimatedCost     *float64 `json:"estimated_cost,omitempty" binding:"omitempty,min=0"`
	RepairPriority    *int     `json:"repair_priority,omitempty" binding:"omitempty,min=1,max=5"`
	RepairRequired    *bool    `json:"repair_required,omitempty"`
	Status            *string  `json:"status,omitempty" binding:"omitempty,oneof=identified assessed scheduled repairing completed cancelled"`
	DamagePhotosJSON  *string  `json:"damage_photos_json,omitempty"`
	AssessmentNotes   *string  `json:"assessment_notes,omitempty"`
}

// VehicleDamageFilterParams represents filtering parameters for damage queries
type VehicleDamageFilterParams struct {
	TransactionID  *int    `json:"transaction_id,omitempty" form:"transaction_id"`
	DamageCategory string  `json:"damage_category,omitempty" form:"damage_category"`
	DamageSeverity string  `json:"damage_severity,omitempty" form:"damage_severity"`
	Status         string  `json:"status,omitempty" form:"status"`
	RepairPriority *int    `json:"repair_priority,omitempty" form:"repair_priority"`
	RepairRequired *bool   `json:"repair_required,omitempty" form:"repair_required"`
	IdentifiedBy   *int    `json:"identified_by,omitempty" form:"identified_by"`
	MinCost        *float64 `json:"min_cost,omitempty" form:"min_cost"`
	MaxCost        *float64 `json:"max_cost,omitempty" form:"max_cost"`
	Search         string  `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// DamageAssessmentRequest represents a request to assess damage
type DamageAssessmentRequest struct {
	EstimatedCost   float64 `json:"estimated_cost" binding:"required,min=0"`
	RepairPriority  int     `json:"repair_priority" binding:"required,min=1,max=5"`
	RepairRequired  bool    `json:"repair_required"`
	AssessmentNotes string  `json:"assessment_notes" binding:"required"`
	Status          string  `json:"status" binding:"required,oneof=assessed scheduled"`
}

// DamageSummary represents a summary of damages for a transaction
type DamageSummary struct {
	TransactionID       int     `json:"transaction_id" db:"transaction_id"`
	TotalDamages        int     `json:"total_damages" db:"total_damages"`
	TotalEstimatedCost  float64 `json:"total_estimated_cost" db:"total_estimated_cost"`
	CriticalDamages     int     `json:"critical_damages" db:"critical_damages"`
	MajorDamages        int     `json:"major_damages" db:"major_damages"`
	ModerateDamages     int     `json:"moderate_damages" db:"moderate_damages"`
	MinorDamages        int     `json:"minor_damages" db:"minor_damages"`
	RepairsRequired     int     `json:"repairs_required" db:"repairs_required"`
	RepairsCompleted    int     `json:"repairs_completed" db:"repairs_completed"`
	HighPriorityDamages int     `json:"high_priority_damages" db:"high_priority_damages"`
}