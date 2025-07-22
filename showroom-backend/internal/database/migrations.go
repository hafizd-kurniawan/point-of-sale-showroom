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

-- Customers table indexes
CREATE INDEX IF NOT EXISTS idx_customers_code ON customers(customer_code);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(customer_name);
CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_type ON customers(customer_type);
CREATE INDEX IF NOT EXISTS idx_customers_city ON customers(city);
CREATE INDEX IF NOT EXISTS idx_customers_province ON customers(province);
CREATE INDEX IF NOT EXISTS idx_customers_is_active ON customers(is_active);
CREATE INDEX IF NOT EXISTS idx_customers_created_by ON customers(created_by);

-- Suppliers table indexes
CREATE INDEX IF NOT EXISTS idx_suppliers_code ON suppliers(supplier_code);
CREATE INDEX IF NOT EXISTS idx_suppliers_name ON suppliers(supplier_name);
CREATE INDEX IF NOT EXISTS idx_suppliers_phone ON suppliers(phone);
CREATE INDEX IF NOT EXISTS idx_suppliers_email ON suppliers(email);
CREATE INDEX IF NOT EXISTS idx_suppliers_type ON suppliers(supplier_type);
CREATE INDEX IF NOT EXISTS idx_suppliers_city ON suppliers(city);
CREATE INDEX IF NOT EXISTS idx_suppliers_province ON suppliers(province);
CREATE INDEX IF NOT EXISTS idx_suppliers_is_active ON suppliers(is_active);
CREATE INDEX IF NOT EXISTS idx_suppliers_created_by ON suppliers(created_by);

-- Vehicle brands table indexes
CREATE INDEX IF NOT EXISTS idx_vehicle_brands_code ON vehicle_brands(brand_code);
CREATE INDEX IF NOT EXISTS idx_vehicle_brands_name ON vehicle_brands(brand_name);
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
CREATE INDEX IF NOT EXISTS idx_vehicle_models_fuel_type ON vehicle_models(fuel_type);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_transmission ON vehicle_models(transmission);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_year_start ON vehicle_models(year_start);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_is_active ON vehicle_models(is_active);
CREATE INDEX IF NOT EXISTS idx_vehicle_models_created_by ON vehicle_models(created_by);

-- Product categories table indexes
CREATE INDEX IF NOT EXISTS idx_product_categories_code ON product_categories(category_code);
CREATE INDEX IF NOT EXISTS idx_product_categories_name ON product_categories(category_name);
CREATE INDEX IF NOT EXISTS idx_product_categories_is_active ON product_categories(is_active);
CREATE INDEX IF NOT EXISTS idx_product_categories_created_by ON product_categories(created_by);

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
    customer_code VARCHAR(50) UNIQUE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    id_card_number VARCHAR(20) UNIQUE,
    tax_number VARCHAR(50),
    customer_type VARCHAR(20) NOT NULL CHECK (customer_type IN ('individual','corporate')),
    birth_date TIMESTAMP,
    occupation VARCHAR(100),
    income_range DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id),
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT
);`

const createSuppliersTable = `
CREATE TABLE IF NOT EXISTS suppliers (
    supplier_id SERIAL PRIMARY KEY,
    supplier_code VARCHAR(50) UNIQUE NOT NULL,
    supplier_name VARCHAR(255) NOT NULL,
    contact_person VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    tax_number VARCHAR(50),
    supplier_type VARCHAR(20) NOT NULL CHECK (supplier_type IN ('parts','vehicle','both')),
    credit_limit DECIMAL(15,2),
    credit_term_days INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id),
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT
);`

const createVehicleBrandsTable = `
CREATE TABLE IF NOT EXISTS vehicle_brands (
    brand_id SERIAL PRIMARY KEY,
    brand_code VARCHAR(50) UNIQUE NOT NULL,
    brand_name VARCHAR(255) NOT NULL,
    country_origin VARCHAR(100),
    logo_image VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createVehicleCategoriesTable = `
CREATE TABLE IF NOT EXISTS vehicle_categories (
    category_id SERIAL PRIMARY KEY,
    category_code VARCHAR(50) UNIQUE NOT NULL,
    category_name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createVehicleModelsTable = `
CREATE TABLE IF NOT EXISTS vehicle_models (
    model_id SERIAL PRIMARY KEY,
    brand_id INTEGER NOT NULL REFERENCES vehicle_brands(brand_id),
    category_id INTEGER NOT NULL REFERENCES vehicle_categories(category_id),
    model_code VARCHAR(50) UNIQUE NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    year_start INTEGER NOT NULL,
    year_end INTEGER,
    fuel_type VARCHAR(20) NOT NULL CHECK (fuel_type IN ('gasoline','diesel','electric','hybrid')),
    transmission VARCHAR(20) NOT NULL CHECK (transmission IN ('manual','automatic','cvt')),
    engine_capacity INTEGER NOT NULL,
    specifications_json JSON,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`

const createProductCategoriesTable = `
CREATE TABLE IF NOT EXISTS product_categories (
    category_id SERIAL PRIMARY KEY,
    category_code VARCHAR(50) UNIQUE NOT NULL,
    category_name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES users(user_id)
);`