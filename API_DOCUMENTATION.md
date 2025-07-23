# Phase 4 & 5 API Documentation

## Phase 4: Vehicle Purchase Flow

### Overview
The Vehicle Purchase Flow manages the process of purchasing vehicles from customers, including inspection, damage assessment, and payment processing.

### Workflow
1. **Vehicle Purchase Transaction Creation**: Create a new transaction when a customer brings a vehicle to sell
2. **Vehicle Inspection**: Inspect the vehicle and record its condition 
3. **Damage Assessment**: Identify and record any damages found during inspection
4. **Transaction Approval**: Approve or reject the purchase based on inspection results
5. **Payment Processing**: Process payments to the customer for approved purchases

### API Endpoints

#### Vehicle Purchase Transactions

```
POST   /api/v1/vehicle-purchases/transactions          # Create new purchase transaction
GET    /api/v1/vehicle-purchases/transactions          # List all transactions with filtering
GET    /api/v1/vehicle-purchases/transactions/:id      # Get transaction details
PUT    /api/v1/vehicle-purchases/transactions/:id      # Update transaction
DELETE /api/v1/vehicle-purchases/transactions/:id      # Delete transaction
POST   /api/v1/vehicle-purchases/transactions/:id/inspect   # Complete vehicle inspection  
POST   /api/v1/vehicle-purchases/transactions/:id/approve   # Approve/reject transaction
GET    /api/v1/vehicle-purchases/transactions/pending-inspection  # Get pending inspections
GET    /api/v1/vehicle-purchases/transactions/pending-approval    # Get pending approvals
GET    /api/v1/vehicle-purchases/transactions/dashboard           # Dashboard statistics
```

#### Vehicle Purchase Payments

```
POST   /api/v1/vehicle-purchases/payments             # Create new payment
GET    /api/v1/vehicle-purchases/payments             # List all payments with filtering
GET    /api/v1/vehicle-purchases/payments/:id         # Get payment details
PUT    /api/v1/vehicle-purchases/payments/:id         # Update payment
DELETE /api/v1/vehicle-purchases/payments/:id         # Delete payment
POST   /api/v1/vehicle-purchases/payments/:id/process # Process payment
POST   /api/v1/vehicle-purchases/payments/:id/approve # Approve/reject payment
GET    /api/v1/vehicle-purchases/payments/summary/:transaction_id  # Payment summary for transaction
```

### Example API Requests

#### Create Vehicle Purchase Transaction
```json
POST /api/v1/vehicle-purchases/transactions
{
  "customer_id": 1,
  "vin_number": "1HGBH41JXMN109186",
  "vehicle_brand": "Honda",
  "vehicle_model": "Civic",
  "vehicle_year": 2020,
  "vehicle_color": "Blue",
  "engine_number": "ENG12345",
  "registration_number": "ABC123",
  "purchase_price": 15000.00,
  "agreed_value": 15000.00,
  "odometer_reading": 45000,
  "fuel_type": "Gasoline",
  "transmission": "Automatic",
  "purchase_notes": "Vehicle in good condition"
}
```

#### Complete Vehicle Inspection
```json
POST /api/v1/vehicle-purchases/transactions/1/inspect
{
  "condition_rating": 8,
  "inspection_notes": "Minor scratches on rear bumper, otherwise excellent condition",
  "evaluation_notes": "Recommended for purchase",
  "recommended_action": "approve"
}
```

## Phase 5: Repair Management

### Overview
The Repair Management system handles the repair workflow for purchased vehicles, including damage tracking, work order management, parts usage, and quality inspections.

### Workflow
1. **Damage Identification**: Record damages found during vehicle inspection
2. **Work Order Creation**: Create repair work orders for identified damages
3. **Work Detail Planning**: Break down work orders into specific tasks
4. **Parts Issue**: Issue required parts for repair work
5. **Repair Execution**: Mechanics perform repair tasks
6. **Quality Inspection**: Inspect completed repairs for quality assurance

### API Endpoints

#### Vehicle Damages

```
POST   /api/v1/repairs/damages                        # Record new damage
GET    /api/v1/repairs/damages                        # List all damages with filtering
GET    /api/v1/repairs/damages/:id                    # Get damage details
PUT    /api/v1/repairs/damages/:id                    # Update damage
DELETE /api/v1/repairs/damages/:id                    # Delete damage
POST   /api/v1/repairs/damages/:id/assess             # Assess damage cost and priority
GET    /api/v1/repairs/damages/transaction/:id        # Get damages for transaction
GET    /api/v1/repairs/damages/high-priority          # Get high priority damages
GET    /api/v1/repairs/damages/summary/:transaction_id # Damage summary for transaction
```

#### Repair Work Orders

```
POST   /api/v1/repairs/work-orders                    # Create new work order
GET    /api/v1/repairs/work-orders                    # List all work orders with filtering
GET    /api/v1/repairs/work-orders/:id                # Get work order details
PUT    /api/v1/repairs/work-orders/:id                # Update work order
DELETE /api/v1/repairs/work-orders/:id                # Delete work order
POST   /api/v1/repairs/work-orders/:id/assign         # Assign work order to mechanic
POST   /api/v1/repairs/work-orders/:id/approve        # Approve/reject work order
POST   /api/v1/repairs/work-orders/:id/status         # Update work order status
GET    /api/v1/repairs/work-orders/pending-approval   # Get pending approvals
GET    /api/v1/repairs/work-orders/mechanic/:id       # Get work orders for mechanic
```

