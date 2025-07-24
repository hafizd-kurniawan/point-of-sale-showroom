package repair

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// RepairWorkOrder represents a work order for vehicle repairs
type RepairWorkOrder struct {
	WorkOrderID          int        `json:"work_order_id" db:"work_order_id"`
	WorkOrderNumber      string     `json:"work_order_number" db:"work_order_number"`
	TransactionID        int        `json:"transaction_id" db:"transaction_id"`
	WorkOrderType        string     `json:"work_order_type" db:"work_order_type"`
	WorkOrderPriority    int        `json:"work_order_priority" db:"work_order_priority"`
	ScheduledStartDate   *time.Time `json:"scheduled_start_date,omitempty" db:"scheduled_start_date"`
	ScheduledEndDate     *time.Time `json:"scheduled_end_date,omitempty" db:"scheduled_end_date"`
	ActualStartDate      *time.Time `json:"actual_start_date,omitempty" db:"actual_start_date"`
	ActualEndDate        *time.Time `json:"actual_end_date,omitempty" db:"actual_end_date"`
	EstimatedCost        float64    `json:"estimated_cost" db:"estimated_cost"`
	ActualCost           float64    `json:"actual_cost" db:"actual_cost"`
	LaborHoursEstimated  float64    `json:"labor_hours_estimated" db:"labor_hours_estimated"`
	LaborHoursActual     float64    `json:"labor_hours_actual" db:"labor_hours_actual"`
	WorkOrderStatus      string     `json:"work_order_status" db:"work_order_status"`
	WorkDescription      string     `json:"work_description" db:"work_description"`
	SpecialInstructions  *string    `json:"special_instructions,omitempty" db:"special_instructions"`
	CompletionNotes      *string    `json:"completion_notes,omitempty" db:"completion_notes"`
	AssignedMechanicID   *int       `json:"assigned_mechanic_id,omitempty" db:"assigned_mechanic_id"`
	SupervisorID         *int       `json:"supervisor_id,omitempty" db:"supervisor_id"`
	CreatedBy            int        `json:"created_by" db:"created_by"`
	ApprovedBy           *int       `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt           *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`

	// Related data for joins
	TransactionNumber    string `json:"transaction_number,omitempty" db:"transaction_number"`
	VehicleBrand         string `json:"vehicle_brand,omitempty" db:"vehicle_brand"`
	VehicleModel         string `json:"vehicle_model,omitempty" db:"vehicle_model"`
	AssignedMechanicName string `json:"assigned_mechanic_name,omitempty" db:"assigned_mechanic_name"`
	SupervisorName       string `json:"supervisor_name,omitempty" db:"supervisor_name"`
	CreatedByName        string `json:"created_by_name,omitempty" db:"created_by_name"`
	ApprovedByName       string `json:"approved_by_name,omitempty" db:"approved_by_name"`
}

