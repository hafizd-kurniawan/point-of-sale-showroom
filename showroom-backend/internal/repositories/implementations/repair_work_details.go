package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/repair"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// RepairWorkDetailRepository implements interfaces.RepairWorkDetailRepository
type RepairWorkDetailRepository struct {
	db *sql.DB
}

// NewRepairWorkDetailRepository creates a new repair work detail repository
func NewRepairWorkDetailRepository(db *sql.DB) interfaces.RepairWorkDetailRepository {
	return &RepairWorkDetailRepository{db: db}
}

// Create creates a new repair work detail
func (r *RepairWorkDetailRepository) Create(ctx context.Context, workDetail *repair.RepairWorkDetail) (*repair.RepairWorkDetail, error) {
	query := `
		INSERT INTO repair_work_details (
			work_order_id, damage_id, task_sequence, task_description, task_type, 
			estimated_hours, labor_rate, task_status, task_notes, assigned_mechanic_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING work_detail_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		workDetail.WorkOrderID,
		workDetail.DamageID,
		workDetail.TaskSequence,
		workDetail.TaskDescription,
		workDetail.TaskType,
		workDetail.EstimatedHours,
		workDetail.LaborRate,
		workDetail.TaskStatus,
		workDetail.TaskNotes,
		workDetail.AssignedMechanicID,
	).Scan(&workDetail.WorkDetailID, &workDetail.CreatedAt, &workDetail.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create repair work detail: %w", err)
	}

	return workDetail, nil
}

// GetByID retrieves a repair work detail by ID
func (r *RepairWorkDetailRepository) GetByID(ctx context.Context, id int) (*repair.RepairWorkDetail, error) {
	query := `
		SELECT 
			rwd.work_detail_id, rwd.work_order_id, rwd.damage_id, rwd.task_sequence,
			rwd.task_description, rwd.task_type, rwd.estimated_hours, rwd.actual_hours,
			rwd.labor_rate, rwd.task_status, rwd.start_date, rwd.end_date,
			rwd.completion_percentage, rwd.task_notes, rwd.quality_check_passed,
			rwd.assigned_mechanic_id, rwd.verified_by, rwd.verified_at, 
			rwd.created_at, rwd.updated_at,
			rwo.work_order_number,
			vd.damage_description,
			u1.full_name as assigned_mechanic_name,
			u2.full_name as verified_by_name
		FROM repair_work_details rwd
		LEFT JOIN repair_work_orders rwo ON rwd.work_order_id = rwo.work_order_id
		LEFT JOIN vehicle_damages vd ON rwd.damage_id = vd.damage_id
		LEFT JOIN users u1 ON rwd.assigned_mechanic_id = u1.user_id
		LEFT JOIN users u2 ON rwd.verified_by = u2.user_id
		WHERE rwd.work_detail_id = $1`

	workDetail := &repair.RepairWorkDetail{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&workDetail.WorkDetailID,
		&workDetail.WorkOrderID,
		&workDetail.DamageID,
		&workDetail.TaskSequence,
		&workDetail.TaskDescription,
		&workDetail.TaskType,
		&workDetail.EstimatedHours,
		&workDetail.ActualHours,
		&workDetail.LaborRate,
		&workDetail.TaskStatus,
		&workDetail.StartDate,
		&workDetail.EndDate,
		&workDetail.CompletionPercentage,
		&workDetail.TaskNotes,
		&workDetail.QualityCheckPassed,
		&workDetail.AssignedMechanicID,
		&workDetail.VerifiedBy,
		&workDetail.VerifiedAt,
		&workDetail.CreatedAt,
		&workDetail.UpdatedAt,
		&workDetail.WorkOrderNumber,
		&workDetail.DamageDescription,
		&workDetail.AssignedMechanicName,
		&workDetail.VerifiedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("repair work detail not found")
		}
		return nil, fmt.Errorf("failed to get repair work detail: %w", err)
	}

	return workDetail, nil
}

// Update updates a repair work detail
func (r *RepairWorkDetailRepository) Update(ctx context.Context, id int, workDetail *repair.RepairWorkDetail) (*repair.RepairWorkDetail, error) {
	query := `
		UPDATE repair_work_details SET
			task_sequence = $1, task_description = $2, task_type = $3, 
			estimated_hours = $4, actual_hours = $5, labor_rate = $6,
			task_status = $7, start_date = $8, end_date = $9,
			completion_percentage = $10, task_notes = $11, quality_check_passed = $12,
			assigned_mechanic_id = $13, updated_at = NOW()
		WHERE work_detail_id = $14
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		workDetail.TaskSequence,
		workDetail.TaskDescription,
		workDetail.TaskType,
		workDetail.EstimatedHours,
		workDetail.ActualHours,
		workDetail.LaborRate,
		workDetail.TaskStatus,
		workDetail.StartDate,
		workDetail.EndDate,
		workDetail.CompletionPercentage,
		workDetail.TaskNotes,
		workDetail.QualityCheckPassed,
		workDetail.AssignedMechanicID,
		id,
	).Scan(&workDetail.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update repair work detail: %w", err)
	}

	workDetail.WorkDetailID = id
	return workDetail, nil
}

