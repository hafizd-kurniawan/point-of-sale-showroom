package database

import (
	"database/sql"
	"log"
)

// RunMigrations executes database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsersTable,
		createUserSessionsTable,
		createCustomersTable,
		createSuppliersTable,
		createVehicleBrandsTable,
		createVehicleCategoriesTable,
		createVehicleModelsTable,
		createProductCategoriesTable,
		createIndexes,
		createMasterDataIndexes,
		// Phase 3: Products, Inventory & Purchasing System
		createProductsSparePartsTable,
		createPurchaseOrdersPartsTable,
		createPurchaseOrderDetailsTable,
		createGoodsReceiptsTable,
		createGoodsReceiptDetailsTable,
		createStockMovementsTable,
		createStockAdjustmentsTable,
		createSupplierPaymentsTable,
		createPhase3Indexes,
	}

	for i, migration := range migrations {
		log.Printf("Running migration %d...", i+1)
		if _, err := db.Exec(migration); err != nil {
			return err
		}
		log.Printf("Migration %d completed successfully", i+1)
	}

	log.Println("All migrations completed successfully")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address VARCHAR(500),
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin','sales','cashier','mechanic','manager')),
    salary DECIMAL(15,2),
    hire_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    profile_image VARCHAR(500),
    notes TEXT
);`

const createUserSessionsTable = `
CREATE TABLE IF NOT EXISTS user_sessions (
    session_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    session_token VARCHAR(500) UNIQUE NOT NULL,
    login_at TIMESTAMP DEFAULT NOW(),
    logout_at TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE
);`

const createIndexes = `
-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_created_by ON users(created_by);

-- User sessions table indexes
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_is_active ON user_sessions(is_active);
CREATE INDEX IF NOT EXISTS idx_user_sessions_login_at ON user_sessions(login_at);

-- Foreign key constraint for users.created_by (self-referencing)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_users_created_by'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT fk_users_created_by 
        FOREIGN KEY (created_by) REFERENCES users(user_id);
    END IF;
