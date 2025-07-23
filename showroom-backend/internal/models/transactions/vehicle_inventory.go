package transactions

import (
	"time"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// VehicleStatus represents the status of a vehicle in inventory
type VehicleStatus string

const (
	VehicleStatusPurchased              VehicleStatus = "purchased"
	VehicleStatusPendingRepairApproval  VehicleStatus = "pending_repair_approval"
	VehicleStatusApprovedRepair         VehicleStatus = "approved_repair"
	VehicleStatusInRepair               VehicleStatus = "in_repair"
	VehicleStatusReadyToSell            VehicleStatus = "ready_to_sell"
	VehicleStatusReserved               VehicleStatus = "reserved"
	VehicleStatusSold                   VehicleStatus = "sold"
	VehicleStatusScrapped               VehicleStatus = "scrapped"
)

// PurchaseType represents the type of vehicle purchase
type PurchaseType string

const (
	PurchaseTypeDirectSell PurchaseType = "direct_sell"
	PurchaseTypeNeedRepair PurchaseType = "need_repair"
	PurchaseTypeAuction    PurchaseType = "auction"
	PurchaseTypeTradeIn    PurchaseType = "trade_in"
)

// ConditionGrade represents the condition grade of a vehicle
type ConditionGrade string

const (
	ConditionGradeA ConditionGrade = "A"
	ConditionGradeB ConditionGrade = "B"
	ConditionGradeC ConditionGrade = "C"
	ConditionGradeD ConditionGrade = "D"
	ConditionGradeE ConditionGrade = "E"
)

// VehicleInventory represents a vehicle in the inventory
type VehicleInventory struct {
	VehicleID                int             `json:"vehicle_id" db:"vehicle_id"`
	VehicleCode              string          `json:"vehicle_code" db:"vehicle_code"`
	ChassisNumber            string          `json:"chassis_number" db:"chassis_number"`
	EngineNumber             string          `json:"engine_number" db:"engine_number"`
	LicensePlate             *string         `json:"license_plate,omitempty" db:"license_plate"`
	BrandID                  int             `json:"brand_id" db:"brand_id"`
	CategoryID               int             `json:"category_id" db:"category_id"`
	ModelID                  int             `json:"model_id" db:"model_id"`
	ModelVariant             *string         `json:"model_variant,omitempty" db:"model_variant"`
	Year                     int             `json:"year" db:"year"`
	Color                    string          `json:"color" db:"color"`
	Mileage                  int             `json:"mileage" db:"mileage"`
	FuelType                 string          `json:"fuel_type" db:"fuel_type"`
	Transmission             string          `json:"transmission" db:"transmission"`
	EngineCapacity           *int            `json:"engine_capacity,omitempty" db:"engine_capacity"`
	PurchasePrice            float64         `json:"purchase_price" db:"purchase_price"`
	EstimatedSellingPrice    *float64        `json:"estimated_selling_price,omitempty" db:"estimated_selling_price"`
	FinalSellingPrice        *float64        `json:"final_selling_price,omitempty" db:"final_selling_price"`
	PurchaseType             PurchaseType    `json:"purchase_type" db:"purchase_type"`
	ConditionGrade           ConditionGrade  `json:"condition_grade" db:"condition_grade"`
	Status                   VehicleStatus   `json:"status" db:"status"`
	PurchaseDate             time.Time       `json:"purchase_date" db:"purchase_date"`
	ReadyToSellDate          *time.Time      `json:"ready_to_sell_date,omitempty" db:"ready_to_sell_date"`
	SoldDate                 *time.Time      `json:"sold_date,omitempty" db:"sold_date"`
	CreatedAt                time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
	PurchasedFromCustomerID  *int            `json:"purchased_from_customer_id,omitempty" db:"purchased_from_customer_id"`
	CreatedBy                int             `json:"created_by" db:"created_by"`
	VehicleImagesJSON        *string         `json:"vehicle_images_json,omitempty" db:"vehicle_images_json"`
	PurchaseNotes            *string         `json:"purchase_notes,omitempty" db:"purchase_notes"`
	ConditionNotes           *string         `json:"condition_notes,omitempty" db:"condition_notes"`
	HasCompleteDocuments     bool            `json:"has_complete_documents" db:"has_complete_documents"`
	DocumentsJSON            *string         `json:"documents_json,omitempty" db:"documents_json"`

	// Related data for joins
	BrandName           string `json:"brand_name,omitempty" db:"brand_name"`
	CategoryName        string `json:"category_name,omitempty" db:"category_name"`
	ModelName           string `json:"model_name,omitempty" db:"model_name"`
	CustomerName        string `json:"customer_name,omitempty" db:"customer_name"`
	CreatedByName       string `json:"created_by_name,omitempty" db:"created_by_name"`
}

// VehicleInventoryListItem represents a simplified vehicle for list views
type VehicleInventoryListItem struct {
	VehicleID             int           `json:"vehicle_id" db:"vehicle_id"`
	VehicleCode           string        `json:"vehicle_code" db:"vehicle_code"`
	ChassisNumber         string        `json:"chassis_number" db:"chassis_number"`
	LicensePlate          *string       `json:"license_plate,omitempty" db:"license_plate"`
	BrandName             string        `json:"brand_name" db:"brand_name"`
	ModelName             string        `json:"model_name" db:"model_name"`
	Year                  int           `json:"year" db:"year"`
	Color                 string        `json:"color" db:"color"`
	Mileage               int           `json:"mileage" db:"mileage"`
	PurchasePrice         float64       `json:"purchase_price" db:"purchase_price"`
	EstimatedSellingPrice *float64      `json:"estimated_selling_price,omitempty" db:"estimated_selling_price"`
	ConditionGrade        ConditionGrade `json:"condition_grade" db:"condition_grade"`
	Status                VehicleStatus `json:"status" db:"status"`
	PurchaseDate          time.Time     `json:"purchase_date" db:"purchase_date"`
	CreatedAt             time.Time     `json:"created_at" db:"created_at"`
}

// VehicleInventoryCreateRequest represents a request to add a vehicle to inventory
type VehicleInventoryCreateRequest struct {
	ChassisNumber            string          `json:"chassis_number" binding:"required,max=50"`
	EngineNumber             string          `json:"engine_number" binding:"required,max=50"`
	LicensePlate             *string         `json:"license_plate,omitempty" binding:"omitempty,max=20"`
	BrandID                  int             `json:"brand_id" binding:"required,min=1"`
	CategoryID               int             `json:"category_id" binding:"required,min=1"`
	ModelID                  int             `json:"model_id" binding:"required,min=1"`
	ModelVariant             *string         `json:"model_variant,omitempty" binding:"omitempty,max=100"`
	Year                     int             `json:"year" binding:"required,min=1900,max=2100"`
	Color                    string          `json:"color" binding:"required,max=50"`
	Mileage                  int             `json:"mileage" binding:"min=0"`
	FuelType                 string          `json:"fuel_type" binding:"required,max=50"`
	Transmission             string          `json:"transmission" binding:"required,max=50"`
	EngineCapacity           *int            `json:"engine_capacity,omitempty" binding:"omitempty,min=1"`
	PurchasePrice            float64         `json:"purchase_price" binding:"required,min=0"`
	EstimatedSellingPrice    *float64        `json:"estimated_selling_price,omitempty" binding:"omitempty,min=0"`
	PurchaseType             PurchaseType    `json:"purchase_type" binding:"required"`
	ConditionGrade           ConditionGrade  `json:"condition_grade" binding:"required"`
	PurchaseDate             time.Time       `json:"purchase_date" binding:"required"`
	PurchasedFromCustomerID  *int            `json:"purchased_from_customer_id,omitempty" binding:"omitempty,min=1"`
	PurchaseNotes            *string         `json:"purchase_notes,omitempty"`
	ConditionNotes           *string         `json:"condition_notes,omitempty"`
	HasCompleteDocuments     bool            `json:"has_complete_documents"`
}

// VehicleInventoryUpdateRequest represents a request to update a vehicle in inventory
type VehicleInventoryUpdateRequest struct {
	LicensePlate             *string         `json:"license_plate,omitempty" binding:"omitempty,max=20"`
	ModelVariant             *string         `json:"model_variant,omitempty" binding:"omitempty,max=100"`
	Color                    *string         `json:"color,omitempty" binding:"omitempty,max=50"`
	Mileage                  *int            `json:"mileage,omitempty" binding:"omitempty,min=0"`
	EstimatedSellingPrice    *float64        `json:"estimated_selling_price,omitempty" binding:"omitempty,min=0"`
	FinalSellingPrice        *float64        `json:"final_selling_price,omitempty" binding:"omitempty,min=0"`
	ConditionGrade           *ConditionGrade `json:"condition_grade,omitempty"`
	Status                   *VehicleStatus  `json:"status,omitempty"`
	ReadyToSellDate          *time.Time      `json:"ready_to_sell_date,omitempty"`
	PurchaseNotes            *string         `json:"purchase_notes,omitempty"`
	ConditionNotes           *string         `json:"condition_notes,omitempty"`
	HasCompleteDocuments     *bool           `json:"has_complete_documents,omitempty"`
}

// VehicleInventoryFilterParams represents filtering parameters for vehicle inventory queries
type VehicleInventoryFilterParams struct {
	BrandID       *int           `json:"brand_id,omitempty" form:"brand_id"`
	CategoryID    *int           `json:"category_id,omitempty" form:"category_id"`
	ModelID       *int           `json:"model_id,omitempty" form:"model_id"`
	Status        *VehicleStatus `json:"status,omitempty" form:"status"`
	PurchaseType  *PurchaseType  `json:"purchase_type,omitempty" form:"purchase_type"`
	ConditionGrade *ConditionGrade `json:"condition_grade,omitempty" form:"condition_grade"`
	MinYear       *int           `json:"min_year,omitempty" form:"min_year"`
	MaxYear       *int           `json:"max_year,omitempty" form:"max_year"`
	MinPrice      *float64       `json:"min_price,omitempty" form:"min_price"`
	MaxPrice      *float64       `json:"max_price,omitempty" form:"max_price"`
	Search        string         `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// IsValid validates the vehicle status
func (s VehicleStatus) IsValid() bool {
	switch s {
	case VehicleStatusPurchased, VehicleStatusPendingRepairApproval, VehicleStatusApprovedRepair,
		 VehicleStatusInRepair, VehicleStatusReadyToSell, VehicleStatusReserved, 
		 VehicleStatusSold, VehicleStatusScrapped:
		return true
	default:
		return false
	}
}

// IsValid validates the purchase type
func (p PurchaseType) IsValid() bool {
	switch p {
	case PurchaseTypeDirectSell, PurchaseTypeNeedRepair, PurchaseTypeAuction, PurchaseTypeTradeIn:
		return true
	default:
		return false
	}
}

// IsValid validates the condition grade
func (c ConditionGrade) IsValid() bool {
	switch c {
	case ConditionGradeA, ConditionGradeB, ConditionGradeC, ConditionGradeD, ConditionGradeE:
		return true
	default:
		return false
	}
}

// CanBeDeleted checks if the vehicle can be deleted (soft delete)
func (v *VehicleInventory) CanBeDeleted() bool {
	return v.Status != VehicleStatusSold
}

// IsSold checks if the vehicle has been sold
func (v *VehicleInventory) IsSold() bool {
	return v.Status == VehicleStatusSold
}

// IsReadyToSell checks if the vehicle is ready to be sold
func (v *VehicleInventory) IsReadyToSell() bool {
	return v.Status == VehicleStatusReadyToSell
}

// NeedsRepair checks if the vehicle needs repair
func (v *VehicleInventory) NeedsRepair() bool {
	return v.PurchaseType == PurchaseTypeNeedRepair && 
		   (v.Status == VehicleStatusPurchased || v.Status == VehicleStatusPendingRepairApproval)
}