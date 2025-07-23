package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/repair"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// RepairPartsUsageRepository implements interfaces.RepairPartsUsageRepository
type RepairPartsUsageRepository struct {
	db *sql.DB
}

// NewRepairPartsUsageRepository creates a new repair parts usage repository
func NewRepairPartsUsageRepository(db *sql.DB) interfaces.RepairPartsUsageRepository {
	return &RepairPartsUsageRepository{db: db}
}

// Create creates a new repair parts usage record
func (r *RepairPartsUsageRepository) Create(ctx context.Context, usage *repair.RepairPartsUsage) (*repair.RepairPartsUsage, error) {
	query := `
		INSERT INTO repair_parts_usage (
			work_detail_id, product_id, quantity_used, unit_cost, total_cost,
			usage_type, usage_notes, usage_status, issued_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING parts_usage_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		usage.WorkDetailID,
		usage.ProductID,
		usage.QuantityUsed,
		usage.UnitCost,
		usage.TotalCost,
		usage.UsageType,
		usage.UsageNotes,
		usage.UsageStatus,
		usage.IssuedBy,
	).Scan(&usage.PartsUsageID, &usage.CreatedAt, &usage.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create repair parts usage: %w", err)
	}

	return usage, nil
}

// GetByID retrieves a repair parts usage by ID
func (r *RepairPartsUsageRepository) GetByID(ctx context.Context, id int) (*repair.RepairPartsUsage, error) {
	query := `
		SELECT 
			rpu.parts_usage_id, rpu.work_detail_id, rpu.product_id, rpu.quantity_used,
			rpu.unit_cost, rpu.total_cost, rpu.usage_type, rpu.usage_notes,
			rpu.usage_status, rpu.issued_by, rpu.approved_by, rpu.issued_at,
			rpu.approved_at, rpu.created_at, rpu.updated_at,
			p.product_name, p.product_code,
			rwd.task_description,
			u1.full_name as issued_by_name,
			u2.full_name as approved_by_name
		FROM repair_parts_usage rpu
		LEFT JOIN products p ON rpu.product_id = p.product_id
		LEFT JOIN repair_work_details rwd ON rpu.work_detail_id = rwd.work_detail_id
		LEFT JOIN users u1 ON rpu.issued_by = u1.user_id
		LEFT JOIN users u2 ON rpu.approved_by = u2.user_id
		WHERE rpu.parts_usage_id = $1`

	usage := &repair.RepairPartsUsage{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&usage.PartsUsageID,
		&usage.WorkDetailID,
		&usage.ProductID,
		&usage.QuantityUsed,
		&usage.UnitCost,
		&usage.TotalCost,
		&usage.UsageType,
		&usage.UsageNotes,
		&usage.UsageStatus,
		&usage.IssuedBy,
		&usage.ApprovedBy,
		&usage.IssuedAt,
		&usage.ApprovedAt,
		&usage.CreatedAt,
		&usage.UpdatedAt,
		&usage.ProductName,
		&usage.ProductCode,
		&usage.TaskDescription,
		&usage.IssuedByName,
		&usage.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("repair parts usage not found")
		}
		return nil, fmt.Errorf("failed to get repair parts usage: %w", err)
	}

	return usage, nil
}

// Update updates a repair parts usage record
func (r *RepairPartsUsageRepository) Update(ctx context.Context, id int, usage *repair.RepairPartsUsage) (*repair.RepairPartsUsage, error) {
	query := `
		UPDATE repair_parts_usage SET
			quantity_used = $1, unit_cost = $2, total_cost = $3,
			usage_type = $4, usage_notes = $5, usage_status = $6, updated_at = NOW()
		WHERE parts_usage_id = $7
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		usage.QuantityUsed,
		usage.UnitCost,
		usage.TotalCost,
		usage.UsageType,
		usage.UsageNotes,
		usage.UsageStatus,
		id,
	).Scan(&usage.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update repair parts usage: %w", err)
	}

	usage.PartsUsageID = id
	return usage, nil
}

// Delete deletes a repair parts usage record
func (r *RepairPartsUsageRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM repair_parts_usage WHERE parts_usage_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete repair parts usage: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("repair parts usage not found")
	}

	return nil
}

// ProcessApproval processes parts usage approval
func (r *RepairPartsUsageRepository) ProcessApproval(ctx context.Context, id int, request *repair.PartsUsageApprovalRequest, approvedBy int) error {
	query := `
		UPDATE repair_parts_usage 
		SET usage_status = $1, approved_by = $2, approved_at = NOW(), updated_at = NOW()
		WHERE parts_usage_id = $3`

	_, err := r.db.ExecContext(ctx, query, request.Status, approvedBy, id)
	if err != nil {
		return fmt.Errorf("failed to process approval: %w", err)
	}

	return nil
}

