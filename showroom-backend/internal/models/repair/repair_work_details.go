package repair

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// RepairWorkDetail represents detailed tasks within a work order
type RepairWorkDetail struct {
	WorkDetailID        int        `json:"work_detail_id" db:"work_detail_id"`
	WorkOrderID         int        `json:"work_order_id" db:"work_order_id"`
	DamageID            *int       `json:"damage_id,omitempty" db:"damage_id"`
	TaskSequence        int        `json:"task_sequence" db:"task_sequence"`
	TaskDescription     string     `json:"task_description" db:"task_description"`
	TaskType            string     `json:"task_type" db:"task_type"`
	EstimatedHours      float64    `json:"estimated_hours" db:"estimated_hours"`
	ActualHours         float64    `json:"actual_hours" db:"actual_hours"`
	LaborRate           float64    `json:"labor_rate" db:"labor_rate"`
	TaskStatus          string     `json:"task_status" db:"task_status"`
	StartDate           *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate             *time.Time `json:"end_date,omitempty" db:"end_date"`
	CompletionPercentage int       `json:"completion_percentage" db:"completion_percentage"`
	TaskNotes           *string    `json:"task_notes,omitempty" db:"task_notes"`
	QualityCheckPassed  *bool      `json:"quality_check_passed,omitempty" db:"quality_check_passed"`
	AssignedMechanicID  *int       `json:"assigned_mechanic_id,omitempty" db:"assigned_mechanic_id"`
	VerifiedBy          *int       `json:"verified_by,omitempty" db:"verified_by"`
	VerifiedAt          *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`

	// Related data for joins
	WorkOrderNumber      string `json:"work_order_number,omitempty" db:"work_order_number"`
	DamageDescription    string `json:"damage_description,omitempty" db:"damage_description"`
	AssignedMechanicName string `json:"assigned_mechanic_name,omitempty" db:"assigned_mechanic_name"`
	VerifiedByName       string `json:"verified_by_name,omitempty" db:"verified_by_name"`
}

// RepairWorkDetailListItem represents a simplified work detail for list views
type RepairWorkDetailListItem struct {
	WorkDetailID         int     `json:"work_detail_id" db:"work_detail_id"`
	WorkOrderNumber      string  `json:"work_order_number" db:"work_order_number"`
	TaskSequence         int     `json:"task_sequence" db:"task_sequence"`
	TaskDescription      string  `json:"task_description" db:"task_description"`
	TaskType             string  `json:"task_type" db:"task_type"`
	TaskStatus           string  `json:"task_status" db:"task_status"`
	EstimatedHours       float64 `json:"estimated_hours" db:"estimated_hours"`
	ActualHours          float64 `json:"actual_hours" db:"actual_hours"`
	CompletionPercentage int     `json:"completion_percentage" db:"completion_percentage"`
	AssignedMechanicName string  `json:"assigned_mechanic_name" db:"assigned_mechanic_name"`
	StartDate            *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate              *time.Time `json:"end_date,omitempty" db:"end_date"`
}

// RepairWorkDetailCreateRequest represents a request to create a work detail
type RepairWorkDetailCreateRequest struct {
	WorkOrderID        int     `json:"work_order_id" binding:"required"`
	DamageID           *int    `json:"damage_id,omitempty"`
	TaskSequence       int     `json:"task_sequence" binding:"required,min=1"`
	TaskDescription    string  `json:"task_description" binding:"required"`
	TaskType           string  `json:"task_type" binding:"required,oneof=diagnosis disassembly repair replacement assembly testing quality_check"`
	EstimatedHours     float64 `json:"estimated_hours" binding:"min=0"`
	LaborRate          float64 `json:"labor_rate" binding:"min=0"`
	TaskNotes          *string `json:"task_notes,omitempty"`
	AssignedMechanicID *int    `json:"assigned_mechanic_id,omitempty"`
}

// RepairWorkDetailUpdateRequest represents a request to update a work detail
type RepairWorkDetailUpdateRequest struct {
	TaskSequence         *int     `json:"task_sequence,omitempty" binding:"omitempty,min=1"`
	TaskDescription      *string  `json:"task_description,omitempty"`
	EstimatedHours       *float64 `json:"estimated_hours,omitempty" binding:"omitempty,min=0"`
	ActualHours          *float64 `json:"actual_hours,omitempty" binding:"omitempty,min=0"`
	LaborRate            *float64 `json:"labor_rate,omitempty" binding:"omitempty,min=0"`
	TaskStatus           *string  `json:"task_status,omitempty" binding:"omitempty,oneof=pending in_progress completed cancelled on_hold"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CompletionPercentage *int     `json:"completion_percentage,omitempty" binding:"omitempty,min=0,max=100"`
	TaskNotes            *string  `json:"task_notes,omitempty"`
	QualityCheckPassed   *bool    `json:"quality_check_passed,omitempty"`
	AssignedMechanicID   *int     `json:"assigned_mechanic_id,omitempty"`
}