END $$;`

const createCustomersTable = `
CREATE TABLE IF NOT EXISTS customers (
    customer_id SERIAL PRIMARY KEY,
    customer_code VARCHAR(20) UNIQUE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_type VARCHAR(20) NOT NULL CHECK (customer_type IN ('individual','corporate')),
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    address VARCHAR(500) NOT NULL,
    city VARCHAR(100) NOT NULL,
    postal_code VARCHAR(10),
    tax_number VARCHAR(30),
    contact_person VARCHAR(255),
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createSuppliersTable = `
CREATE TABLE IF NOT EXISTS suppliers (
    supplier_id SERIAL PRIMARY KEY,
    supplier_code VARCHAR(20) UNIQUE NOT NULL,
    supplier_name VARCHAR(255) NOT NULL,
    supplier_type VARCHAR(20) NOT NULL CHECK (supplier_type IN ('parts','vehicle','both')),
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    address VARCHAR(500) NOT NULL,
    city VARCHAR(100) NOT NULL,
    postal_code VARCHAR(10),
    tax_number VARCHAR(30),
    contact_person VARCHAR(255) NOT NULL,
    bank_account VARCHAR(100),
    payment_terms VARCHAR(255),
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createVehicleBrandsTable = `
CREATE TABLE IF NOT EXISTS vehicle_brands (
    brand_id SERIAL PRIMARY KEY,
    brand_code VARCHAR(20) UNIQUE NOT NULL,
    brand_name VARCHAR(100) UNIQUE NOT NULL,
    country_origin VARCHAR(100) NOT NULL,
    description TEXT,
    logo_url VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createVehicleCategoriesTable = `
CREATE TABLE IF NOT EXISTS vehicle_categories (
    category_id SERIAL PRIMARY KEY,
    category_code VARCHAR(20) UNIQUE NOT NULL,
    category_name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createVehicleModelsTable = `
CREATE TABLE IF NOT EXISTS vehicle_models (
    model_id SERIAL PRIMARY KEY,
    model_code VARCHAR(20) UNIQUE NOT NULL,
    model_name VARCHAR(100) NOT NULL,
    brand_id INTEGER NOT NULL REFERENCES vehicle_brands(brand_id),
    category_id INTEGER NOT NULL REFERENCES vehicle_categories(category_id),
    model_year INTEGER NOT NULL CHECK (model_year >= 1900 AND model_year <= 2100),
    engine_capacity DECIMAL(5,2),
    fuel_type VARCHAR(50) NOT NULL,
    transmission VARCHAR(50) NOT NULL,
    seat_capacity INTEGER NOT NULL CHECK (seat_capacity >= 1 AND seat_capacity <= 50),
    color VARCHAR(50) NOT NULL,
    price DECIMAL(15,2) NOT NULL CHECK (price >= 0),
    description TEXT,
    image_url VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createProductCategoriesTable = `
CREATE TABLE IF NOT EXISTS product_categories (
    category_id SERIAL PRIMARY KEY,
    category_code VARCHAR(20) UNIQUE NOT NULL,
    category_name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id INTEGER REFERENCES product_categories(category_id),
    level INTEGER NOT NULL DEFAULT 1,
    path VARCHAR(500) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createMasterDataIndexes = `
-- Customers table indexes
CREATE INDEX IF NOT EXISTS idx_customers_code ON customers(customer_code);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(customer_name);
CREATE INDEX IF NOT EXISTS idx_customers_type ON customers(customer_type);
CREATE INDEX IF NOT EXISTS idx_customers_city ON customers(city);
CREATE INDEX IF NOT EXISTS idx_customers_is_active ON customers(is_active);
CREATE INDEX IF NOT EXISTS idx_customers_created_by ON customers(created_by);

-- Suppliers table indexes
CREATE INDEX IF NOT EXISTS idx_suppliers_code ON suppliers(supplier_code);
CREATE INDEX IF NOT EXISTS idx_suppliers_name ON suppliers(supplier_name);
CREATE INDEX IF NOT EXISTS idx_suppliers_type ON suppliers(supplier_type);
CREATE INDEX IF NOT EXISTS idx_suppliers_city ON suppliers(city);
CREATE INDEX IF NOT EXISTS idx_suppliers_is_active ON suppliers(is_active);
CREATE INDEX IF NOT EXISTS idx_suppliers_created_by ON suppliers(created_by);

-- Vehicle brands table indexes
CREATE INDEX IF NOT EXISTS idx_vehicle_brands_code ON vehicle_brands(brand_code);
CREATE INDEX IF NOT EXISTS idx_vehicle_brands_name ON vehicle_brands(brand_name);
CREATE INDEX IF NOT EXISTS idx_vehicle_brands_country ON vehicle_brands(country_origin);
CREATE INDEX IF NOT EXISTS idx_vehicle_brands_is_active ON vehicle_brands(is_active);
CREATE INDEX IF NOT EXISTS idx_vehicle_brands_created_by ON vehicle_brands(created_by);

-- Vehicle categories table indexes
CREATE INDEX IF NOT EXISTS idx_vehicle_categories_code ON vehicle_categories(category_code);
CREATE INDEX IF NOT EXISTS idx_vehicle_categories_name ON vehicle_categories(category_name);
CREATE INDEX IF NOT EXISTS idx_vehicle_categories_is_active ON vehicle_categories(is_active);
CREATE INDEX IF NOT EXISTS idx_vehicle_categories_created_by ON vehicle_categories(created_by);

-- Vehicle models table indexes
CREATE INDEX IF NOT EXISTS idx_vehicle_models_code ON vehicle_models(model_code);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_name ON vehicle_models(model_name);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_brand_id ON vehicle_models(brand_id);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_category_id ON vehicle_models(category_id);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_year ON vehicle_models(model_year);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_fuel_type ON vehicle_models(fuel_type);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_transmission ON vehicle_models(transmission);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_price ON vehicle_models(price);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_is_active ON vehicle_models(is_active);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_created_by ON vehicle_models(created_by);

-- Product categories table indexes
CREATE INDEX IF NOT EXISTS idx_product_categories_code ON product_categories(category_code);
CREATE INDEX IF NOT EXISTS idx_product_categories_name ON product_categories(category_name);
CREATE INDEX IF NOT EXISTS idx_product_categories_parent_id ON product_categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_product_categories_level ON product_categories(level);
CREATE INDEX IF NOT EXISTS idx_product_categories_path ON product_categories(path);
CREATE INDEX IF NOT EXISTS idx_product_categories_is_active ON product_categories(is_active);
CREATE INDEX IF NOT EXISTS idx_product_categories_created_by ON product_categories(created_by);`

// Phase 3: Products, Inventory & Purchasing System Tables

const createProductsSparePartsTable = `
CREATE TABLE IF NOT EXISTS products_spare_parts (
    product_id SERIAL PRIMARY KEY,
    product_code VARCHAR(20) UNIQUE NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    description TEXT,
    brand_id INTEGER NOT NULL REFERENCES vehicle_brands(brand_id),
    category_id INTEGER NOT NULL REFERENCES product_categories(category_id),
    unit_measure VARCHAR(50) NOT NULL,
    cost_price DECIMAL(15,2) NOT NULL CHECK (cost_price >= 0),
    selling_price DECIMAL(15,2) NOT NULL CHECK (selling_price >= 0),
    markup_percentage DECIMAL(5,2) NOT NULL CHECK (markup_percentage >= 0),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    min_stock_level INTEGER NOT NULL DEFAULT 0 CHECK (min_stock_level >= 0),
    max_stock_level INTEGER NOT NULL DEFAULT 0 CHECK (max_stock_level >= 0),
    location_rack VARCHAR(100),
    barcode VARCHAR(100),
    weight DECIMAL(10,3),
    dimensions VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id),
    is_active BOOLEAN DEFAULT TRUE,
    product_image VARCHAR(500),
    notes TEXT
);`

const createPurchaseOrdersPartsTable = `
CREATE TABLE IF NOT EXISTS purchase_orders_parts (
    po_id SERIAL PRIMARY KEY,
    po_number VARCHAR(20) UNIQUE NOT NULL,
    supplier_id INTEGER NOT NULL REFERENCES suppliers(supplier_id),
    po_date TIMESTAMP NOT NULL DEFAULT NOW(),
    required_date TIMESTAMP,
    expected_delivery_date TIMESTAMP,
    po_type VARCHAR(20) NOT NULL CHECK (po_type IN ('regular','urgent','blanket','contract')) DEFAULT 'regular',
    subtotal DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
    discount_amount DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
    shipping_cost DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (shipping_cost >= 0),
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft','sent','acknowledged','partial_received','received','completed','cancelled')) DEFAULT 'draft',
    payment_terms VARCHAR(20) NOT NULL CHECK (payment_terms IN ('cod','net_30','net_60','advance')) DEFAULT 'net_30',
    payment_due_date TIMESTAMP,
    created_by INTEGER NOT NULL REFERENCES users(user_id),
    approved_by INTEGER REFERENCES users(user_id),
    approved_at TIMESTAMP,
    delivery_address VARCHAR(500),
    po_notes TEXT,
    terms_and_conditions TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);`

const createPurchaseOrderDetailsTable = `
CREATE TABLE IF NOT EXISTS purchase_order_details (
    po_detail_id SERIAL PRIMARY KEY,
    po_id INTEGER NOT NULL REFERENCES purchase_orders_parts(po_id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products_spare_parts(product_id),
    item_description VARCHAR(500),
    quantity_ordered INTEGER NOT NULL CHECK (quantity_ordered > 0),
    quantity_received INTEGER NOT NULL DEFAULT 0 CHECK (quantity_received >= 0),
    quantity_pending INTEGER NOT NULL DEFAULT 0 CHECK (quantity_pending >= 0),
    unit_cost DECIMAL(15,2) NOT NULL CHECK (unit_cost >= 0),
    total_cost DECIMAL(15,2) NOT NULL CHECK (total_cost >= 0),
    expected_date TIMESTAMP,
    received_date TIMESTAMP,
    line_status VARCHAR(20) NOT NULL CHECK (line_status IN ('pending','partial','received','cancelled')) DEFAULT 'pending',
    item_notes TEXT
);`

const createGoodsReceiptsTable = `
CREATE TABLE IF NOT EXISTS goods_receipts (
    receipt_id SERIAL PRIMARY KEY,
    po_id INTEGER NOT NULL REFERENCES purchase_orders_parts(po_id),
    receipt_number VARCHAR(20) UNIQUE NOT NULL,
    receipt_date TIMESTAMP NOT NULL DEFAULT NOW(),
    received_by INTEGER NOT NULL REFERENCES users(user_id),
    supplier_delivery_note VARCHAR(100),
    supplier_invoice_number VARCHAR(100),
    total_received_value DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (total_received_value >= 0),
    receipt_status VARCHAR(20) NOT NULL CHECK (receipt_status IN ('partial','complete','with_discrepancy')) DEFAULT 'partial',
    receipt_notes TEXT,
    discrepancy_notes TEXT,
    receipt_documents_json TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);`

const createGoodsReceiptDetailsTable = `
CREATE TABLE IF NOT EXISTS goods_receipt_details (
    receipt_detail_id SERIAL PRIMARY KEY,
    receipt_id INTEGER NOT NULL REFERENCES goods_receipts(receipt_id) ON DELETE CASCADE,
    po_detail_id INTEGER NOT NULL REFERENCES purchase_order_details(po_detail_id),
    product_id INTEGER NOT NULL REFERENCES products_spare_parts(product_id),
    quantity_received INTEGER NOT NULL CHECK (quantity_received >= 0),
    quantity_accepted INTEGER NOT NULL CHECK (quantity_accepted >= 0),
    quantity_rejected INTEGER NOT NULL CHECK (quantity_rejected >= 0),
    unit_cost DECIMAL(15,2) NOT NULL CHECK (unit_cost >= 0),
    total_cost DECIMAL(15,2) NOT NULL CHECK (total_cost >= 0),
    condition_received VARCHAR(20) NOT NULL CHECK (condition_received IN ('good','damaged','expired','wrong_item')) DEFAULT 'good',
    inspection_notes TEXT,
    rejection_reason TEXT,
    expiry_date TIMESTAMP,
    batch_number VARCHAR(100),
    serial_numbers_json TEXT
);`

const createStockMovementsTable = `
CREATE TABLE IF NOT EXISTS stock_movements (
    movement_id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products_spare_parts(product_id),
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('in','out','transfer','adjustment','damage','expired','return')),
    reference_type VARCHAR(20) NOT NULL CHECK (reference_type IN ('purchase','sales','repair','adjustment','transfer','return')),
    reference_id INTEGER NOT NULL,
    quantity_before INTEGER NOT NULL CHECK (quantity_before >= 0),
    quantity_moved INTEGER NOT NULL,
    quantity_after INTEGER NOT NULL CHECK (quantity_after >= 0),
    unit_cost DECIMAL(15,2) NOT NULL CHECK (unit_cost >= 0),
    total_value DECIMAL(15,2) NOT NULL,
    location_from VARCHAR(100),
    location_to VARCHAR(100),
    movement_date TIMESTAMP NOT NULL DEFAULT NOW(),
    processed_by INTEGER NOT NULL REFERENCES users(user_id),
    movement_reason VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);`

const createStockAdjustmentsTable = `
CREATE TABLE IF NOT EXISTS stock_adjustments (
    adjustment_id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products_spare_parts(product_id),
    adjustment_type VARCHAR(20) NOT NULL CHECK (adjustment_type IN ('physical_count','damage','expired','theft','correction','write_off')),
    quantity_system INTEGER NOT NULL CHECK (quantity_system >= 0),
    quantity_physical INTEGER NOT NULL CHECK (quantity_physical >= 0),
    quantity_variance INTEGER NOT NULL,
    cost_impact DECIMAL(15,2) NOT NULL,
    adjustment_reason VARCHAR(255) NOT NULL,
    notes TEXT,
    approved_by INTEGER REFERENCES users(user_id),
    adjustment_date TIMESTAMP NOT NULL DEFAULT NOW(),
    approved_at TIMESTAMP,
    supporting_documents_json TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createSupplierPaymentsTable = `
CREATE TABLE IF NOT EXISTS supplier_payments (
    payment_id SERIAL PRIMARY KEY,
    supplier_id INTEGER NOT NULL REFERENCES suppliers(supplier_id),
    po_id INTEGER REFERENCES purchase_orders_parts(po_id),
    payment_number VARCHAR(20) UNIQUE NOT NULL,
    invoice_amount DECIMAL(15,2) NOT NULL CHECK (invoice_amount >= 0),
    payment_amount DECIMAL(15,2) NOT NULL CHECK (payment_amount >= 0),
    discount_taken DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (discount_taken >= 0),
    outstanding_amount DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (outstanding_amount >= 0),
    invoice_date TIMESTAMP NOT NULL,
    payment_date TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP NOT NULL,
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('cash','transfer','check','credit')) DEFAULT 'transfer',
    payment_reference VARCHAR(100),
    invoice_number VARCHAR(100) NOT NULL,
    payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN ('pending','partial','paid','overdue','disputed')) DEFAULT 'pending',
    days_overdue INTEGER NOT NULL DEFAULT 0 CHECK (days_overdue >= 0),
    penalty_amount DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (penalty_amount >= 0),
    processed_by INTEGER NOT NULL REFERENCES users(user_id),
    payment_notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);`

const createPhase3Indexes = `
-- Products spare parts table indexes
CREATE INDEX IF NOT EXISTS idx_products_spare_parts_code ON products_spare_parts(product_code);
CREATE INDEX IF NOT EXISTS idx_products_spare_parts_name ON products_spare_parts(product_name);
CREATE INDEX IF NOT EXISTS idx_products_spare_parts_brand_id ON products_spare_parts(brand_id);
CREATE INDEX IF NOT EXISTS idx_products_spare_parts_category_id ON products_spare_parts(category_id);
CREATE INDEX IF NOT EXISTS idx_products_spare_parts_is_active ON products_spare_parts(is_active);
CREATE INDEX IF NOT EXISTS idx_products_spare_parts_stock_quantity ON products_spare_parts(stock_quantity);
CREATE INDEX IF NOT EXISTS idx_products_spare_parts_min_stock ON products_spare_parts(min_stock_level);
CREATE INDEX IF NOT EXISTS idx_products_spare_parts_barcode ON products_spare_parts(barcode);
CREATE INDEX IF NOT EXISTS idx_products_spare_parts_created_by ON products_spare_parts(created_by);

-- Purchase orders parts table indexes
CREATE INDEX IF NOT EXISTS idx_purchase_orders_parts_number ON purchase_orders_parts(po_number);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_parts_supplier_id ON purchase_orders_parts(supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_parts_status ON purchase_orders_parts(status);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_parts_po_date ON purchase_orders_parts(po_date);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_parts_required_date ON purchase_orders_parts(required_date);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_parts_created_by ON purchase_orders_parts(created_by);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_parts_approved_by ON purchase_orders_parts(approved_by);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_parts_po_type ON purchase_orders_parts(po_type);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_parts_payment_terms ON purchase_orders_parts(payment_terms);

-- Purchase order details table indexes
CREATE INDEX IF NOT EXISTS idx_purchase_order_details_po_id ON purchase_order_details(po_id);
CREATE INDEX IF NOT EXISTS idx_purchase_order_details_product_id ON purchase_order_details(product_id);
CREATE INDEX IF NOT EXISTS idx_purchase_order_details_line_status ON purchase_order_details(line_status);

-- Goods receipts table indexes
CREATE INDEX IF NOT EXISTS idx_goods_receipts_number ON goods_receipts(receipt_number);
CREATE INDEX IF NOT EXISTS idx_goods_receipts_po_id ON goods_receipts(po_id);
CREATE INDEX IF NOT EXISTS idx_goods_receipts_receipt_date ON goods_receipts(receipt_date);
CREATE INDEX IF NOT EXISTS idx_goods_receipts_received_by ON goods_receipts(received_by);
CREATE INDEX IF NOT EXISTS idx_goods_receipts_status ON goods_receipts(receipt_status);

-- Goods receipt details table indexes
CREATE INDEX IF NOT EXISTS idx_goods_receipt_details_receipt_id ON goods_receipt_details(receipt_id);
CREATE INDEX IF NOT EXISTS idx_goods_receipt_details_po_detail_id ON goods_receipt_details(po_detail_id);
CREATE INDEX IF NOT EXISTS idx_goods_receipt_details_product_id ON goods_receipt_details(product_id);
CREATE INDEX IF NOT EXISTS idx_goods_receipt_details_condition ON goods_receipt_details(condition_received);

-- Stock movements table indexes
CREATE INDEX IF NOT EXISTS idx_stock_movements_product_id ON stock_movements(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_type ON stock_movements(movement_type);
CREATE INDEX IF NOT EXISTS idx_stock_movements_reference_type ON stock_movements(reference_type);
CREATE INDEX IF NOT EXISTS idx_stock_movements_reference_id ON stock_movements(reference_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_date ON stock_movements(movement_date);
CREATE INDEX IF NOT EXISTS idx_stock_movements_processed_by ON stock_movements(processed_by);

-- Stock adjustments table indexes
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_product_id ON stock_adjustments(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_type ON stock_adjustments(adjustment_type);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_date ON stock_adjustments(adjustment_date);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_approved_by ON stock_adjustments(approved_by);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_created_by ON stock_adjustments(created_by);

-- Supplier payments table indexes
CREATE INDEX IF NOT EXISTS idx_supplier_payments_number ON supplier_payments(payment_number);
CREATE INDEX IF NOT EXISTS idx_supplier_payments_supplier_id ON supplier_payments(supplier_id);
CREATE INDEX IF NOT EXISTS idx_supplier_payments_po_id ON supplier_payments(po_id);
CREATE INDEX IF NOT EXISTS idx_supplier_payments_status ON supplier_payments(payment_status);
CREATE INDEX IF NOT EXISTS idx_supplier_payments_invoice_date ON supplier_payments(invoice_date);
CREATE INDEX IF NOT EXISTS idx_supplier_payments_payment_date ON supplier_payments(payment_date);
CREATE INDEX IF NOT EXISTS idx_supplier_payments_due_date ON supplier_payments(due_date);
CREATE INDEX IF NOT EXISTS idx_supplier_payments_method ON supplier_payments(payment_method);
CREATE INDEX IF NOT EXISTS idx_supplier_payments_processed_by ON supplier_payments(processed_by);
CREATE INDEX IF NOT EXISTS idx_supplier_payments_invoice_number ON supplier_payments(invoice_number);`