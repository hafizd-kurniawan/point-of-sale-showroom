package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/repair"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// QualityInspectionRepository implements interfaces.QualityInspectionRepository
type QualityInspectionRepository struct {
	db *sql.DB
}

// NewQualityInspectionRepository creates a new quality inspection repository
func NewQualityInspectionRepository(db *sql.DB) interfaces.QualityInspectionRepository {
	return &QualityInspectionRepository{db: db}
}

// Create creates a new quality inspection
func (r *QualityInspectionRepository) Create(ctx context.Context, inspection *repair.QualityInspection) (*repair.QualityInspection, error) {
	query := `
		INSERT INTO quality_inspections (
			work_order_id, inspection_type, inspection_checklist_json, 
			quality_standards_json, inspection_status, scheduled_date, inspector_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING inspection_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		inspection.WorkOrderID,
		inspection.InspectionType,
		inspection.InspectionChecklistJSON,
		inspection.QualityStandardsJSON,
		inspection.InspectionStatus,
		inspection.ScheduledDate,
		inspection.InspectorID,
	).Scan(&inspection.InspectionID, &inspection.CreatedAt, &inspection.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create quality inspection: %w", err)
	}

	return inspection, nil
}

// GetByID retrieves a quality inspection by ID
func (r *QualityInspectionRepository) GetByID(ctx context.Context, id int) (*repair.QualityInspection, error) {
	query := `
		SELECT 
			qi.inspection_id, qi.work_order_id, qi.inspection_type, qi.inspection_checklist_json,
			qi.quality_standards_json, qi.inspection_status, qi.quality_score, qi.pass_fail_status,
			qi.inspection_notes, qi.defects_found_json, qi.rework_required, qi.rework_instructions,
			qi.scheduled_date, qi.inspection_date, qi.completion_date, qi.inspector_id,
			qi.signed_off_by, qi.signed_off_at, qi.created_at, qi.updated_at,
			rwo.work_order_number,
			u1.full_name as inspector_name,
			u2.full_name as signed_off_by_name
		FROM quality_inspections qi
		LEFT JOIN repair_work_orders rwo ON qi.work_order_id = rwo.work_order_id
		LEFT JOIN users u1 ON qi.inspector_id = u1.user_id
		LEFT JOIN users u2 ON qi.signed_off_by = u2.user_id
		WHERE qi.inspection_id = $1`

	inspection := &repair.QualityInspection{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inspection.InspectionID,
		&inspection.WorkOrderID,
		&inspection.InspectionType,
		&inspection.InspectionChecklistJSON,
		&inspection.QualityStandardsJSON,
		&inspection.InspectionStatus,
		&inspection.QualityScore,
		&inspection.PassFailStatus,
		&inspection.InspectionNotes,
		&inspection.DefectsFoundJSON,
		&inspection.ReworkRequired,
		&inspection.ReworkInstructions,
		&inspection.ScheduledDate,
		&inspection.InspectionDate,
		&inspection.CompletionDate,
		&inspection.InspectorID,
		&inspection.SignedOffBy,
		&inspection.SignedOffAt,
		&inspection.CreatedAt,
		&inspection.UpdatedAt,
		&inspection.WorkOrderNumber,
		&inspection.InspectorName,
		&inspection.SignedOffByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("quality inspection not found")
		}
		return nil, fmt.Errorf("failed to get quality inspection: %w", err)
	}

	return inspection, nil
}

// Update updates a quality inspection
func (r *QualityInspectionRepository) Update(ctx context.Context, id int, inspection *repair.QualityInspection) (*repair.QualityInspection, error) {
	query := `
		UPDATE quality_inspections SET
			inspection_type = $1, inspection_checklist_json = $2, quality_standards_json = $3,
			inspection_status = $4, quality_score = $5, pass_fail_status = $6,
			inspection_notes = $7, defects_found_json = $8, rework_required = $9,
			rework_instructions = $10, scheduled_date = $11, inspection_date = $12,
			completion_date = $13, updated_at = NOW()
		WHERE inspection_id = $14
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		inspection.InspectionType,
		inspection.InspectionChecklistJSON,
		inspection.QualityStandardsJSON,
		inspection.InspectionStatus,
		inspection.QualityScore,
		inspection.PassFailStatus,
		inspection.InspectionNotes,
		inspection.DefectsFoundJSON,
		inspection.ReworkRequired,
		inspection.ReworkInstructions,
		inspection.ScheduledDate,
		inspection.InspectionDate,
		inspection.CompletionDate,
		id,
	).Scan(&inspection.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update quality inspection: %w", err)
	}

	inspection.InspectionID = id
	return inspection, nil
}

// Delete deletes a quality inspection
func (r *QualityInspectionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM quality_inspections WHERE inspection_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete quality inspection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quality inspection not found")
	}

	return nil
}

// SignOffInspection signs off an inspection
func (r *QualityInspectionRepository) SignOffInspection(ctx context.Context, id int, request *repair.InspectionSignOffRequest, signedOffBy int) error {
	query := `
		UPDATE quality_inspections 
		SET inspection_status = $1, pass_fail_status = $2, quality_score = $3,
			inspection_notes = $4, signed_off_by = $5, signed_off_at = NOW(),
			completion_date = NOW(), updated_at = NOW()
		WHERE inspection_id = $6`

	_, err := r.db.ExecContext(ctx, query,
		request.InspectionStatus,
		request.PassFailStatus,
		request.QualityScore,
		request.InspectionNotes,
		signedOffBy,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to sign off inspection: %w", err)
	}

	return nil
}

