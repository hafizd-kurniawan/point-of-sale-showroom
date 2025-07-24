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
			usage_type, part_condition, warranty_period_days, installation_notes, 
			issued_by, used_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING usage_id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		usage.WorkDetailID,
		usage.ProductID,
		usage.QuantityUsed,
		usage.UnitCost,
		usage.TotalCost,
		usage.UsageType,
		usage.PartCondition,
		usage.WarrantyPeriodDays,
		usage.InstallationNotes,
		usage.IssuedBy,
		usage.UsedBy,
	).Scan(&usage.UsageID, &usage.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create repair parts usage: %w", err)
	}

	return usage, nil
}

// GetByID retrieves a repair parts usage by ID
func (r *RepairPartsUsageRepository) GetByID(ctx context.Context, id int) (*repair.RepairPartsUsage, error) {
	query := `
		SELECT 
			rpu.usage_id, rpu.work_detail_id, rpu.product_id, rpu.quantity_used,
			rpu.unit_cost, rpu.total_cost, rpu.usage_date, rpu.usage_type,
			rpu.part_condition, rpu.warranty_period_days, rpu.installation_notes,
			rpu.issued_by, rpu.used_by, rpu.approved_by, rpu.approved_at, rpu.created_at,
			p.product_name, p.product_code,
			rwd.task_description,
			rwo.work_order_number,
			u1.full_name as issued_by_name,
			u2.full_name as used_by_name,
			u3.full_name as approved_by_name
		FROM repair_parts_usage rpu
		LEFT JOIN products_spare_parts p ON rpu.product_id = p.product_id
		LEFT JOIN repair_work_details rwd ON rpu.work_detail_id = rwd.work_detail_id
		LEFT JOIN repair_work_orders rwo ON rwd.work_order_id = rwo.work_order_id
		LEFT JOIN users u1 ON rpu.issued_by = u1.user_id
		LEFT JOIN users u2 ON rpu.used_by = u2.user_id
		LEFT JOIN users u3 ON rpu.approved_by = u3.user_id
		WHERE rpu.usage_id = $1`

	usage := &repair.RepairPartsUsage{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&usage.UsageID,
		&usage.WorkDetailID,
		&usage.ProductID,
		&usage.QuantityUsed,
		&usage.UnitCost,
		&usage.TotalCost,
		&usage.UsageDate,
		&usage.UsageType,
		&usage.PartCondition,
		&usage.WarrantyPeriodDays,
		&usage.InstallationNotes,
		&usage.IssuedBy,
		&usage.UsedBy,
		&usage.ApprovedBy,
		&usage.ApprovedAt,
		&usage.CreatedAt,
		&usage.ProductName,
		&usage.ProductCode,
		&usage.TaskDescription,
		&usage.WorkOrderNumber,
		&usage.IssuedByName,
		&usage.UsedByName,
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
			usage_type = $4, part_condition = $5, warranty_period_days = $6,
			installation_notes = $7
		WHERE usage_id = $8`

	_, err := r.db.ExecContext(ctx, query,
		usage.QuantityUsed,
		usage.UnitCost,
		usage.TotalCost,
		usage.UsageType,
		usage.PartCondition,
		usage.WarrantyPeriodDays,
		usage.InstallationNotes,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update repair parts usage: %w", err)
	}

	usage.UsageID = id
	return usage, nil
}

// Delete deletes a repair parts usage record
func (r *RepairPartsUsageRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM repair_parts_usage WHERE usage_id = $1`

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

	// Create parts usage records for each item
	query := `
		INSERT INTO repair_parts_usage (
			work_detail_id, product_id, quantity_used, unit_cost, total_cost,
			usage_type, part_condition, warranty_period_days, installation_notes, 
			issued_by, used_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	// Get unit cost for each product to calculate total cost
	costQuery := `SELECT cost_price FROM products_spare_parts WHERE product_id = $1`

	for _, item := range request.Items {
		var unitCost float64
		err = tx.QueryRowContext(ctx, costQuery, item.ProductID).Scan(&unitCost)
		if err != nil {
			return fmt.Errorf("failed to get product cost for product %d: %w", item.ProductID, err)
		}

		totalCost := unitCost * float64(item.QuantityRequested)

		_, err = tx.ExecContext(ctx, query,
			workDetailID,
			item.ProductID,
			item.QuantityRequested,
			unitCost,
			totalCost,
			item.UsageType,
			item.PartCondition,
			item.WarrantyPeriodDays,
			item.InstallationNotes,
			issuedBy,
			issuedBy, // For now, issued_by and used_by are the same
		)

		if err != nil {
			return fmt.Errorf("failed to create parts usage record for product %d: %w", item.ProductID, err)
		}

		// Update stock movement (this would integrate with stock system)
		stockQuery := `
			INSERT INTO stock_movements (
				product_id, movement_type, reference_type, reference_id,
				quantity_before, quantity_moved, quantity_after, unit_cost, total_value,
				movement_reason, processed_by
			) VALUES ($1, 'out', 'repair', $2, 0, $3, 0, $4, $5, 'Parts used for repair', $6)`

		_, err = tx.ExecContext(ctx, stockQuery,
			item.ProductID,
			workDetailID,
			item.QuantityRequested,
			unitCost,
			totalCost,
			issuedBy,
		)

		if err != nil {
			return fmt.Errorf("failed to create stock movement for product %d: %w", item.ProductID, err)
		}
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
			COUNT(DISTINCT rpu.usage_id) as total_parts_used,
			COALESCE(SUM(rpu.total_cost), 0) as total_parts_cost,
			COUNT(CASE WHEN rpu.usage_type = 'new' THEN 1 END) as new_parts,
			COUNT(CASE WHEN rpu.usage_type = 'replacement' THEN 1 END) as replacement_parts,
			COUNT(CASE WHEN rpu.usage_type = 'warranty' THEN 1 END) as warranty_parts,
			COUNT(CASE WHEN rpu.part_condition = 'oem' THEN 1 END) as oem_parts,
			COUNT(CASE WHEN rpu.part_condition = 'aftermarket' THEN 1 END) as aftermarket_parts,
			COUNT(CASE WHEN rpu.approved_by IS NOT NULL THEN 1 END) as approved_usage,
			COUNT(CASE WHEN rpu.approved_by IS NULL THEN 1 END) as pending_approval
		FROM repair_parts_usage rpu
		JOIN repair_work_details rwd ON rpu.work_detail_id = rwd.work_detail_id
		WHERE rwd.work_order_id = $1`

	summary := &repair.PartsUsageSummary{}
	err := r.db.QueryRowContext(ctx, query, workOrderID).Scan(
		&summary.TotalPartsUsed,
		&summary.TotalPartsCost,
		&summary.NewParts,
		&summary.ReplacementParts,
		&summary.WarrantyParts,
		&summary.OEMParts,
		&summary.AftermarketParts,
		&summary.ApprovedUsage,
		&summary.PendingApproval,
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
			COALESCE(SUM(rpu.quantity_used), 0) as quantity_used,
			COALESCE(SUM(rpu.total_cost), 0) as total_value,
			p.stock_quantity as quantity_before,
			(p.stock_quantity - COALESCE(SUM(rpu.quantity_used), 0)) as quantity_after,
			CASE WHEN (p.stock_quantity - COALESCE(SUM(rpu.quantity_used), 0)) < p.min_stock_level THEN true ELSE false END as low_stock_alert,
			CASE WHEN (p.stock_quantity - COALESCE(SUM(rpu.quantity_used), 0)) <= p.min_stock_level THEN true ELSE false END as reorder_required
		FROM repair_parts_usage rpu
		JOIN repair_work_details rwd ON rpu.work_detail_id = rwd.work_detail_id
		JOIN products_spare_parts p ON rpu.product_id = p.product_id
		WHERE rwd.work_order_id = $1
		GROUP BY p.product_id, p.product_name, p.product_code, p.stock_quantity, p.min_stock_level`

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
			&impact.QuantityUsed,
			&impact.TotalValue,
			&impact.QuantityBefore,
			&impact.QuantityAfter,
			&impact.LowStockAlert,
			&impact.ReorderRequired,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory impact: %w", err)
		}

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