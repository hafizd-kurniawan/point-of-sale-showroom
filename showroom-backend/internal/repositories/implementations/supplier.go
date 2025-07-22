package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type supplierRepository struct {
	db *sql.DB
}

// NewSupplierRepository creates a new supplier repository
func NewSupplierRepository(db *sql.DB) interfaces.SupplierRepository {
	return &supplierRepository{db: db}
}

// Create creates a new supplier
func (r *supplierRepository) Create(ctx context.Context, supplier *master.Supplier) (*master.Supplier, error) {
	query := `
		INSERT INTO suppliers (supplier_code, supplier_name, contact_person, phone, email, address, 
		                       city, province, postal_code, tax_number, supplier_type, credit_limit,
		                       credit_term_days, created_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING supplier_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		supplier.SupplierCode, supplier.SupplierName, supplier.ContactPerson, supplier.Phone,
		supplier.Email, supplier.Address, supplier.City, supplier.Province, supplier.PostalCode,
		supplier.TaxNumber, supplier.SupplierType, supplier.CreditLimit, supplier.CreditTermDays,
		supplier.CreatedBy, supplier.Notes,
	).Scan(&supplier.SupplierID, &supplier.CreatedAt, &supplier.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	return supplier, nil
}

// GetByID retrieves a supplier by ID
func (r *supplierRepository) GetByID(ctx context.Context, id int) (*master.Supplier, error) {
	query := `
		SELECT s.supplier_id, s.supplier_code, s.supplier_name, s.contact_person, s.phone, s.email,
		       s.address, s.city, s.province, s.postal_code, s.tax_number, s.supplier_type,
		       s.credit_limit, s.credit_term_days, s.created_at, s.updated_at, s.created_by,
		       s.is_active, s.notes,
		       u.user_id as creator_user_id, u.username as creator_username, u.full_name as creator_full_name
		FROM suppliers s
		LEFT JOIN users u ON s.created_by = u.user_id
		WHERE s.supplier_id = $1 AND s.is_active = true`

	supplier := &master.Supplier{}
	creator := &user.UserCreatorInfo{}
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&supplier.SupplierID, &supplier.SupplierCode, &supplier.SupplierName, &supplier.ContactPerson,
		&supplier.Phone, &supplier.Email, &supplier.Address, &supplier.City, &supplier.Province,
		&supplier.PostalCode, &supplier.TaxNumber, &supplier.SupplierType, &supplier.CreditLimit,
		&supplier.CreditTermDays, &supplier.CreatedAt, &supplier.UpdatedAt, &supplier.CreatedBy,
		&supplier.IsActive, &supplier.Notes,
		&creator.UserID, &creator.Username, &creator.FullName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}

	supplier.Creator = creator
	return supplier, nil
}

// GetByCode retrieves a supplier by code
func (r *supplierRepository) GetByCode(ctx context.Context, code string) (*master.Supplier, error) {
	query := `
		SELECT s.supplier_id, s.supplier_code, s.supplier_name, s.contact_person, s.phone, s.email,
		       s.address, s.city, s.province, s.postal_code, s.tax_number, s.supplier_type,
		       s.credit_limit, s.credit_term_days, s.created_at, s.updated_at, s.created_by,
		       s.is_active, s.notes
		FROM suppliers s
		WHERE s.supplier_code = $1 AND s.is_active = true`

	supplier := &master.Supplier{}
	
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&supplier.SupplierID, &supplier.SupplierCode, &supplier.SupplierName, &supplier.ContactPerson,
		&supplier.Phone, &supplier.Email, &supplier.Address, &supplier.City, &supplier.Province,
		&supplier.PostalCode, &supplier.TaxNumber, &supplier.SupplierType, &supplier.CreditLimit,
		&supplier.CreditTermDays, &supplier.CreatedAt, &supplier.UpdatedAt, &supplier.CreatedBy,
		&supplier.IsActive, &supplier.Notes,
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
func (r *supplierRepository) Update(ctx context.Context, id int, supplier *master.Supplier) (*master.Supplier, error) {
	query := `
		UPDATE suppliers 
		SET supplier_name = $1, contact_person = $2, phone = $3, email = $4, address = $5, 
		    city = $6, province = $7, postal_code = $8, supplier_type = $9, credit_limit = $10,
		    credit_term_days = $11, is_active = $12, notes = $13, updated_at = NOW()
		WHERE supplier_id = $14 AND is_active = true
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		supplier.SupplierName, supplier.ContactPerson, supplier.Phone, supplier.Email,
		supplier.Address, supplier.City, supplier.Province, supplier.PostalCode,
		supplier.SupplierType, supplier.CreditLimit, supplier.CreditTermDays,
		supplier.IsActive, supplier.Notes, id,
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
func (r *supplierRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE suppliers SET is_active = false, updated_at = NOW() WHERE supplier_id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete supplier: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("supplier not found")
	}

	return nil
}

