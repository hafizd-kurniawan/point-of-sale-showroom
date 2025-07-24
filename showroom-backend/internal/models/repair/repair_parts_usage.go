package repair

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// RepairPartsUsage represents parts used in repair work
type RepairPartsUsage struct {
	UsageID            int       `json:"usage_id" db:"usage_id"`
	WorkDetailID       int       `json:"work_detail_id" db:"work_detail_id"`
	ProductID          int       `json:"product_id" db:"product_id"`
	QuantityUsed       int       `json:"quantity_used" db:"quantity_used"`
	UnitCost           float64   `json:"unit_cost" db:"unit_cost"`
	TotalCost          float64   `json:"total_cost" db:"total_cost"`
	UsageDate          time.Time `json:"usage_date" db:"usage_date"`
	UsageType          string    `json:"usage_type" db:"usage_type"`
	PartCondition      string    `json:"part_condition" db:"part_condition"`
	WarrantyPeriodDays int       `json:"warranty_period_days" db:"warranty_period_days"`
	InstallationNotes  *string   `json:"installation_notes,omitempty" db:"installation_notes"`
	IssuedBy           int       `json:"issued_by" db:"issued_by"`
	UsedBy             int       `json:"used_by" db:"used_by"`
	ApprovedBy         *int      `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`

	// Related data for joins
	ProductName        string `json:"product_name,omitempty" db:"product_name"`
	ProductCode        string `json:"product_code,omitempty" db:"product_code"`
	WorkOrderNumber    string `json:"work_order_number,omitempty" db:"work_order_number"`
	TaskDescription    string `json:"task_description,omitempty" db:"task_description"`
	IssuedByName       string `json:"issued_by_name,omitempty" db:"issued_by_name"`
	UsedByName         string `json:"used_by_name,omitempty" db:"used_by_name"`
	ApprovedByName     string `json:"approved_by_name,omitempty" db:"approved_by_name"`
}

// RepairPartsUsageListItem represents a simplified parts usage for list views
type RepairPartsUsageListItem struct {
	UsageID         int       `json:"usage_id" db:"usage_id"`
	WorkOrderNumber string    `json:"work_order_number" db:"work_order_number"`
	ProductName     string    `json:"product_name" db:"product_name"`
	ProductCode     string    `json:"product_code" db:"product_code"`
	QuantityUsed    int       `json:"quantity_used" db:"quantity_used"`
	UnitCost        float64   `json:"unit_cost" db:"unit_cost"`
	TotalCost       float64   `json:"total_cost" db:"total_cost"`
	UsageType       string    `json:"usage_type" db:"usage_type"`
	PartCondition   string    `json:"part_condition" db:"part_condition"`
	UsageDate       time.Time `json:"usage_date" db:"usage_date"`
	UsedByName      string    `json:"used_by_name" db:"used_by_name"`
	IssuedByName    string    `json:"issued_by_name" db:"issued_by_name"`
}

// RepairPartsUsageCreateRequest represents a request to create parts usage
type RepairPartsUsageCreateRequest struct {
	WorkDetailID       int     `json:"work_detail_id" binding:"required"`
	ProductID          int     `json:"product_id" binding:"required"`
	QuantityUsed       int     `json:"quantity_used" binding:"required,min=1"`
	UnitCost           float64 `json:"unit_cost" binding:"required,min=0"`
	UsageType          string  `json:"usage_type" binding:"required,oneof=new replacement additional warranty"`
	PartCondition      string  `json:"part_condition" binding:"required,oneof=new refurbished used oem aftermarket"`
	WarrantyPeriodDays int     `json:"warranty_period_days" binding:"min=0"`
	InstallationNotes  *string `json:"installation_notes,omitempty"`
	UsedBy             int     `json:"used_by" binding:"required"`
}

