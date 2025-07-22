package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// SupplierRepository implements interfaces.SupplierRepository
type SupplierRepository struct {
	db *sql.DB
}

// NewSupplierRepository creates a new supplier repository
func NewSupplierRepository(db *sql.DB) interfaces.SupplierRepository {
	return &SupplierRepository{db: db}
}

// Create creates a new supplier
func (r *SupplierRepository) Create(ctx context.Context, supplier *master.Supplier) (*master.Supplier, error) {
	query := `
		INSERT INTO suppliers (supplier_code, supplier_name, supplier_type, phone, email, address, city, postal_code, tax_number, contact_person, bank_account, payment_terms, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING supplier_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		supplier.SupplierCode,
		supplier.SupplierName,
		supplier.SupplierType,
		supplier.Phone,
		supplier.Email,
		supplier.Address,
		supplier.City,
		supplier.PostalCode,
		supplier.TaxNumber,
		supplier.ContactPerson,
		supplier.BankAccount,
		supplier.PaymentTerms,
		supplier.Notes,
		supplier.CreatedBy,
	).Scan(&supplier.SupplierID, &supplier.CreatedAt, &supplier.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	supplier.IsActive = true
	return supplier, nil
}

// GetByID retrieves a supplier by ID
func (r *SupplierRepository) GetByID(ctx context.Context, id int) (*master.Supplier, error) {
	query := `
		SELECT supplier_id, supplier_code, supplier_name, supplier_type, phone, email, address, city, postal_code, tax_number, contact_person, bank_account, payment_terms, notes, is_active, created_at, updated_at, created_by
		FROM suppliers
		WHERE supplier_id = $1`

	supplier := &master.Supplier{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&supplier.SupplierID,
		&supplier.SupplierCode,
		&supplier.SupplierName,
		&supplier.SupplierType,
		&supplier.Phone,
		&supplier.Email,
		&supplier.Address,
		&supplier.City,
		&supplier.PostalCode,
		&supplier.TaxNumber,
		&supplier.ContactPerson,
		&supplier.BankAccount,
		&supplier.PaymentTerms,
		&supplier.Notes,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
		&supplier.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}

	return supplier, nil
}

// GetByCode retrieves a supplier by code
func (r *SupplierRepository) GetByCode(ctx context.Context, code string) (*master.Supplier, error) {
	query := `
		SELECT supplier_id, supplier_code, supplier_name, supplier_type, phone, email, address, city, postal_code, tax_number, contact_person, bank_account, payment_terms, notes, is_active, created_at, updated_at, created_by
		FROM suppliers
		WHERE supplier_code = $1`

	supplier := &master.Supplier{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&supplier.SupplierID,
		&supplier.SupplierCode,
		&supplier.SupplierName,
		&supplier.SupplierType,
		&supplier.Phone,
		&supplier.Email,
		&supplier.Address,
		&supplier.City,
		&supplier.PostalCode,
		&supplier.TaxNumber,
		&supplier.ContactPerson,
		&supplier.BankAccount,
		&supplier.PaymentTerms,
		&supplier.Notes,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
		&supplier.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}

	return supplier, nil
}

// Update updates a supplier
func (r *SupplierRepository) Update(ctx context.Context, id int, supplier *master.Supplier) (*master.Supplier, error) {
	query := `
		UPDATE suppliers
		SET supplier_name = $1, supplier_type = $2, phone = $3, email = $4, address = $5, city = $6, postal_code = $7, tax_number = $8, contact_person = $9, bank_account = $10, payment_terms = $11, notes = $12, is_active = $13, updated_at = NOW()
		WHERE supplier_id = $14
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		supplier.SupplierName,
		supplier.SupplierType,
		supplier.Phone,
		supplier.Email,
		supplier.Address,
		supplier.City,
		supplier.PostalCode,
		supplier.TaxNumber,
		supplier.ContactPerson,
		supplier.BankAccount,
		supplier.PaymentTerms,
		supplier.Notes,
		supplier.IsActive,
		id,
	).Scan(&supplier.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}

	supplier.SupplierID = id
	return supplier, nil
}

// Delete soft deletes a supplier
func (r *SupplierRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE suppliers SET is_active = FALSE, updated_at = NOW() WHERE supplier_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete supplier: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("supplier not found")
	}

	return nil
}

// List retrieves suppliers with filtering and pagination
func (r *SupplierRepository) List(ctx context.Context, params *master.SupplierFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.SupplierType != nil {
		conditions = append(conditions, fmt.Sprintf("supplier_type = $%d", argIndex))
		args = append(args, *params.SupplierType)
		argIndex++
	}

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.City != "" {
		conditions = append(conditions, fmt.Sprintf("city ILIKE $%d", argIndex))
		args = append(args, "%"+params.City+"%")
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf("(supplier_name ILIKE $%d OR supplier_code ILIKE $%d OR phone ILIKE $%d OR email ILIKE $%d OR contact_person ILIKE $%d)", argIndex, argIndex, argIndex, argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM suppliers %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count suppliers: %w", err)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT supplier_id, supplier_code, supplier_name, supplier_type, phone, email, city, contact_person, is_active, created_at
		FROM suppliers
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}
	defer rows.Close()

	var suppliers []master.SupplierListItem
	for rows.Next() {
		var supplier master.SupplierListItem
		err := rows.Scan(
			&supplier.SupplierID,
			&supplier.SupplierCode,
			&supplier.SupplierName,
			&supplier.SupplierType,
			&supplier.Phone,
			&supplier.Email,
			&supplier.City,
			&supplier.ContactPerson,
			&supplier.IsActive,
			&supplier.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan supplier: %w", err)
		}
		suppliers = append(suppliers, supplier)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate suppliers: %w", err)
	}

	return &common.PaginatedResponse{
		Data:       suppliers,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// GenerateCode generates a new supplier code
func (r *SupplierRepository) GenerateCode(ctx context.Context) (string, error) {
	query := `
		SELECT supplier_code
		FROM suppliers
		WHERE supplier_code LIKE 'SUPP-%'
		ORDER BY supplier_code DESC
		LIMIT 1`

	var lastCode sql.NullString
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get last supplier code: %w", err)
	}

	nextNumber := 1
	if lastCode.Valid {
		// Extract number from code (e.g., "SUPP-001" -> "001" -> 1)
		parts := strings.Split(lastCode.String, "-")
		if len(parts) == 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil {
				nextNumber = num + 1
			}
		}
	}

	return fmt.Sprintf("SUPP-%03d", nextNumber), nil
}

// IsCodeExists checks if a supplier code already exists
func (r *SupplierRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM suppliers WHERE supplier_code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check supplier code existence: %w", err)
	}

	return exists, nil
}

// IsEmailExists checks if an email already exists (excluding a specific ID)
func (r *SupplierRepository) IsEmailExists(ctx context.Context, email string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM suppliers WHERE email = $1 AND supplier_id != $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check supplier email existence: %w", err)
	}

	return exists, nil
}

// IsPhoneExists checks if a phone already exists (excluding a specific ID)
func (r *SupplierRepository) IsPhoneExists(ctx context.Context, phone string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM suppliers WHERE phone = $1 AND supplier_id != $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, phone, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check supplier phone existence: %w", err)
	}

	return exists, nil
}