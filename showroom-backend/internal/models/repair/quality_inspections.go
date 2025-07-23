package repair

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// QualityInspection represents quality inspections for repair work
type QualityInspection struct {
	InspectionID         int        `json:"inspection_id" db:"inspection_id"`
	WorkOrderID          int        `json:"work_order_id" db:"work_order_id"`
	InspectionType       string     `json:"inspection_type" db:"inspection_type"`
	InspectionDate       time.Time  `json:"inspection_date" db:"inspection_date"`
	InspectorID          int        `json:"inspector_id" db:"inspector_id"`
	OverallRating        int        `json:"overall_rating" db:"overall_rating"`
	WorkmanshipRating    int        `json:"workmanship_rating" db:"workmanship_rating"`
	SafetyRating         int        `json:"safety_rating" db:"safety_rating"`
	AppearanceRating     int        `json:"appearance_rating" db:"appearance_rating"`
	FunctionalityRating  int        `json:"functionality_rating" db:"functionality_rating"`
	InspectionStatus     string     `json:"inspection_status" db:"inspection_status"`
	InspectionNotes      *string    `json:"inspection_notes,omitempty" db:"inspection_notes"`
	DefectsFound         *string    `json:"defects_found,omitempty" db:"defects_found"`
	Recommendations      *string    `json:"recommendations,omitempty" db:"recommendations"`
	PhotosJSON           *string    `json:"photos_json,omitempty" db:"photos_json"`
	SignedOffBy          *int       `json:"signed_off_by,omitempty" db:"signed_off_by"`
	SignedOffAt          *time.Time `json:"signed_off_at,omitempty" db:"signed_off_at"`
	ReworkRequired       bool       `json:"rework_required" db:"rework_required"`
	NextInspectionDate   *time.Time `json:"next_inspection_date,omitempty" db:"next_inspection_date"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`

	// Related data for joins
	WorkOrderNumber   string `json:"work_order_number,omitempty" db:"work_order_number"`
	TransactionNumber string `json:"transaction_number,omitempty" db:"transaction_number"`
	VehicleBrand      string `json:"vehicle_brand,omitempty" db:"vehicle_brand"`
	VehicleModel      string `json:"vehicle_model,omitempty" db:"vehicle_model"`
	InspectorName     string `json:"inspector_name,omitempty" db:"inspector_name"`
	SignedOffByName   string `json:"signed_off_by_name,omitempty" db:"signed_off_by_name"`
}

// QualityInspectionListItem represents a simplified inspection for list views
type QualityInspectionListItem struct {
	InspectionID        int       `json:"inspection_id" db:"inspection_id"`
	WorkOrderNumber     string    `json:"work_order_number" db:"work_order_number"`
	TransactionNumber   string    `json:"transaction_number" db:"transaction_number"`
	VehicleBrand        string    `json:"vehicle_brand" db:"vehicle_brand"`
	VehicleModel        string    `json:"vehicle_model" db:"vehicle_model"`
	InspectionType      string    `json:"inspection_type" db:"inspection_type"`
	InspectionStatus    string    `json:"inspection_status" db:"inspection_status"`
	OverallRating       int       `json:"overall_rating" db:"overall_rating"`
	InspectionDate      time.Time `json:"inspection_date" db:"inspection_date"`
	InspectorName       string    `json:"inspector_name" db:"inspector_name"`
	ReworkRequired      bool      `json:"rework_required" db:"rework_required"`
	SignedOffBy         *int      `json:"signed_off_by,omitempty" db:"signed_off_by"`
	SignedOffAt         *time.Time `json:"signed_off_at,omitempty" db:"signed_off_at"`
}

// QualityInspectionCreateRequest represents a request to create a quality inspection
type QualityInspectionCreateRequest struct {
	WorkOrderID         int     `json:"work_order_id" binding:"required"`
	InspectionType      string  `json:"inspection_type" binding:"required,oneof=pre_repair during_repair post_repair final_inspection"`
	OverallRating       int     `json:"overall_rating" binding:"required,min=1,max=10"`
	WorkmanshipRating   int     `json:"workmanship_rating" binding:"required,min=1,max=10"`
	SafetyRating        int     `json:"safety_rating" binding:"required,min=1,max=10"`
	AppearanceRating    int     `json:"appearance_rating" binding:"required,min=1,max=10"`
	FunctionalityRating int     `json:"functionality_rating" binding:"required,min=1,max=10"`
	InspectionNotes     *string `json:"inspection_notes,omitempty"`
	DefectsFound        *string `json:"defects_found,omitempty"`
	Recommendations     *string `json:"recommendations,omitempty"`
	PhotosJSON          *string `json:"photos_json,omitempty"`
	ReworkRequired      bool    `json:"rework_required"`
}