// RepairPartsUsageUpdateRequest represents a request to update parts usage
type RepairPartsUsageUpdateRequest struct {
	QuantityUsed       *int     `json:"quantity_used,omitempty" binding:"omitempty,min=1"`
	UnitCost           *float64 `json:"unit_cost,omitempty" binding:"omitempty,min=0"`
	UsageType          *string  `json:"usage_type,omitempty" binding:"omitempty,oneof=new replacement additional warranty"`
	PartCondition      *string  `json:"part_condition,omitempty" binding:"omitempty,oneof=new refurbished used oem aftermarket"`
	WarrantyPeriodDays *int     `json:"warranty_period_days,omitempty" binding:"omitempty,min=0"`
	InstallationNotes  *string  `json:"installation_notes,omitempty"`
}

// RepairPartsUsageFilterParams represents filtering parameters for parts usage queries
type RepairPartsUsageFilterParams struct {
	WorkDetailID  *int       `json:"work_detail_id,omitempty" form:"work_detail_id"`
	ProductID     *int       `json:"product_id,omitempty" form:"product_id"`
	UsageType     string     `json:"usage_type,omitempty" form:"usage_type"`
	PartCondition string     `json:"part_condition,omitempty" form:"part_condition"`
	IssuedBy      *int       `json:"issued_by,omitempty" form:"issued_by"`
	UsedBy        *int       `json:"used_by,omitempty" form:"used_by"`
	StartDate     *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate       *time.Time `json:"end_date,omitempty" form:"end_date"`
	Search        string     `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// PartsUsageApprovalRequest represents a request to approve parts usage
type PartsUsageApprovalRequest struct {
	Status        string  `json:"status" binding:"required,oneof=approved rejected"`
	ApprovalNotes *string `json:"approval_notes,omitempty"`
}

// PartsUsageIssueRequest represents a request to issue parts for repair
type PartsUsageIssueRequest struct {
	Items []PartsIssueItem `json:"items" binding:"required,dive"`
}

// PartsIssueItem represents an individual part to be issued
type PartsIssueItem struct {
	ProductID          int     `json:"product_id" binding:"required"`
	QuantityRequested  int     `json:"quantity_requested" binding:"required,min=1"`
	UsageType          string  `json:"usage_type" binding:"required,oneof=new replacement additional warranty"`
	PartCondition      string  `json:"part_condition" binding:"required,oneof=new refurbished used oem aftermarket"`
	WarrantyPeriodDays int     `json:"warranty_period_days" binding:"min=0"`
	InstallationNotes  *string `json:"installation_notes,omitempty"`
}

// PartsUsageSummary represents a summary of parts usage for a work order
type PartsUsageSummary struct {
	WorkOrderID       int     `json:"work_order_id" db:"work_order_id"`
	TotalPartsUsed    int     `json:"total_parts_used" db:"total_parts_used"`
	TotalPartsCost    float64 `json:"total_parts_cost" db:"total_parts_cost"`
	NewParts          int     `json:"new_parts" db:"new_parts"`
	ReplacementParts  int     `json:"replacement_parts" db:"replacement_parts"`
	WarrantyParts     int     `json:"warranty_parts" db:"warranty_parts"`
	OEMParts          int     `json:"oem_parts" db:"oem_parts"`
	AftermarketParts  int     `json:"aftermarket_parts" db:"aftermarket_parts"`
	ApprovedUsage     int     `json:"approved_usage" db:"approved_usage"`
	PendingApproval   int     `json:"pending_approval" db:"pending_approval"`
}

// PartsInventoryImpact represents the impact of parts usage on inventory
type PartsInventoryImpact struct {
	ProductID          int     `json:"product_id" db:"product_id"`
	ProductName        string  `json:"product_name" db:"product_name"`
	ProductCode        string  `json:"product_code" db:"product_code"`
	QuantityBefore     int     `json:"quantity_before" db:"quantity_before"`
	QuantityUsed       int     `json:"quantity_used" db:"quantity_used"`
	QuantityAfter      int     `json:"quantity_after" db:"quantity_after"`
	TotalValue         float64 `json:"total_value" db:"total_value"`
	LowStockAlert      bool    `json:"low_stock_alert" db:"low_stock_alert"`
	ReorderRequired    bool    `json:"reorder_required" db:"reorder_required"`
}