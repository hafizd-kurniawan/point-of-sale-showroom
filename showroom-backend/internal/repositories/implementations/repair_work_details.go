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
			work_order_id, damage_id, task_description, estimated_hours,
			task_priority, task_category, work_instructions, task_status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING work_detail_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		workDetail.WorkOrderID,
		workDetail.DamageID,
		workDetail.TaskDescription,
		workDetail.EstimatedHours,
		workDetail.TaskPriority,
		workDetail.TaskCategory,
		workDetail.WorkInstructions,
		workDetail.TaskStatus,
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
			rwd.work_detail_id, rwd.work_order_id, rwd.damage_id, rwd.task_description,
			rwd.estimated_hours, rwd.actual_hours, rwd.task_priority, rwd.task_category,
			rwd.work_instructions, rwd.task_status, rwd.assigned_mechanic_id, rwd.progress_percentage,
			rwd.quality_check_status, rwd.completion_notes, rwd.start_time, rwd.end_time,
			rwd.verified_by, rwd.verified_at, rwd.created_at, rwd.updated_at,
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
		&workDetail.TaskDescription,
		&workDetail.EstimatedHours,
		&workDetail.ActualHours,
		&workDetail.TaskPriority,
		&workDetail.TaskCategory,
		&workDetail.WorkInstructions,
		&workDetail.TaskStatus,
		&workDetail.AssignedMechanicID,
		&workDetail.ProgressPercentage,
		&workDetail.QualityCheckStatus,
		&workDetail.CompletionNotes,
		&workDetail.StartTime,
		&workDetail.EndTime,
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
			task_description = $1, estimated_hours = $2, actual_hours = $3,
			task_priority = $4, task_category = $5, work_instructions = $6,
			task_status = $7, progress_percentage = $8, quality_check_status = $9,
			completion_notes = $10, start_time = $11, end_time = $12, updated_at = NOW()
		WHERE work_detail_id = $13
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		workDetail.TaskDescription,
		workDetail.EstimatedHours,
		workDetail.ActualHours,
		workDetail.TaskPriority,
		workDetail.TaskCategory,
		workDetail.WorkInstructions,
		workDetail.TaskStatus,
		workDetail.ProgressPercentage,
		workDetail.QualityCheckStatus,
		workDetail.CompletionNotes,
		workDetail.StartTime,
		workDetail.EndTime,
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
		SET progress_percentage = $1, task_status = $2, completion_notes = $3, updated_at = NOW()
		WHERE work_detail_id = $4`

	_, err := r.db.ExecContext(ctx, query,
		request.ProgressPercentage,
		request.TaskStatus,
		request.CompletionNotes,
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
		SET assigned_mechanic_id = $1, start_time = $2, task_status = 'assigned', updated_at = NOW()
		WHERE work_detail_id = $3`

	_, err := r.db.ExecContext(ctx, query,
		request.MechanicID,
		request.StartTime,
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
		SET quality_check_status = $1, verified_by = $2, verified_at = NOW(), updated_at = NOW()
		WHERE work_detail_id = $3`

	_, err := r.db.ExecContext(ctx, query,
		request.QualityCheckStatus,
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
			COALESCE(AVG(progress_percentage), 0) as average_progress,
			COUNT(CASE WHEN task_status = 'completed' THEN 1 END) as completed_tasks,
			COUNT(CASE WHEN task_status = 'in_progress' THEN 1 END) as in_progress_tasks
		FROM repair_work_details 
		WHERE work_order_id = $1`

	summary := &repair.WorkDetailSummary{}
	err := r.db.QueryRowContext(ctx, query, workOrderID).Scan(
		&summary.TotalTasks,
		&summary.TotalEstimatedHours,
		&summary.TotalActualHours,
		&summary.AverageProgress,
		&summary.CompletedTasks,
		&summary.InProgressTasks,
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