// QualityInspectionUpdateRequest represents a request to update a quality inspection
type QualityInspectionUpdateRequest struct {
	OverallRating        *int    `json:"overall_rating,omitempty" binding:"omitempty,min=1,max=10"`
	WorkmanshipRating    *int    `json:"workmanship_rating,omitempty" binding:"omitempty,min=1,max=10"`
	SafetyRating         *int    `json:"safety_rating,omitempty" binding:"omitempty,min=1,max=10"`
	AppearanceRating     *int    `json:"appearance_rating,omitempty" binding:"omitempty,min=1,max=10"`
	FunctionalityRating  *int    `json:"functionality_rating,omitempty" binding:"omitempty,min=1,max=10"`
	InspectionStatus     *string `json:"inspection_status,omitempty" binding:"omitempty,oneof=passed failed conditional_pass needs_rework"`
	InspectionNotes      *string `json:"inspection_notes,omitempty"`
	DefectsFound         *string `json:"defects_found,omitempty"`
	Recommendations      *string `json:"recommendations,omitempty"`
	PhotosJSON           *string `json:"photos_json,omitempty"`
	ReworkRequired       *bool   `json:"rework_required,omitempty"`
	NextInspectionDate   *time.Time `json:"next_inspection_date,omitempty"`
}

// QualityInspectionFilterParams represents filtering parameters for inspection queries
type QualityInspectionFilterParams struct {
	WorkOrderID      *int       `json:"work_order_id,omitempty" form:"work_order_id"`
	InspectionType   string     `json:"inspection_type,omitempty" form:"inspection_type"`
	InspectionStatus string     `json:"inspection_status,omitempty" form:"inspection_status"`
	InspectorID      *int       `json:"inspector_id,omitempty" form:"inspector_id"`
	MinOverallRating *int       `json:"min_overall_rating,omitempty" form:"min_overall_rating"`
	MaxOverallRating *int       `json:"max_overall_rating,omitempty" form:"max_overall_rating"`
	ReworkRequired   *bool      `json:"rework_required,omitempty" form:"rework_required"`
	SignedOffBy      *int       `json:"signed_off_by,omitempty" form:"signed_off_by"`
	StartDate        *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate          *time.Time `json:"end_date,omitempty" form:"end_date"`
	Search           string     `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// InspectionSignOffRequest represents a request to sign off an inspection
type InspectionSignOffRequest struct {
	Status         string     `json:"status" binding:"required,oneof=passed failed conditional_pass"`
	SignOffNotes   *string    `json:"sign_off_notes,omitempty"`
	ReworkRequired *bool      `json:"rework_required,omitempty"`
	NextInspectionDate *time.Time `json:"next_inspection_date,omitempty"`
}

// InspectionReworkRequest represents a request to schedule rework
type InspectionReworkRequest struct {
	ReworkDescription   string     `json:"rework_description" binding:"required"`
	ReworkPriority     int        `json:"rework_priority" binding:"required,min=1,max=5"`
	EstimatedHours     float64    `json:"estimated_hours" binding:"required,min=0"`
	AssignedMechanicID *int       `json:"assigned_mechanic_id,omitempty"`
	ScheduledDate      *time.Time `json:"scheduled_date,omitempty"`
	ReworkNotes        *string    `json:"rework_notes,omitempty"`
}

// QualityMetrics represents quality metrics for inspections
type QualityMetrics struct {
	WorkOrderID             int     `json:"work_order_id" db:"work_order_id"`
	TotalInspections        int     `json:"total_inspections" db:"total_inspections"`
	PassedInspections       int     `json:"passed_inspections" db:"passed_inspections"`
	FailedInspections       int     `json:"failed_inspections" db:"failed_inspections"`
	ConditionalPasses       int     `json:"conditional_passes" db:"conditional_passes"`
	ReworksRequired         int     `json:"reworks_required" db:"reworks_required"`
	AverageOverallRating    float64 `json:"average_overall_rating" db:"average_overall_rating"`
	AverageWorkmanshipRating float64 `json:"average_workmanship_rating" db:"average_workmanship_rating"`
	AverageSafetyRating     float64 `json:"average_safety_rating" db:"average_safety_rating"`
	AverageAppearanceRating float64 `json:"average_appearance_rating" db:"average_appearance_rating"`
	AverageFunctionalityRating float64 `json:"average_functionality_rating" db:"average_functionality_rating"`
	FirstTimePassRate       float64 `json:"first_time_pass_rate" db:"first_time_pass_rate"`
}

// InspectionScheduleRequest represents a request to schedule an inspection
type InspectionScheduleRequest struct {
	InspectionType     string     `json:"inspection_type" binding:"required,oneof=pre_repair during_repair post_repair final_inspection"`
	ScheduledDate      time.Time  `json:"scheduled_date" binding:"required"`
	InspectorID        int        `json:"inspector_id" binding:"required"`
	SpecialInstructions *string   `json:"special_instructions,omitempty"`
}

// InspectionDashboard represents dashboard data for quality inspections
type InspectionDashboard struct {
	TotalInspections       int     `json:"total_inspections" db:"total_inspections"`
	PendingInspections     int     `json:"pending_inspections" db:"pending_inspections"`
	CompletedToday         int     `json:"completed_today" db:"completed_today"`
	PassRate               float64 `json:"pass_rate" db:"pass_rate"`
	ReworkRate             float64 `json:"rework_rate" db:"rework_rate"`
	AverageRating          float64 `json:"average_rating" db:"average_rating"`
	CriticalDefects        int     `json:"critical_defects" db:"critical_defects"`
	OverdueInspections     int     `json:"overdue_inspections" db:"overdue_inspections"`
}