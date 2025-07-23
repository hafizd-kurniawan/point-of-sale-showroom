package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/repair"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// RepairWorkOrderRepository implements interfaces.RepairWorkOrderRepository
type RepairWorkOrderRepository struct {
	db *sql.DB
}

// NewRepairWorkOrderRepository creates a new repair work order repository
func NewRepairWorkOrderRepository(db *sql.DB) interfaces.RepairWorkOrderRepository {
	return &RepairWorkOrderRepository{db: db}
}

// Create creates a new repair work order
func (r *RepairWorkOrderRepository) Create(ctx context.Context, workOrder *repair.RepairWorkOrder) (*repair.RepairWorkOrder, error) {
	query := `
		INSERT INTO repair_work_orders (
			work_order_number, transaction_id, priority_level, estimated_cost,
			estimated_duration_hours, work_description, special_instructions,
			work_order_status, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING work_order_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		workOrder.WorkOrderNumber,
		workOrder.TransactionID,
		workOrder.PriorityLevel,
		workOrder.EstimatedCost,
		workOrder.EstimatedDurationHours,
		workOrder.WorkDescription,
		workOrder.SpecialInstructions,
		workOrder.WorkOrderStatus,
		workOrder.CreatedBy,
	).Scan(&workOrder.WorkOrderID, &workOrder.CreatedAt, &workOrder.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create repair work order: %w", err)
	}

	return workOrder, nil
}

// GetByID retrieves a repair work order by ID
func (r *RepairWorkOrderRepository) GetByID(ctx context.Context, id int) (*repair.RepairWorkOrder, error) {
	query := `
		SELECT 
			rwo.work_order_id, rwo.work_order_number, rwo.transaction_id, rwo.priority_level,
			rwo.estimated_cost, rwo.actual_cost, rwo.estimated_duration_hours, rwo.actual_duration_hours,
			rwo.work_description, rwo.special_instructions, rwo.work_order_status, rwo.assigned_mechanic_id,
			rwo.start_date, rwo.target_completion_date, rwo.actual_completion_date, rwo.created_by,
			rwo.approved_by, rwo.approved_at, rwo.created_at, rwo.updated_at,
			vpt.transaction_number,
			u1.full_name as assigned_mechanic_name,
			u2.full_name as created_by_name,
			u3.full_name as approved_by_name
		FROM repair_work_orders rwo
		LEFT JOIN vehicle_purchase_transactions vpt ON rwo.transaction_id = vpt.transaction_id
		LEFT JOIN users u1 ON rwo.assigned_mechanic_id = u1.user_id
		LEFT JOIN users u2 ON rwo.created_by = u2.user_id
		LEFT JOIN users u3 ON rwo.approved_by = u3.user_id
		WHERE rwo.work_order_id = $1`

	workOrder := &repair.RepairWorkOrder{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&workOrder.WorkOrderID,
		&workOrder.WorkOrderNumber,
		&workOrder.TransactionID,
		&workOrder.PriorityLevel,
		&workOrder.EstimatedCost,
		&workOrder.ActualCost,
		&workOrder.EstimatedDurationHours,
		&workOrder.ActualDurationHours,
		&workOrder.WorkDescription,
		&workOrder.SpecialInstructions,
		&workOrder.WorkOrderStatus,
		&workOrder.AssignedMechanicID,
		&workOrder.StartDate,
		&workOrder.TargetCompletionDate,
		&workOrder.ActualCompletionDate,
		&workOrder.CreatedBy,
		&workOrder.ApprovedBy,
		&workOrder.ApprovedAt,
		&workOrder.CreatedAt,
		&workOrder.UpdatedAt,
		&workOrder.TransactionNumber,
		&workOrder.AssignedMechanicName,
		&workOrder.CreatedByName,
		&workOrder.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("repair work order not found")
		}
		return nil, fmt.Errorf("failed to get repair work order: %w", err)
	}

	return workOrder, nil
}

// GetByNumber retrieves a repair work order by number
func (r *RepairWorkOrderRepository) GetByNumber(ctx context.Context, number string) (*repair.RepairWorkOrder, error) {
	query := `
		SELECT 
			rwo.work_order_id, rwo.work_order_number, rwo.transaction_id, rwo.priority_level,
			rwo.estimated_cost, rwo.actual_cost, rwo.estimated_duration_hours, rwo.actual_duration_hours,
			rwo.work_description, rwo.special_instructions, rwo.work_order_status, rwo.assigned_mechanic_id,
			rwo.start_date, rwo.target_completion_date, rwo.actual_completion_date, rwo.created_by,
			rwo.approved_by, rwo.approved_at, rwo.created_at, rwo.updated_at,
			vpt.transaction_number,
			u1.full_name as assigned_mechanic_name,
			u2.full_name as created_by_name,
			u3.full_name as approved_by_name
		FROM repair_work_orders rwo
		LEFT JOIN vehicle_purchase_transactions vpt ON rwo.transaction_id = vpt.transaction_id
		LEFT JOIN users u1 ON rwo.assigned_mechanic_id = u1.user_id
		LEFT JOIN users u2 ON rwo.created_by = u2.user_id
		LEFT JOIN users u3 ON rwo.approved_by = u3.user_id
		WHERE rwo.work_order_number = $1`

	workOrder := &repair.RepairWorkOrder{}
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&workOrder.WorkOrderID,
		&workOrder.WorkOrderNumber,
		&workOrder.TransactionID,
		&workOrder.PriorityLevel,
		&workOrder.EstimatedCost,
		&workOrder.ActualCost,
		&workOrder.EstimatedDurationHours,
		&workOrder.ActualDurationHours,
		&workOrder.WorkDescription,
		&workOrder.SpecialInstructions,
		&workOrder.WorkOrderStatus,
		&workOrder.AssignedMechanicID,
		&workOrder.StartDate,
		&workOrder.TargetCompletionDate,
		&workOrder.ActualCompletionDate,
		&workOrder.CreatedBy,
		&workOrder.ApprovedBy,
		&workOrder.ApprovedAt,
		&workOrder.CreatedAt,
		&workOrder.UpdatedAt,
		&workOrder.TransactionNumber,
		&workOrder.AssignedMechanicName,
		&workOrder.CreatedByName,
		&workOrder.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("repair work order not found")
		}
		return nil, fmt.Errorf("failed to get repair work order: %w", err)
	}

	return workOrder, nil
}

// Update updates a repair work order
func (r *RepairWorkOrderRepository) Update(ctx context.Context, id int, workOrder *repair.RepairWorkOrder) (*repair.RepairWorkOrder, error) {
	query := `
		UPDATE repair_work_orders SET
			priority_level = $1, estimated_cost = $2, actual_cost = $3,
			estimated_duration_hours = $4, actual_duration_hours = $5, work_description = $6,
			special_instructions = $7, work_order_status = $8, start_date = $9,
			target_completion_date = $10, actual_completion_date = $11, updated_at = NOW()
		WHERE work_order_id = $12
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		workOrder.PriorityLevel,
		workOrder.EstimatedCost,
		workOrder.ActualCost,
		workOrder.EstimatedDurationHours,
		workOrder.ActualDurationHours,
		workOrder.WorkDescription,
		workOrder.SpecialInstructions,
		workOrder.WorkOrderStatus,
		workOrder.StartDate,
		workOrder.TargetCompletionDate,
		workOrder.ActualCompletionDate,
		id,
	).Scan(&workOrder.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update repair work order: %w", err)
	}

	workOrder.WorkOrderID = id
	return workOrder, nil
}