// ScheduleRework schedules rework for failed inspection
func (r *QualityInspectionRepository) ScheduleRework(ctx context.Context, id int, request *repair.InspectionReworkRequest) error {
	query := `
		UPDATE quality_inspections 
		SET rework_required = true, rework_instructions = $1,
			inspection_status = 'rework_required', updated_at = NOW()
		WHERE inspection_id = $2`

	_, err := r.db.ExecContext(ctx, query, request.ReworkInstructions, id)
	if err != nil {
		return fmt.Errorf("failed to schedule rework: %w", err)
	}

	return nil
}

// ScheduleInspection schedules a new inspection
func (r *QualityInspectionRepository) ScheduleInspection(ctx context.Context, request *repair.InspectionScheduleRequest) (*repair.QualityInspection, error) {
	query := `
		INSERT INTO quality_inspections (
			work_order_id, inspection_type, inspection_checklist_json,
			quality_standards_json, inspection_status, scheduled_date, inspector_id
		) VALUES ($1, $2, $3, $4, 'scheduled', $5, $6)
		RETURNING inspection_id, created_at, updated_at`

	inspection := &repair.QualityInspection{
		WorkOrderID:             request.WorkOrderID,
		InspectionType:         request.InspectionType,
		InspectionChecklistJSON: request.InspectionChecklistJSON,
		QualityStandardsJSON:   request.QualityStandardsJSON,
		ScheduledDate:          request.ScheduledDate,
		InspectorID:            request.InspectorID,
	}

	err := r.db.QueryRowContext(ctx, query,
		request.WorkOrderID,
		request.InspectionType,
		request.InspectionChecklistJSON,
		request.QualityStandardsJSON,
		request.ScheduledDate,
		request.InspectorID,
	).Scan(&inspection.InspectionID, &inspection.CreatedAt, &inspection.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to schedule inspection: %w", err)
	}

	return inspection, nil
}

// GetQualityMetrics gets quality metrics for a work order
func (r *QualityInspectionRepository) GetQualityMetrics(ctx context.Context, workOrderID int) (*repair.QualityMetrics, error) {
	query := `
		SELECT 
			COUNT(*) as total_inspections,
			COALESCE(AVG(quality_score), 0) as average_quality_score,
			COUNT(CASE WHEN pass_fail_status = 'pass' THEN 1 END) as passed_inspections,
			COUNT(CASE WHEN pass_fail_status = 'fail' THEN 1 END) as failed_inspections,
			COUNT(CASE WHEN rework_required = true THEN 1 END) as rework_required_count
		FROM quality_inspections 
		WHERE work_order_id = $1`

	metrics := &repair.QualityMetrics{}
	err := r.db.QueryRowContext(ctx, query, workOrderID).Scan(
		&metrics.TotalInspections,
		&metrics.AverageQualityScore,
		&metrics.PassedInspections,
		&metrics.FailedInspections,
		&metrics.ReworkRequiredCount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get quality metrics: %w", err)
	}

	metrics.WorkOrderID = workOrderID
	if metrics.TotalInspections > 0 {
		metrics.PassRate = float64(metrics.PassedInspections) / float64(metrics.TotalInspections) * 100
	}

	return metrics, nil
}

// GetInspectionDashboard gets inspection dashboard data
func (r *QualityInspectionRepository) GetInspectionDashboard(ctx context.Context) (*repair.InspectionDashboard, error) {
	query := `
		SELECT 
			COUNT(*) as total_inspections,
			COUNT(CASE WHEN inspection_status = 'scheduled' THEN 1 END) as scheduled_inspections,
			COUNT(CASE WHEN inspection_status = 'in_progress' THEN 1 END) as in_progress_inspections,
			COUNT(CASE WHEN inspection_status = 'completed' THEN 1 END) as completed_inspections,
			COUNT(CASE WHEN pass_fail_status = 'pass' THEN 1 END) as passed_inspections,
			COUNT(CASE WHEN pass_fail_status = 'fail' THEN 1 END) as failed_inspections,
			COUNT(CASE WHEN rework_required = true THEN 1 END) as rework_required_count,
			COALESCE(AVG(quality_score), 0) as overall_quality_score
		FROM quality_inspections 
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'`

	dashboard := &repair.InspectionDashboard{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&dashboard.TotalInspections,
		&dashboard.ScheduledInspections,
		&dashboard.InProgressInspections,
		&dashboard.CompletedInspections,
		&dashboard.PassedInspections,
		&dashboard.FailedInspections,
		&dashboard.ReworkRequiredCount,
		&dashboard.OverallQualityScore,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get inspection dashboard: %w", err)
	}

	if dashboard.TotalInspections > 0 {
		dashboard.PassRate = float64(dashboard.PassedInspections) / float64(dashboard.TotalInspections) * 100
	}

	return dashboard, nil
}

// Stub implementations for list methods
func (r *QualityInspectionRepository) List(ctx context.Context, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *QualityInspectionRepository) GetByWorkOrderID(ctx context.Context, workOrderID int, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *QualityInspectionRepository) GetByInspector(ctx context.Context, inspectorID int, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *QualityInspectionRepository) GetByStatus(ctx context.Context, status string, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *QualityInspectionRepository) GetByType(ctx context.Context, inspectionType string, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *QualityInspectionRepository) GetFailedInspections(ctx context.Context, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *QualityInspectionRepository) GetReworkRequired(ctx context.Context, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *QualityInspectionRepository) Search(ctx context.Context, query string, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}