#### Repair Work Details

```
POST   /api/v1/repairs/work-details                   # Create new work detail
GET    /api/v1/repairs/work-details                   # List all work details with filtering
GET    /api/v1/repairs/work-details/:id               # Get work detail details
PUT    /api/v1/repairs/work-details/:id               # Update work detail
DELETE /api/v1/repairs/work-details/:id               # Delete work detail
POST   /api/v1/repairs/work-details/:id/progress      # Update work progress
POST   /api/v1/repairs/work-details/:id/assign        # Assign task to mechanic
POST   /api/v1/repairs/work-details/:id/quality-check # Perform quality check
GET    /api/v1/repairs/work-details/work-order/:id    # Get details for work order
```

#### Repair Parts Usage

```
POST   /api/v1/repairs/parts-usage                    # Record parts usage
GET    /api/v1/repairs/parts-usage                    # List all parts usage with filtering
GET    /api/v1/repairs/parts-usage/:id                # Get parts usage details
PUT    /api/v1/repairs/parts-usage/:id                # Update parts usage
DELETE /api/v1/repairs/parts-usage/:id                # Delete parts usage
POST   /api/v1/repairs/parts-usage/:id/approve        # Approve/reject parts usage
POST   /api/v1/repairs/parts-usage/issue              # Issue parts for repair
GET    /api/v1/repairs/parts-usage/work-detail/:id    # Get parts usage for work detail
GET    /api/v1/repairs/parts-usage/summary/:work_order_id # Parts usage summary
```

#### Quality Inspections

```
POST   /api/v1/repairs/inspections                    # Create new inspection
GET    /api/v1/repairs/inspections                    # List all inspections with filtering
GET    /api/v1/repairs/inspections/:id                # Get inspection details
PUT    /api/v1/repairs/inspections/:id                # Update inspection
DELETE /api/v1/repairs/inspections/:id                # Delete inspection
POST   /api/v1/repairs/inspections/:id/sign-off       # Sign off inspection
POST   /api/v1/repairs/inspections/:id/rework         # Schedule rework
POST   /api/v1/repairs/inspections/schedule           # Schedule new inspection
GET    /api/v1/repairs/inspections/work-order/:id     # Get inspections for work order
GET    /api/v1/repairs/inspections/failed             # Get failed inspections
GET    /api/v1/repairs/inspections/dashboard          # Inspection dashboard
```

### Example API Requests

#### Record Vehicle Damage
```json
POST /api/v1/repairs/damages
{
  "transaction_id": 1,
  "damage_category": "body",
  "damage_type": "Scratch",
  "damage_severity": "minor",
  "damage_location": "Rear bumper",
  "damage_description": "Light scratch approximately 6 inches long on rear bumper",
  "estimated_cost": 150.00,
  "repair_priority": 2,
  "repair_required": true,
  "assessment_notes": "Can be buffed and polished"
}
```

#### Create Repair Work Order
```json
POST /api/v1/repairs/work-orders
{
  "transaction_id": 1,
  "work_order_type": "repair",
  "work_order_priority": 3,
  "scheduled_start_date": "2024-07-25T08:00:00Z",
  "scheduled_end_date": "2024-07-25T17:00:00Z",
  "estimated_cost": 500.00,
  "labor_hours_estimated": 8.0,
  "work_description": "Repair body damages and perform cosmetic touch-ups",
  "assigned_mechanic_id": 5,
  "supervisor_id": 3
}
```

#### Record Quality Inspection
```json
POST /api/v1/repairs/inspections
{
  "work_order_id": 1,
  "inspection_type": "post_repair",
  "overall_rating": 9,
  "workmanship_rating": 9,
  "safety_rating": 10,
  "appearance_rating": 8,
  "functionality_rating": 9,
  "inspection_notes": "Excellent repair work, minor color matching issue",
  "rework_required": false
}
```

## Integration Points

### Stock Management Integration
- Parts usage automatically updates inventory levels
- Low stock alerts triggered when parts usage reduces inventory below minimum levels
- Integration with existing purchase order system for parts replenishment

### User Role Integration
- **Admin**: Full access to all operations
- **Manager**: Approve transactions, work orders, and high-value operations
- **Sales**: Create and manage purchase transactions
- **Mechanic**: Perform repairs, update work progress, record parts usage
- **Cashier**: Process payments

### Workflow States

#### Purchase Transaction States
- `pending` → `inspection` → `approved`/`rejected` → `completed`

#### Work Order States  
- `draft` → `scheduled` → `in_progress` → `completed`

#### Quality Inspection States
- `passed` → `conditional_pass` → `failed` → `needs_rework`

## Business Rules

### Purchase Flow
1. Vehicle inspection must be completed before approval
2. Only approved transactions can proceed to payment
3. VIN numbers must be unique across all transactions
4. Condition rating affects purchase price recommendations

### Repair Flow
1. Work orders can only be created for approved purchase transactions
2. High-priority damages must be addressed first
3. Quality inspections are mandatory for completed work orders
4. Parts usage requires approval for high-value items
5. Rework automatically creates new work orders

### Payment Processing
1. Payment amounts cannot exceed agreed purchase value
2. Multiple partial payments allowed
3. Payment approval required for amounts above threshold
4. Automatic payment reconciliation with transaction totals