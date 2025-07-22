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

type customerRepository struct {
	db *sql.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *sql.DB) interfaces.CustomerRepository {
	return &customerRepository{db: db}
}

// Create creates a new customer
func (r *customerRepository) Create(ctx context.Context, customer *master.Customer) (*master.Customer, error) {
	query := `
		INSERT INTO customers (customer_code, customer_name, phone, email, address, city, province, 
		                       postal_code, id_card_number, tax_number, customer_type, birth_date,
		                       occupation, income_range, created_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING customer_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		customer.CustomerCode, customer.CustomerName, customer.Phone, customer.Email,
		customer.Address, customer.City, customer.Province, customer.PostalCode,
		customer.IDCardNumber, customer.TaxNumber, customer.CustomerType, customer.BirthDate,
		customer.Occupation, customer.IncomeRange, customer.CreatedBy, customer.Notes,
	).Scan(&customer.CustomerID, &customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// GetByID retrieves a customer by ID
func (r *customerRepository) GetByID(ctx context.Context, id int) (*master.Customer, error) {
	query := `
		SELECT c.customer_id, c.customer_code, c.customer_name, c.phone, c.email, c.address,
		       c.city, c.province, c.postal_code, c.id_card_number, c.tax_number, c.customer_type,
		       c.birth_date, c.occupation, c.income_range, c.created_at, c.updated_at, c.created_by,
		       c.is_active, c.notes,
		       u.user_id as creator_user_id, u.username as creator_username, u.full_name as creator_full_name
		FROM customers c
		LEFT JOIN users u ON c.created_by = u.user_id
		WHERE c.customer_id = $1 AND c.is_active = true`

	customer := &master.Customer{}
	creator := &user.UserCreatorInfo{}
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.CustomerID, &customer.CustomerCode, &customer.CustomerName, &customer.Phone,
		&customer.Email, &customer.Address, &customer.City, &customer.Province,
		&customer.PostalCode, &customer.IDCardNumber, &customer.TaxNumber, &customer.CustomerType,
		&customer.BirthDate, &customer.Occupation, &customer.IncomeRange, &customer.CreatedAt,
		&customer.UpdatedAt, &customer.CreatedBy, &customer.IsActive, &customer.Notes,
		&creator.UserID, &creator.Username, &creator.FullName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	customer.Creator = creator
	return customer, nil
}

// GetByCode retrieves a customer by code
func (r *customerRepository) GetByCode(ctx context.Context, code string) (*master.Customer, error) {
	query := `
		SELECT c.customer_id, c.customer_code, c.customer_name, c.phone, c.email, c.address,
		       c.city, c.province, c.postal_code, c.id_card_number, c.tax_number, c.customer_type,
		       c.birth_date, c.occupation, c.income_range, c.created_at, c.updated_at, c.created_by,
		       c.is_active, c.notes
		FROM customers c
		WHERE c.customer_code = $1 AND c.is_active = true`

	customer := &master.Customer{}
	
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&customer.CustomerID, &customer.CustomerCode, &customer.CustomerName, &customer.Phone,
		&customer.Email, &customer.Address, &customer.City, &customer.Province,
		&customer.PostalCode, &customer.IDCardNumber, &customer.TaxNumber, &customer.CustomerType,
		&customer.BirthDate, &customer.Occupation, &customer.IncomeRange, &customer.CreatedAt,
		&customer.UpdatedAt, &customer.CreatedBy, &customer.IsActive, &customer.Notes,
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
func (r *customerRepository) Update(ctx context.Context, id int, customer *master.Customer) (*master.Customer, error) {
	query := `
		UPDATE customers 
		SET customer_name = $1, phone = $2, email = $3, address = $4, city = $5, province = $6,
		    postal_code = $7, customer_type = $8, occupation = $9, income_range = $10,
		    is_active = $11, notes = $12, updated_at = NOW()
		WHERE customer_id = $13 AND is_active = true
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		customer.CustomerName, customer.Phone, customer.Email, customer.Address,
		customer.City, customer.Province, customer.PostalCode, customer.CustomerType,
		customer.Occupation, customer.IncomeRange, customer.IsActive, customer.Notes, id,
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
func (r *customerRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE customers SET is_active = false, updated_at = NOW() WHERE customer_id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// List retrieves customers with filtering and pagination
func (r *customerRepository) List(ctx context.Context, params *master.CustomerFilterParams) ([]master.CustomerListItem, int, error) {
	params.Validate()

	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, "c.is_active = true")

	if params.CustomerType != nil {
		conditions = append(conditions, fmt.Sprintf("c.customer_type = $%d", argIndex))
		args = append(args, *params.CustomerType)
		argIndex++
	}

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("c.is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.City != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(c.city) = LOWER($%d)", argIndex))
		args = append(args, params.City)
		argIndex++
	}

	if params.Province != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(c.province) = LOWER($%d)", argIndex))
		args = append(args, params.Province)
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf(`(
			LOWER(c.customer_name) LIKE LOWER($%d) OR 
			LOWER(c.phone) LIKE LOWER($%d) OR 
			LOWER(c.email) LIKE LOWER($%d) OR 
			LOWER(c.id_card_number) LIKE LOWER($%d)
		)`, argIndex, argIndex, argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total records
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM customers c WHERE %s`, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT c.customer_id, c.customer_code, c.customer_name, c.phone, c.email,
		       c.customer_type, c.city, c.is_active, c.created_at
		FROM customers c
		WHERE %s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var customers []master.CustomerListItem
	for rows.Next() {
		var customer master.CustomerListItem
		err := rows.Scan(
			&customer.CustomerID, &customer.CustomerCode, &customer.CustomerName,
			&customer.Phone, &customer.Email, &customer.CustomerType,
			&customer.City, &customer.IsActive, &customer.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, customer)
	}

	return customers, total, nil
}

// ExistsByCode checks if a customer with the given code exists
func (r *customerRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE customer_code = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check customer code existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByIDCardNumber checks if a customer with the given ID card number exists
func (r *customerRepository) ExistsByIDCardNumber(ctx context.Context, idCard string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE id_card_number = $1 AND is_active = true)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, idCard).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check customer ID card existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByCodeExcludingID checks if a customer with the given code exists excluding a specific ID
func (r *customerRepository) ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE customer_code = $1 AND customer_id != $2)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check customer code existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByIDCardNumberExcludingID checks if a customer with the given ID card exists excluding a specific ID
func (r *customerRepository) ExistsByIDCardNumberExcludingID(ctx context.Context, idCard string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE id_card_number = $1 AND customer_id != $2 AND is_active = true)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, idCard, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check customer ID card existence: %w", err)
	}
	
	return exists, nil
}

// GetNextCustomerCode generates the next customer code
func (r *customerRepository) GetNextCustomerCode(ctx context.Context) (string, error) {
	query := `
		SELECT customer_code 
		FROM customers 
		WHERE customer_code LIKE 'CUST-%' 
		ORDER BY customer_id DESC 
		LIMIT 1`
	
	var lastCode string
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// First customer
			return "CUST-001", nil
		}
		return "", fmt.Errorf("failed to get last customer code: %w", err)
	}
	
	// Extract number from code (e.g., "CUST-001" -> 1)
	parts := strings.Split(lastCode, "-")
	if len(parts) != 2 {
		return "CUST-001", nil
	}
	
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "CUST-001", nil
	}
	
	// Generate next code
	nextNum := num + 1
	return fmt.Sprintf("CUST-%03d", nextNum), nil
}