// UpdateStatus updates the status of a repair work order
func (r *RepairWorkOrderRepository) UpdateStatus(ctx context.Context, id int, request *repair.WorkOrderStatusRequest) error {
	query := `
		UPDATE repair_work_orders 
		SET work_order_status = $1, updated_at = NOW()
		WHERE work_order_id = $2`

	_, err := r.db.ExecContext(ctx, query, request.Status, id)
	if err != nil {
		return fmt.Errorf("failed to update work order status: %w", err)
	}

	return nil
}

// Delete deletes a repair work order
func (r *RepairWorkOrderRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM repair_work_orders WHERE work_order_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete repair work order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("repair work order not found")
	}

	return nil
}

// AssignMechanic assigns a mechanic to a work order
func (r *RepairWorkOrderRepository) AssignMechanic(ctx context.Context, id int, request *repair.WorkOrderAssignmentRequest) error {
	query := `
		UPDATE repair_work_orders 
		SET assigned_mechanic_id = $1, start_date = $2, target_completion_date = $3,
			work_order_status = 'assigned', updated_at = NOW()
		WHERE work_order_id = $4`

	_, err := r.db.ExecContext(ctx, query,
		request.MechanicID,
		request.StartDate,
		request.TargetCompletionDate,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to assign mechanic: %w", err)
	}

	return nil
}