// RepairWorkOrderListItem represents a simplified work order for list views
type RepairWorkOrderListItem struct {
	WorkOrderID          int        `json:"work_order_id" db:"work_order_id"`
	WorkOrderNumber      string     `json:"work_order_number" db:"work_order_number"`
	TransactionNumber    string     `json:"transaction_number" db:"transaction_number"`
	VehicleBrand         string     `json:"vehicle_brand" db:"vehicle_brand"`
	VehicleModel         string     `json:"vehicle_model" db:"vehicle_model"`
	WorkOrderType        string     `json:"work_order_type" db:"work_order_type"`
	WorkOrderPriority    int        `json:"work_order_priority" db:"work_order_priority"`
	WorkOrderStatus      string     `json:"work_order_status" db:"work_order_status"`
	EstimatedCost        float64    `json:"estimated_cost" db:"estimated_cost"`
	AssignedMechanicName string     `json:"assigned_mechanic_name" db:"assigned_mechanic_name"`
	ScheduledStartDate   *time.Time `json:"scheduled_start_date,omitempty" db:"scheduled_start_date"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
}

// RepairWorkOrderCreateRequest represents a request to create a work order
type RepairWorkOrderCreateRequest struct {
	TransactionID        int        `json:"transaction_id" binding:"required"`
	WorkOrderType        string     `json:"work_order_type" binding:"required,oneof=inspection repair maintenance improvement"`
	WorkOrderPriority    int        `json:"work_order_priority" binding:"required,min=1,max=5"`
	ScheduledStartDate   *time.Time `json:"scheduled_start_date,omitempty"`
	ScheduledEndDate     *time.Time `json:"scheduled_end_date,omitempty"`
	EstimatedCost        float64    `json:"estimated_cost" binding:"min=0"`
	LaborHoursEstimated  float64    `json:"labor_hours_estimated" binding:"min=0"`
	WorkDescription      string     `json:"work_description" binding:"required"`
	SpecialInstructions  *string    `json:"special_instructions,omitempty"`
	AssignedMechanicID   *int       `json:"assigned_mechanic_id,omitempty"`
	SupervisorID         *int       `json:"supervisor_id,omitempty"`
}

// RepairWorkOrderUpdateRequest represents a request to update a work order
type RepairWorkOrderUpdateRequest struct {
	WorkOrderPriority    *int       `json:"work_order_priority,omitempty" binding:"omitempty,min=1,max=5"`
	ScheduledStartDate   *time.Time `json:"scheduled_start_date,omitempty"`
	ScheduledEndDate     *time.Time `json:"scheduled_end_date,omitempty"`
	EstimatedCost        *float64   `json:"estimated_cost,omitempty" binding:"omitempty,min=0"`
	ActualCost           *float64   `json:"actual_cost,omitempty" binding:"omitempty,min=0"`
	LaborHoursEstimated  *float64   `json:"labor_hours_estimated,omitempty" binding:"omitempty,min=0"`
	LaborHoursActual     *float64   `json:"labor_hours_actual,omitempty" binding:"omitempty,min=0"`
	WorkOrderStatus      *string    `json:"work_order_status,omitempty" binding:"omitempty,oneof=draft scheduled in_progress suspended completed cancelled"`
	WorkDescription      *string    `json:"work_description,omitempty"`
	SpecialInstructions  *string    `json:"special_instructions,omitempty"`
	CompletionNotes      *string    `json:"completion_notes,omitempty"`
	AssignedMechanicID   *int       `json:"assigned_mechanic_id,omitempty"`
	SupervisorID         *int       `json:"supervisor_id,omitempty"`
}

// RepairWorkOrderFilterParams represents filtering parameters for work order queries
type RepairWorkOrderFilterParams struct {
	TransactionID      *int       `json:"transaction_id,omitempty" form:"transaction_id"`
	WorkOrderType      string     `json:"work_order_type,omitempty" form:"work_order_type"`
	WorkOrderStatus    string     `json:"work_order_status,omitempty" form:"work_order_status"`
	WorkOrderPriority  *int       `json:"work_order_priority,omitempty" form:"work_order_priority"`
	AssignedMechanicID *int       `json:"assigned_mechanic_id,omitempty" form:"assigned_mechanic_id"`
	SupervisorID       *int       `json:"supervisor_id,omitempty" form:"supervisor_id"`
	CreatedBy          *int       `json:"created_by,omitempty" form:"created_by"`
	StartDate          *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate            *time.Time `json:"end_date,omitempty" form:"end_date"`
	Search             string     `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// WorkOrderStatusRequest represents a request to update work order status
type WorkOrderStatusRequest struct {
	Status          string     `json:"status" binding:"required,oneof=draft scheduled in_progress suspended completed cancelled"`
	StatusNotes     *string    `json:"status_notes,omitempty"`
	ActualStartDate *time.Time `json:"actual_start_date,omitempty"`
	ActualEndDate   *time.Time `json:"actual_end_date,omitempty"`
}

// WorkOrderAssignmentRequest represents a request to assign work order to mechanic
type WorkOrderAssignmentRequest struct {
	AssignedMechanicID *int    `json:"assigned_mechanic_id,omitempty"`
	SupervisorID       *int    `json:"supervisor_id,omitempty"`
	AssignmentNotes    *string `json:"assignment_notes,omitempty"`
}

// WorkOrderApprovalRequest represents a request to approve work order
type WorkOrderApprovalRequest struct {
	Status        string  `json:"status" binding:"required,oneof=approved rejected"`
	ApprovalNotes *string `json:"approval_notes,omitempty"`
}

// WorkOrderSummary represents a summary of work orders for a transaction
type WorkOrderSummary struct {
	TransactionID       int     `json:"transaction_id" db:"transaction_id"`
	TotalWorkOrders     int     `json:"total_work_orders" db:"total_work_orders"`
	CompletedOrders     int     `json:"completed_orders" db:"completed_orders"`
	InProgressOrders    int     `json:"in_progress_orders" db:"in_progress_orders"`
	ScheduledOrders     int     `json:"scheduled_orders" db:"scheduled_orders"`
	TotalEstimatedCost  float64 `json:"total_estimated_cost" db:"total_estimated_cost"`
	TotalActualCost     float64 `json:"total_actual_cost" db:"total_actual_cost"`
	TotalEstimatedHours float64 `json:"total_estimated_hours" db:"total_estimated_hours"`
	TotalActualHours    float64 `json:"total_actual_hours" db:"total_actual_hours"`
}