// List retrieves suppliers with filtering and pagination
func (r *supplierRepository) List(ctx context.Context, params *master.SupplierFilterParams) ([]master.SupplierListItem, int, error) {
	params.Validate()

	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, "s.is_active = true")

	if params.SupplierType != nil {
		conditions = append(conditions, fmt.Sprintf("s.supplier_type = $%d", argIndex))
		args = append(args, *params.SupplierType)
		argIndex++
	}

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("s.is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.City != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(s.city) = LOWER($%d)", argIndex))
		args = append(args, params.City)
		argIndex++
	}

	if params.Province != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(s.province) = LOWER($%d)", argIndex))
		args = append(args, params.Province)
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf(`(
			LOWER(s.supplier_name) LIKE LOWER($%d) OR 
			LOWER(s.contact_person) LIKE LOWER($%d) OR 
			LOWER(s.phone) LIKE LOWER($%d)
		)`, argIndex, argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total records
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM suppliers s WHERE %s`, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count suppliers: %w", err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT s.supplier_id, s.supplier_code, s.supplier_name, s.contact_person, s.phone,
		       s.email, s.supplier_type, s.city, s.is_active, s.created_at
		FROM suppliers s
		WHERE %s
		ORDER BY s.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list suppliers: %w", err)
	}
	defer rows.Close()

	var suppliers []master.SupplierListItem
	for rows.Next() {
		var supplier master.SupplierListItem
		err := rows.Scan(
			&supplier.SupplierID, &supplier.SupplierCode, &supplier.SupplierName,
			&supplier.ContactPerson, &supplier.Phone, &supplier.Email,
			&supplier.SupplierType, &supplier.City, &supplier.IsActive, &supplier.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan supplier: %w", err)
		}
		suppliers = append(suppliers, supplier)
	}

	return suppliers, total, nil
}

// ExistsByCode checks if a supplier with the given code exists
func (r *supplierRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM suppliers WHERE supplier_code = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check supplier code existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByCodeExcludingID checks if a supplier with the given code exists excluding a specific ID
func (r *supplierRepository) ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM suppliers WHERE supplier_code = $1 AND supplier_id != $2)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check supplier code existence: %w", err)
	}
	
	return exists, nil
}

// GetNextSupplierCode generates the next supplier code
func (r *supplierRepository) GetNextSupplierCode(ctx context.Context) (string, error) {
	query := `
		SELECT supplier_code 
		FROM suppliers 
		WHERE supplier_code LIKE 'SUPP-%' 
		ORDER BY supplier_id DESC 
		LIMIT 1`
	
	var lastCode string
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// First supplier
			return "SUPP-001", nil
		}
		return "", fmt.Errorf("failed to get last supplier code: %w", err)
	}
	
	// Extract number from code (e.g., "SUPP-001" -> 1)
	parts := strings.Split(lastCode, "-")
	if len(parts) != 2 {
		return "SUPP-001", nil
	}
	
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "SUPP-001", nil
	}
	
	// Generate next code
	nextNum := num + 1
	return fmt.Sprintf("SUPP-%03d", nextNum), nil
}