// ProcessApproval processes work order approval
func (r *RepairWorkOrderRepository) ProcessApproval(ctx context.Context, id int, request *repair.WorkOrderApprovalRequest, approvedBy int) error {
	query := `
		UPDATE repair_work_orders 
		SET work_order_status = $1, approved_by = $2, approved_at = NOW(), updated_at = NOW()
		WHERE work_order_id = $3`

	_, err := r.db.ExecContext(ctx, query, request.Status, approvedBy, id)
	if err != nil {
		return fmt.Errorf("failed to process approval: %w", err)
	}

	return nil
}

// GenerateNumber generates a unique work order number
func (r *RepairWorkOrderRepository) GenerateNumber(ctx context.Context) (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("WO%d%02d", now.Year(), now.Month())

	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(work_order_number FROM LENGTH($1) + 1) AS INTEGER)), 0) + 1
		FROM repair_work_orders 
		WHERE work_order_number LIKE $1 || '%'`

	var nextNum int
	err := r.db.QueryRowContext(ctx, query, prefix).Scan(&nextNum)
	if err != nil {
		return "", fmt.Errorf("failed to generate work order number: %w", err)
	}

	return fmt.Sprintf("%s%04d", prefix, nextNum), nil
}

// IsNumberExists checks if a work order number already exists
func (r *RepairWorkOrderRepository) IsNumberExists(ctx context.Context, number string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM repair_work_orders WHERE work_order_number = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, number).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if work order number exists: %w", err)
	}

	return exists, nil
}

// GetWorkOrderSummary gets work order summary for a transaction
func (r *RepairWorkOrderRepository) GetWorkOrderSummary(ctx context.Context, transactionID int) (*repair.WorkOrderSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_work_orders,
			COALESCE(SUM(estimated_cost), 0) as total_estimated_cost,
			COALESCE(SUM(actual_cost), 0) as total_actual_cost,
			COUNT(CASE WHEN work_order_status = 'completed' THEN 1 END) as completed_orders,
			COUNT(CASE WHEN work_order_status = 'in_progress' THEN 1 END) as in_progress_orders
		FROM repair_work_orders 
		WHERE transaction_id = $1`

	summary := &repair.WorkOrderSummary{}
	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(
		&summary.TotalWorkOrders,
		&summary.TotalEstimatedCost,
		&summary.TotalActualCost,
		&summary.CompletedOrders,
		&summary.InProgressOrders,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get work order summary: %w", err)
	}

	summary.TransactionID = transactionID

	return summary, nil
}

// Stub implementations for list methods
func (r *RepairWorkOrderRepository) List(ctx context.Context, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkOrderRepository) GetByTransactionID(ctx context.Context, transactionID int, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkOrderRepository) GetByStatus(ctx context.Context, status string, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkOrderRepository) GetByMechanic(ctx context.Context, mechanicID int, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkOrderRepository) GetPendingApproval(ctx context.Context, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *RepairWorkOrderRepository) Search(ctx context.Context, query string, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}