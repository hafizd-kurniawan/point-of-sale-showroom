package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/repair"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// VehicleDamageRepository implements interfaces.VehicleDamageRepository
type VehicleDamageRepository struct {
	db *sql.DB
}

// NewVehicleDamageRepository creates a new vehicle damage repository
func NewVehicleDamageRepository(db *sql.DB) interfaces.VehicleDamageRepository {
	return &VehicleDamageRepository{db: db}
}

// Create creates a new vehicle damage record
func (r *VehicleDamageRepository) Create(ctx context.Context, damage *repair.VehicleDamage) (*repair.VehicleDamage, error) {
	query := `
		INSERT INTO vehicle_damages (
			transaction_id, damage_category, damage_severity, damage_description,
			damage_location, damage_cause, estimated_cost, assessment_notes,
			photos_json, damage_status, assessed_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING damage_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		damage.TransactionID,
		damage.DamageCategory,
		damage.DamageSeverity,
		damage.DamageDescription,
		damage.DamageLocation,
		damage.DamageCause,
		damage.EstimatedCost,
		damage.AssessmentNotes,
		damage.PhotosJSON,
		damage.DamageStatus,
		damage.AssessedBy,
	).Scan(&damage.DamageID, &damage.CreatedAt, &damage.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle damage: %w", err)
	}

	return damage, nil
}

// GetByID retrieves a vehicle damage by ID
func (r *VehicleDamageRepository) GetByID(ctx context.Context, id int) (*repair.VehicleDamage, error) {
	query := `
		SELECT 
			vd.damage_id, vd.transaction_id, vd.damage_category, vd.damage_severity,
			vd.damage_description, vd.damage_location, vd.damage_cause, vd.estimated_cost,
			vd.assessment_notes, vd.photos_json, vd.damage_status, vd.assessed_by,
			vd.verified_by, vd.assessed_at, vd.verified_at, vd.created_at, vd.updated_at,
			vpt.transaction_number,
			u1.full_name as assessed_by_name,
			u2.full_name as verified_by_name
		FROM vehicle_damages vd
		LEFT JOIN vehicle_purchase_transactions vpt ON vd.transaction_id = vpt.transaction_id
		LEFT JOIN users u1 ON vd.assessed_by = u1.user_id
		LEFT JOIN users u2 ON vd.verified_by = u2.user_id
		WHERE vd.damage_id = $1`

	damage := &repair.VehicleDamage{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&damage.DamageID,
		&damage.TransactionID,
		&damage.DamageCategory,
		&damage.DamageSeverity,
		&damage.DamageDescription,
		&damage.DamageLocation,
		&damage.DamageCause,
		&damage.EstimatedCost,
		&damage.AssessmentNotes,
		&damage.PhotosJSON,
		&damage.DamageStatus,
		&damage.AssessedBy,
		&damage.VerifiedBy,
		&damage.AssessedAt,
		&damage.VerifiedAt,
		&damage.CreatedAt,
		&damage.UpdatedAt,
		&damage.TransactionNumber,
		&damage.AssessedByName,
		&damage.VerifiedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle damage not found")
		}
		return nil, fmt.Errorf("failed to get vehicle damage: %w", err)
	}

	return damage, nil
}

// Update updates a vehicle damage record
func (r *VehicleDamageRepository) Update(ctx context.Context, id int, damage *repair.VehicleDamage) (*repair.VehicleDamage, error) {
	query := `
		UPDATE vehicle_damages SET
			damage_category = $1, damage_severity = $2, damage_description = $3,
			damage_location = $4, damage_cause = $5, estimated_cost = $6,
			assessment_notes = $7, photos_json = $8, damage_status = $9,
			updated_at = NOW()
		WHERE damage_id = $10
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		damage.DamageCategory,
		damage.DamageSeverity,
		damage.DamageDescription,
		damage.DamageLocation,
		damage.DamageCause,
		damage.EstimatedCost,
		damage.AssessmentNotes,
		damage.PhotosJSON,
		damage.DamageStatus,
		id,
	).Scan(&damage.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update vehicle damage: %w", err)
	}

	damage.DamageID = id
	return damage, nil
}

// UpdateStatus updates the status of a vehicle damage
func (r *VehicleDamageRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `
		UPDATE vehicle_damages 
		SET damage_status = $1, updated_at = NOW()
		WHERE damage_id = $2`

	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update damage status: %w", err)
	}

	return nil
}

// Delete deletes a vehicle damage record
func (r *VehicleDamageRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM vehicle_damages WHERE damage_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle damage: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle damage not found")
	}

	return nil
}

// AssessDamage assesses damage with specific request data
func (r *VehicleDamageRepository) AssessDamage(ctx context.Context, id int, request *repair.DamageAssessmentRequest) error {
	query := `
		UPDATE vehicle_damages 
		SET damage_severity = $1, estimated_cost = $2, assessment_notes = $3,
			damage_status = 'assessed', assessed_at = NOW(), updated_at = NOW()
		WHERE damage_id = $4`

	_, err := r.db.ExecContext(ctx, query,
		request.DamageSeverity,
		request.EstimatedCost,
		request.AssessmentNotes,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to assess damage: %w", err)
	}

	return nil
}

// GetDamageSummary gets damage summary for a transaction
func (r *VehicleDamageRepository) GetDamageSummary(ctx context.Context, transactionID int) (*repair.DamageSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_damages,
			COALESCE(SUM(estimated_cost), 0) as total_estimated_cost,
			COUNT(CASE WHEN damage_severity = 'critical' THEN 1 END) as critical_damages,
			COUNT(CASE WHEN damage_severity = 'high' THEN 1 END) as high_damages,
			COUNT(CASE WHEN damage_severity = 'medium' THEN 1 END) as medium_damages,
			COUNT(CASE WHEN damage_severity = 'low' THEN 1 END) as low_damages
		FROM vehicle_damages 
		WHERE transaction_id = $1`

	summary := &repair.DamageSummary{}
	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(
		&summary.TotalDamages,
		&summary.TotalEstimatedCost,
		&summary.CriticalDamages,
		&summary.HighDamages,
		&summary.MediumDamages,
		&summary.LowDamages,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get damage summary: %w", err)
	}

	summary.TransactionID = transactionID

	return summary, nil
}

// Stub implementations for list methods
func (r *VehicleDamageRepository) List(ctx context.Context, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehicleDamageRepository) GetByTransactionID(ctx context.Context, transactionID int, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehicleDamageRepository) GetByCategory(ctx context.Context, category string, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehicleDamageRepository) GetBySeverity(ctx context.Context, severity string, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehicleDamageRepository) GetByStatus(ctx context.Context, status string, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehicleDamageRepository) GetHighPriority(ctx context.Context, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehicleDamageRepository) Search(ctx context.Context, query string, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}