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

// CustomerRepository implements interfaces.CustomerRepository
type CustomerRepository struct {
	db *sql.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *sql.DB) interfaces.CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create creates a new customer
func (r *CustomerRepository) Create(ctx context.Context, customer *master.Customer) (*master.Customer, error) {
	query := `
		INSERT INTO customers (customer_code, customer_name, customer_type, phone, email, address, city, postal_code, tax_number, contact_person, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING customer_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		customer.CustomerCode,
		customer.CustomerName,
		customer.CustomerType,
		customer.Phone,
		customer.Email,
		customer.Address,
		customer.City,
		customer.PostalCode,
		customer.TaxNumber,
		customer.ContactPerson,
		customer.Notes,
		customer.CreatedBy,
	).Scan(&customer.CustomerID, &customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	customer.IsActive = true
	return customer, nil
}

// GetByID retrieves a customer by ID
func (r *CustomerRepository) GetByID(ctx context.Context, id int) (*master.Customer, error) {
	query := `
		SELECT customer_id, customer_code, customer_name, customer_type, phone, email, address, city, postal_code, tax_number, contact_person, notes, is_active, created_at, updated_at, created_by
		FROM customers
		WHERE customer_id = $1`

	customer := &master.Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.CustomerID,
		&customer.CustomerCode,
		&customer.CustomerName,
		&customer.CustomerType,
		&customer.Phone,
		&customer.Email,
		&customer.Address,
		&customer.City,
		&customer.PostalCode,
		&customer.TaxNumber,
		&customer.ContactPerson,
		&customer.Notes,
		&customer.IsActive,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// GetByCode retrieves a customer by code
func (r *CustomerRepository) GetByCode(ctx context.Context, code string) (*master.Customer, error) {
	query := `
		SELECT customer_id, customer_code, customer_name, customer_type, phone, email, address, city, postal_code, tax_number, contact_person, notes, is_active, created_at, updated_at, created_by
		FROM customers
		WHERE customer_code = $1`

	customer := &master.Customer{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&customer.CustomerID,
		&customer.CustomerCode,
		&customer.CustomerName,
		&customer.CustomerType,
		&customer.Phone,
		&customer.Email,
		&customer.Address,
		&customer.City,
		&customer.PostalCode,
		&customer.TaxNumber,
		&customer.ContactPerson,
		&customer.Notes,
		&customer.IsActive,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// Update updates a customer
func (r *CustomerRepository) Update(ctx context.Context, id int, customer *master.Customer) (*master.Customer, error) {
	query := `
		UPDATE customers
		SET customer_name = $1, customer_type = $2, phone = $3, email = $4, address = $5, city = $6, postal_code = $7, tax_number = $8, contact_person = $9, notes = $10, is_active = $11, updated_at = NOW()
		WHERE customer_id = $12
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		customer.CustomerName,
		customer.CustomerType,
		customer.Phone,
		customer.Email,
		customer.Address,
		customer.City,
		customer.PostalCode,
		customer.TaxNumber,
		customer.ContactPerson,
		customer.Notes,
		customer.IsActive,
		id,
	).Scan(&customer.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	customer.CustomerID = id
	return customer, nil
}

// Delete soft deletes a customer
func (r *CustomerRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE customers SET is_active = FALSE, updated_at = NOW() WHERE customer_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// List retrieves customers with filtering and pagination
func (r *CustomerRepository) List(ctx context.Context, params *master.CustomerFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.CustomerType != nil {
		conditions = append(conditions, fmt.Sprintf("customer_type = $%d", argIndex))
		args = append(args, *params.CustomerType)
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
		searchCondition := fmt.Sprintf("(customer_name ILIKE $%d OR customer_code ILIKE $%d OR phone ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex, argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM customers %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count customers: %w", err)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT customer_id, customer_code, customer_name, customer_type, phone, email, city, is_active, created_at
		FROM customers
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var customers []master.CustomerListItem
	for rows.Next() {
		var customer master.CustomerListItem
		err := rows.Scan(
			&customer.CustomerID,
			&customer.CustomerCode,
			&customer.CustomerName,
			&customer.CustomerType,
			&customer.Phone,
			&customer.Email,
			&customer.City,
			&customer.IsActive,
			&customer.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate customers: %w", err)
	}

	return &common.PaginatedResponse{
		Data:       customers,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// GenerateCode generates a new customer code
func (r *CustomerRepository) GenerateCode(ctx context.Context) (string, error) {
	query := `
		SELECT customer_code
		FROM customers
		WHERE customer_code LIKE 'CUST-%'
		ORDER BY customer_code DESC
		LIMIT 1`

	var lastCode sql.NullString
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get last customer code: %w", err)
	}

	nextNumber := 1
	if lastCode.Valid {
		// Extract number from code (e.g., "CUST-001" -> "001" -> 1)
		parts := strings.Split(lastCode.String, "-")
		if len(parts) == 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil {
				nextNumber = num + 1
			}
		}
	}

	return fmt.Sprintf("CUST-%03d", nextNumber), nil
}

// IsCodeExists checks if a customer code already exists
func (r *CustomerRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE customer_code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check customer code existence: %w", err)
	}

	return exists, nil
}

// IsEmailExists checks if an email already exists (excluding a specific ID)
func (r *CustomerRepository) IsEmailExists(ctx context.Context, email string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE email = $1 AND customer_id != $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check customer email existence: %w", err)
	}

	return exists, nil
}

// IsPhoneExists checks if a phone already exists (excluding a specific ID)
func (r *CustomerRepository) IsPhoneExists(ctx context.Context, phone string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE phone = $1 AND customer_id != $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, phone, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check customer phone existence: %w", err)
	}

	return exists, nil
}