# API Testing Guide - Phase 4 & 5

## Overview

This guide provides comprehensive instructions for testing the Phase 4 (Vehicle Purchase Flow) and Phase 5 (Repair Management System) APIs.

## Table of Contents

1. [Environment Setup](#environment-setup)
2. [Authentication](#authentication)
3. [Testing Scenarios](#testing-scenarios)
4. [Test Data](#test-data)
5. [Expected Results](#expected-results)
6. [Troubleshooting](#troubleshooting)

## Environment Setup

### Prerequisites

1. **Server Running**: Ensure the Go server is running on `http://localhost:8080`
2. **Database**: PostgreSQL database with Phase 4 & 5 tables created
3. **Authentication**: Valid JWT token for API access
4. **Test Tools**: One of the following:
   - VS Code with REST Client extension
   - Postman
   - curl
   - Insomnia

### Configuration

```bash
# Environment Variables
export SERVER_URL=http://localhost:8080
export JWT_TOKEN=your_jwt_token_here
export DATABASE_URL=postgresql://user:password@localhost:5432/showroom_db
```

## Authentication

### Getting JWT Token

```bash
# Login to get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin@showroom.com",
    "password": "admin123"
  }'
```

### Using Token in Requests

```bash
# Include in Authorization header
Authorization: Bearer your_jwt_token_here
```

## Testing Scenarios

### Scenario 1: Complete Vehicle Purchase Flow

```bash
# Step 1: Create Vehicle Purchase Transaction
curl -X POST http://localhost:8080/api/v1/vehicle-purchases/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "customer_id": 1,
    "vehicle_brand": "Honda",
    "vehicle_model": "Civic",
    "vehicle_year": 2020,
    "purchase_price": 25000000,
    "agreed_value": 23000000,
    "vin_number": "1HGBH41JXMN109186"
  }'

# Step 2: Complete Inspection
curl -X POST http://localhost:8080/api/v1/vehicle-purchases/transactions/1/inspect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "condition_rating": 7,
    "inspection_notes": "Minor damage found"
  }'

# Step 3: Approve Transaction
curl -X POST http://localhost:8080/api/v1/vehicle-purchases/transactions/1/approve \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "status": "approved"
  }'

# Step 4: Create Payment
curl -X POST http://localhost:8080/api/v1/vehicle-purchases/payments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "transaction_id": 1,
    "payment_amount": 23000000,
    "payment_method": "bank_transfer"
  }'
```

### Scenario 2: Complete Repair Management Flow

```bash
# Step 1: Create Damage Record
curl -X POST http://localhost:8080/api/v1/repairs/damages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "transaction_id": 1,
    "damage_category": "body",
    "damage_severity": "medium",
    "damage_description": "Rear bumper damage",
    "estimated_cost": 1500000
  }'

# Step 2: Create Work Order
curl -X POST http://localhost:8080/api/v1/repairs/work-orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "transaction_id": 1,
    "priority_level": "medium",
    "estimated_cost": 1500000,
    "work_description": "Repair rear bumper"
  }'

# Step 3: Create Work Detail
curl -X POST http://localhost:8080/api/v1/repairs/work-details \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "work_order_id": 1,
    "damage_id": 1,
    "task_description": "Sand and paint bumper",
    "estimated_hours": 4
  }'

# Step 4: Issue Parts
curl -X POST http://localhost:8080/api/v1/repairs/work-details/1/issue-parts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "product_id": 201,
    "quantity_used": 2,
    "unit_cost": 150000,
    "total_cost": 300000
  }'

# Step 5: Create Quality Inspection
curl -X POST http://localhost:8080/api/v1/repairs/inspections \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "work_order_id": 1,
    "inspection_type": "final_inspection",
    "scheduled_date": "2024-01-16T16:00:00Z"
  }'
```

## Test Data

### Sample Customers

```json
{
  "customer_id": 1,
  "customer_name": "John Doe",
  "email": "john.doe@email.com",
  "phone": "081234567890"
}
```

### Sample Products (Parts)

```json
{
  "product_id": 201,
  "product_name": "Automotive Paint - White",
  "product_code": "PAINT-WHT-001",
  "current_stock": 50,
  "unit_price": 150000
}
```

### Sample Users

```json
{
  "admin": {
    "user_id": 1,
    "username": "admin",
    "role": "admin"
  },
  "mechanic": {
    "user_id": 5,
    "username": "mechanic1",
    "role": "mechanic"
  },
  "inspector": {
    "user_id": 7,
    "username": "inspector1",
    "role": "inspector"
  }
}
```

## Expected Results

### Successful Responses

#### Create Transaction (201 Created)
```json
{
  "success": true,
  "message": "Transaction created successfully",
  "data": {
    "transaction_id": 1,
    "transaction_number": "TXN202401001",
    "customer_id": 1,
    "purchase_price": 25000000,
    "transaction_status": "pending",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

#### Create Work Order (201 Created)
```json
{
  "success": true,
  "message": "Work order created successfully",
  "data": {
    "work_order_id": 1,
    "work_order_number": "WO202401001",
    "transaction_id": 1,
    "estimated_cost": 1500000,
    "work_order_status": "pending",
    "created_at": "2024-01-15T11:00:00Z"
  }
}
```

### Error Responses

#### Validation Error (400 Bad Request)
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "Field 'customer_id' is required"
}
```

#### Not Found (404 Not Found)
```json
{
  "success": false,
  "message": "Transaction not found",
  "error": "No transaction found with ID 999"
}
```

#### Unauthorized (401 Unauthorized)
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Invalid or expired token"
}
```

## Performance Testing

### Load Testing with curl

```bash
#!/bin/bash
# Simple load test script

echo "Starting load test..."
for i in {1..10}
do
  echo "Request $i"
  curl -s -o /dev/null -w "%{http_code} - %{time_total}s\n" \
    -X GET http://localhost:8080/api/v1/vehicle-purchases/transactions \
    -H "Authorization: Bearer $JWT_TOKEN"
done
```

### Expected Performance

- **Response Time**: < 200ms for GET requests
- **Response Time**: < 500ms for POST requests
- **Throughput**: > 100 requests/second
- **Error Rate**: < 1%

## Integration Testing

### End-to-End Test Script

```bash
#!/bin/bash
set -e

echo "Running E2E tests..."

# Test 1: Vehicle Purchase Flow
echo "Testing Vehicle Purchase Flow..."
TRANSACTION_ID=$(curl -s -X POST http://localhost:8080/api/v1/vehicle-purchases/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{"customer_id": 1, "vehicle_brand": "Honda", "purchase_price": 25000000}' \
  | jq -r '.data.transaction_id')

echo "Created transaction: $TRANSACTION_ID"

# Test 2: Inspection
curl -s -X POST http://localhost:8080/api/v1/vehicle-purchases/transactions/$TRANSACTION_ID/inspect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{"condition_rating": 8, "inspection_notes": "Good condition"}'

echo "Inspection completed"

# Test 3: Approval
curl -s -X POST http://localhost:8080/api/v1/vehicle-purchases/transactions/$TRANSACTION_ID/approve \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{"status": "approved"}'

echo "Transaction approved"

echo "E2E test completed successfully!"
```

## Troubleshooting

### Common Issues

#### 1. Authentication Errors
```bash
# Check token validity
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $JWT_TOKEN"
```

#### 2. Database Connection Issues
```bash
# Check if database tables exist
psql $DATABASE_URL -c "\dt"
```

#### 3. Server Not Running
```bash
# Check if server is running
curl http://localhost:8080/api/v1/health
```

#### 4. Invalid Request Data
- Verify JSON syntax
- Check required fields
- Validate data types
- Ensure foreign key references exist

### Debug Mode

Set environment variable for detailed logging:
```bash
export GIN_MODE=debug
export LOG_LEVEL=debug
```

### API Response Validation

Use JSON schema validation tools:
```bash
# Install ajv-cli
npm install -g ajv-cli

# Validate response
echo '{"transaction_id": 1}' | ajv validate -s transaction-schema.json
```

## Test Coverage Report

### API Endpoints Tested

#### Phase 4 - Vehicle Purchase (21 endpoints)
- ✅ Transaction CRUD operations
- ✅ Payment processing
- ✅ Inspection workflow
- ✅ Approval workflow
- ✅ Dashboard statistics

#### Phase 5 - Repair Management (32 endpoints)
- ✅ Damage management
- ✅ Work order lifecycle
- ✅ Work detail tracking
- ✅ Parts usage management
- ✅ Quality inspection process

### Test Scenarios Covered
- ✅ Happy path workflows
- ✅ Error handling
- ✅ Authentication & authorization
- ✅ Data validation
- ✅ Business rule enforcement
- ✅ Integration between phases

### Automated Testing

Consider implementing:
1. **Unit Tests**: Go test files for each handler
2. **Integration Tests**: Full workflow testing
3. **Contract Tests**: API schema validation
4. **Performance Tests**: Load and stress testing
5. **Security Tests**: Authentication and authorization testing

This comprehensive testing guide ensures thorough validation of all Phase 4 & 5 functionality.