// UpdateProgress updates the progress of a work detail
func (r *RepairWorkDetailRepository) UpdateProgress(ctx context.Context, id int, request *repair.WorkDetailProgressRequest) error {
	query := `
		UPDATE repair_work_details 
		SET completion_percentage = $1, task_status = $2, actual_hours = $3, 
			start_date = $4, end_date = $5, task_notes = $6, updated_at = NOW()
		WHERE work_detail_id = $7`

	_, err := r.db.ExecContext(ctx, query,
		request.CompletionPercentage,
		request.TaskStatus,
		request.ActualHours,
		request.StartDate,
		request.EndDate,
		request.ProgressNotes,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	return nil
}

// UpdateStatus updates the status of a work detail
func (r *RepairWorkDetailRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `
		UPDATE repair_work_details 
		SET task_status = $1, updated_at = NOW()
		WHERE work_detail_id = $2`

	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update work detail status: %w", err)
	}

	return nil
}

// Delete deletes a repair work detail
func (r *RepairWorkDetailRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM repair_work_details WHERE work_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete repair work detail: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("repair work detail not found")
	}

	return nil
}

// AssignMechanic assigns a mechanic to a work detail
func (r *RepairWorkDetailRepository) AssignMechanic(ctx context.Context, id int, request *repair.WorkDetailAssignmentRequest) error {
	query := `
		UPDATE repair_work_details 
		SET assigned_mechanic_id = $1, task_status = 'in_progress', updated_at = NOW()
		WHERE work_detail_id = $2`

	_, err := r.db.ExecContext(ctx, query,
		request.AssignedMechanicID,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to assign mechanic: %w", err)
	}

	return nil
}

// PerformQualityCheck performs quality check on work detail
func (r *RepairWorkDetailRepository) PerformQualityCheck(ctx context.Context, id int, request *repair.WorkDetailQualityCheckRequest, verifiedBy int) error {
	query := `
		UPDATE repair_work_details 
		SET quality_check_passed = $1, actual_hours = $2, task_notes = $3, 
			verified_by = $4, verified_at = NOW(), updated_at = NOW()
		WHERE work_detail_id = $5`

	_, err := r.db.ExecContext(ctx, query,
		request.QualityCheckPassed,
		request.ActualHours,
		request.QualityNotes,
		verifiedBy,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to perform quality check: %w", err)
	}

	return nil
}

// GetWorkDetailSummary gets work detail summary for a work order
func (r *RepairWorkDetailRepository) GetWorkDetailSummary(ctx context.Context, workOrderID int) (*repair.WorkDetailSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_tasks,
			COALESCE(SUM(estimated_hours), 0) as total_estimated_hours,
			COALESCE(SUM(actual_hours), 0) as total_actual_hours,
			COALESCE(AVG(completion_percentage), 0) as overall_progress,
			COUNT(CASE WHEN task_status = 'completed' THEN 1 END) as completed_tasks,
			COUNT(CASE WHEN task_status = 'in_progress' THEN 1 END) as in_progress_tasks,
			COUNT(CASE WHEN task_status = 'pending' THEN 1 END) as pending_tasks,
			COUNT(CASE WHEN quality_check_passed = true THEN 1 END) as quality_checks_passed,
			COUNT(CASE WHEN quality_check_passed = false THEN 1 END) as quality_checks_failed
		FROM repair_work_details 
		WHERE work_order_id = $1`

	summary := &repair.WorkDetailSummary{}
	err := r.db.QueryRowContext(ctx, query, workOrderID).Scan(
		&summary.TotalTasks,
		&summary.TotalEstimatedHours,
		&summary.TotalActualHours,
		&summary.OverallProgress,
		&summary.CompletedTasks,
		&summary.InProgressTasks,
		&summary.PendingTasks,
		&summary.QualityChecksPassed,
		&summary.QualityChecksFailed,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get work detail summary: %w", err)
	}

	summary.WorkOrderID = workOrderID

	return summary, nil
}

// Stub implementations for list methods
func (r *RepairWorkDetailRepository) List(ctx context.Context, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkDetailRepository) GetByWorkOrderID(ctx context.Context, workOrderID int, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkDetailRepository) GetByDamageID(ctx context.Context, damageID int, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkDetailRepository) GetByMechanic(ctx context.Context, mechanicID int, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkDetailRepository) GetByStatus(ctx context.Context, status string, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkDetailRepository) Search(ctx context.Context, query string, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}