// IssuePartsForRepair issues parts for a repair work detail
func (r *RepairPartsUsageRepository) IssuePartsForRepair(ctx context.Context, workDetailID int, request *repair.PartsUsageIssueRequest, issuedBy int) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create parts usage record
	query := `
		INSERT INTO repair_parts_usage (
			work_detail_id, product_id, quantity_used, unit_cost, total_cost,
			usage_type, usage_notes, usage_status, issued_by, issued_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 'issued', $8, NOW())`

	_, err = tx.ExecContext(ctx, query,
		workDetailID,
		request.ProductID,
		request.QuantityUsed,
		request.UnitCost,
		request.TotalCost,
		request.UsageType,
		request.UsageNotes,
		issuedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create parts usage record: %w", err)
	}

	// Update stock movement (this would integrate with stock system)
	// This is a simplified version - in reality you'd need to properly integrate with the stock management
	stockQuery := `
		INSERT INTO stock_movements (
			product_id, movement_type, quantity, reference_type, reference_id,
			movement_reason, created_by
		) VALUES ($1, 'OUT', $2, 'repair_usage', $3, 'Parts used for repair', $4)`

	_, err = tx.ExecContext(ctx, stockQuery,
		request.ProductID,
		request.QuantityUsed,
		workDetailID,
		issuedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create stock movement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetPartsUsageSummary gets parts usage summary for a work order
func (r *RepairPartsUsageRepository) GetPartsUsageSummary(ctx context.Context, workOrderID int) (*repair.PartsUsageSummary, error) {
	query := `
		SELECT 
			COUNT(DISTINCT rpu.parts_usage_id) as total_parts_used,
			COALESCE(SUM(rpu.total_cost), 0) as total_parts_cost,
			COUNT(DISTINCT rpu.product_id) as unique_parts_count,
			COALESCE(SUM(rpu.quantity_used), 0) as total_quantity_used
		FROM repair_parts_usage rpu
		JOIN repair_work_details rwd ON rpu.work_detail_id = rwd.work_detail_id
		WHERE rwd.work_order_id = $1`

	summary := &repair.PartsUsageSummary{}
	err := r.db.QueryRowContext(ctx, query, workOrderID).Scan(
		&summary.TotalPartsUsed,
		&summary.TotalPartsCost,
		&summary.UniquePartsCount,
		&summary.TotalQuantityUsed,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get parts usage summary: %w", err)
	}

	summary.WorkOrderID = workOrderID

	return summary, nil
}

// GetInventoryImpact gets inventory impact for a work order
func (r *RepairPartsUsageRepository) GetInventoryImpact(ctx context.Context, workOrderID int) ([]*repair.PartsInventoryImpact, error) {
	query := `
		SELECT 
			p.product_id, p.product_name, p.product_code,
			COALESCE(SUM(rpu.quantity_used), 0) as total_used,
			COALESCE(SUM(rpu.total_cost), 0) as total_cost,
			p.current_stock,
			(p.current_stock - COALESCE(SUM(rpu.quantity_used), 0)) as remaining_stock
		FROM repair_parts_usage rpu
		JOIN repair_work_details rwd ON rpu.work_detail_id = rwd.work_detail_id
		JOIN products p ON rpu.product_id = p.product_id
		WHERE rwd.work_order_id = $1
		GROUP BY p.product_id, p.product_name, p.product_code, p.current_stock`

	rows, err := r.db.QueryContext(ctx, query, workOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory impact: %w", err)
	}
	defer rows.Close()

	var impacts []*repair.PartsInventoryImpact
	for rows.Next() {
		impact := &repair.PartsInventoryImpact{}
		err := rows.Scan(
			&impact.ProductID,
			&impact.ProductName,
			&impact.ProductCode,
			&impact.TotalUsed,
			&impact.TotalCost,
			&impact.CurrentStock,
			&impact.RemainingStock,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory impact: %w", err)
		}

		impact.WorkOrderID = workOrderID
		impacts = append(impacts, impact)
	}

	return impacts, nil
}

// Stub implementations for list methods
func (r *RepairPartsUsageRepository) List(ctx context.Context, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairPartsUsageRepository) GetByWorkDetailID(ctx context.Context, workDetailID int, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairPartsUsageRepository) GetByProductID(ctx context.Context, productID int, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairPartsUsageRepository) GetByUsageType(ctx context.Context, usageType string, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairPartsUsageRepository) GetPendingApproval(ctx context.Context, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairPartsUsageRepository) Search(ctx context.Context, query string, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}