// RepairWorkDetailFilterParams represents filtering parameters for work detail queries
type RepairWorkDetailFilterParams struct {
	WorkOrderID        *int    `json:"work_order_id,omitempty" form:"work_order_id"`
	DamageID           *int    `json:"damage_id,omitempty" form:"damage_id"`
	TaskType           string  `json:"task_type,omitempty" form:"task_type"`
	TaskStatus         string  `json:"task_status,omitempty" form:"task_status"`
	AssignedMechanicID *int    `json:"assigned_mechanic_id,omitempty" form:"assigned_mechanic_id"`
	QualityCheckPassed *bool   `json:"quality_check_passed,omitempty" form:"quality_check_passed"`
	MinCompletionPct   *int    `json:"min_completion_pct,omitempty" form:"min_completion_pct"`
	MaxCompletionPct   *int    `json:"max_completion_pct,omitempty" form:"max_completion_pct"`
	Search             string  `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// WorkDetailProgressRequest represents a request to update task progress
type WorkDetailProgressRequest struct {
	CompletionPercentage int        `json:"completion_percentage" binding:"required,min=0,max=100"`
	ActualHours          *float64   `json:"actual_hours,omitempty" binding:"omitempty,min=0"`
	TaskStatus           string     `json:"task_status" binding:"required,oneof=pending in_progress completed cancelled on_hold"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	ProgressNotes        *string    `json:"progress_notes,omitempty"`
}

// WorkDetailQualityCheckRequest represents a request to perform quality check
type WorkDetailQualityCheckRequest struct {
	QualityCheckPassed bool    `json:"quality_check_passed" binding:"required"`
	QualityNotes       string  `json:"quality_notes" binding:"required"`
	ActualHours        float64 `json:"actual_hours" binding:"required,min=0"`
}

// WorkDetailAssignmentRequest represents a request to assign task to mechanic
type WorkDetailAssignmentRequest struct {
	AssignedMechanicID *int    `json:"assigned_mechanic_id,omitempty"`
	AssignmentNotes    *string `json:"assignment_notes,omitempty"`
}

// WorkDetailSummary represents a summary of work details for a work order
type WorkDetailSummary struct {
	WorkOrderID          int     `json:"work_order_id" db:"work_order_id"`
	TotalTasks           int     `json:"total_tasks" db:"total_tasks"`
	CompletedTasks       int     `json:"completed_tasks" db:"completed_tasks"`
	InProgressTasks      int     `json:"in_progress_tasks" db:"in_progress_tasks"`
	PendingTasks         int     `json:"pending_tasks" db:"pending_tasks"`
	TotalEstimatedHours  float64 `json:"total_estimated_hours" db:"total_estimated_hours"`
	TotalActualHours     float64 `json:"total_actual_hours" db:"total_actual_hours"`
	OverallProgress      float64 `json:"overall_progress" db:"overall_progress"`
	QualityChecksPassed  int     `json:"quality_checks_passed" db:"quality_checks_passed"`
	QualityChecksFailed  int     `json:"quality_checks_failed" db:"quality_checks_failed"`
}