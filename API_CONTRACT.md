# API Contract & Schema Documentation

## Phase 4: Vehicle Purchase Flow APIs

### Vehicle Purchase Transactions

#### Data Schema
```typescript
interface VehiclePurchaseTransaction {
  transaction_id: number;
  transaction_number: string;
  customer_id: number;
  vehicle_id?: number;
  vin_number?: string;
  vehicle_brand: string;
  vehicle_model: string;
  vehicle_year: number;
  vehicle_color: string;
  engine_number?: string;
  registration_number?: string;
  purchase_price: number;
  agreed_value: number;
  odometer_reading: number;
  fuel_type: string;
  transmission: string;
  condition_rating?: number; // 1-10
  purchase_date: string; // ISO 8601
  transaction_status: 'pending' | 'inspection' | 'approved' | 'rejected' | 'completed' | 'cancelled';
  inspection_notes?: string;
  evaluation_notes?: string;
  purchase_notes?: string;
  documents_json?: string;
  processed_by: number;
  inspected_by?: number;
  approved_by?: number;
  approved_at?: string; // ISO 8601
  created_at: string; // ISO 8601
  updated_at: string; // ISO 8601
  
  // Related data
  customer_name?: string;
  processed_by_name?: string;
  inspected_by_name?: string;
  approved_by_name?: string;
}
```

#### API Endpoints

**POST /api/v1/vehicle-purchases/transactions**
```typescript
// Request
interface CreateTransactionRequest {
  customer_id: number;
  vehicle_id?: number;
  vin_number?: string; // max 50 chars
  vehicle_brand: string; // max 100 chars, required
  vehicle_model: string; // max 100 chars, required
  vehicle_year: number; // 1900-2100, required
  vehicle_color: string; // max 50 chars, required
  engine_number?: string; // max 100 chars
  registration_number?: string; // max 50 chars
  purchase_price: number; // min 0, required
  agreed_value: number; // min 0, required
  odometer_reading: number; // min 0
  fuel_type: string; // max 50 chars, required
  transmission: string; // max 50 chars, required
  condition_rating?: number; // 1-10
  inspection_notes?: string;
  evaluation_notes?: string;
  purchase_notes?: string;
  documents_json?: string;
}

// Response
interface CreateTransactionResponse {
  success: boolean;
  data: VehiclePurchaseTransaction;
  message: string;
}
```

**GET /api/v1/vehicle-purchases/transactions**
```typescript
// Query Parameters
interface TransactionFilterParams {
  page?: number; // min 1, default 1
  limit?: number; // min 1, max 100, default 10
  customer_id?: number;
  vehicle_brand?: string;
  vehicle_model?: string;
  vehicle_year?: number;
  transaction_status?: string;
  min_purchase_price?: number;
  max_purchase_price?: number;
  processed_by?: number;
  inspected_by?: number;
  start_date?: string; // ISO 8601
  end_date?: string; // ISO 8601
  search?: string;
}

// Response
interface ListTransactionsResponse {
  success: boolean;
  data: VehiclePurchaseTransaction[];
  pagination: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
    has_more: boolean;
  };
  message: string;
}
```

**POST /api/v1/vehicle-purchases/transactions/:id/inspect**
```typescript
// Request
interface InspectionRequest {
  condition_rating: number; // 1-10, required
  inspection_notes: string; // required
  evaluation_notes?: string;
  recommended_action: 'approve' | 'reject' | 'needs_repair'; // required
}

// Response
interface InspectionResponse {
  success: boolean;
  data: VehiclePurchaseTransaction;
  message: string;
}
```

### Vehicle Purchase Payments

#### Data Schema
```typescript
interface VehiclePurchasePayment {
  payment_id: number;
  transaction_id: number;
  payment_number: string;
  payment_method: 'cash' | 'transfer' | 'check' | 'financing';
  payment_amount: number;
  payment_date: string; // ISO 8601
  payment_status: 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';
  reference_number?: string;
  bank_account?: string;
  payment_description?: string;
  payment_notes?: string;
  processed_by: number;
  approved_by?: number;
  approved_at?: string; // ISO 8601
  created_at: string; // ISO 8601
  updated_at: string; // ISO 8601
  
  // Related data
  transaction_number?: string;
  customer_name?: string;
  vehicle_brand?: string;
  vehicle_model?: string;
  processed_by_name?: string;
  approved_by_name?: string;
}
```

## Phase 5: Repair Management APIs

### Vehicle Damages

#### Data Schema
```typescript
interface VehicleDamage {
  damage_id: number;
  transaction_id: number;
  damage_category: 'body' | 'engine' | 'interior' | 'electrical' | 'suspension' | 'brake' | 'transmission' | 'other';
  damage_type: string; // max 100 chars
  damage_severity: 'minor' | 'moderate' | 'major' | 'critical';
  damage_location: string; // max 100 chars
  damage_description: string;
  estimated_cost: number; // min 0
  repair_priority: number; // 1-5
  repair_required: boolean;
  damage_photos_json?: string;
  assessment_notes?: string;
  identified_by: number;
  identified_at: string; // ISO 8601
  status: 'identified' | 'assessed' | 'scheduled' | 'repairing' | 'completed' | 'cancelled';
  created_at: string; // ISO 8601
  updated_at: string; // ISO 8601
  
  // Related data
  transaction_number?: string;
  vehicle_brand?: string;
  vehicle_model?: string;
  identified_by_name?: string;
}
```

### Repair Work Orders

#### Data Schema
```typescript
interface RepairWorkOrder {
  work_order_id: number;
  work_order_number: string;
  transaction_id: number;
  work_order_type: 'inspection' | 'repair' | 'maintenance' | 'improvement';
  work_order_priority: number; // 1-5
  scheduled_start_date?: string; // ISO 8601
  scheduled_end_date?: string; // ISO 8601
  actual_start_date?: string; // ISO 8601
  actual_end_date?: string; // ISO 8601
  estimated_cost: number; // min 0
  actual_cost: number; // min 0
  labor_hours_estimated: number; // min 0
  labor_hours_actual: number; // min 0
  work_order_status: 'draft' | 'scheduled' | 'in_progress' | 'suspended' | 'completed' | 'cancelled';
  work_description: string;
  special_instructions?: string;
  completion_notes?: string;
  assigned_mechanic_id?: number;
  supervisor_id?: number;
  created_by: number;
  approved_by?: number;
  approved_at?: string; // ISO 8601
  created_at: string; // ISO 8601
  updated_at: string; // ISO 8601
  
  // Related data
  transaction_number?: string;
  vehicle_brand?: string;
  vehicle_model?: string;
  assigned_mechanic_name?: string;
  supervisor_name?: string;
  created_by_name?: string;
  approved_by_name?: string;
}
```

### Repair Work Details

#### Data Schema
```typescript
interface RepairWorkDetail {
  work_detail_id: number;
  work_order_id: number;
  damage_id?: number;
  task_sequence: number; // min 1
  task_description: string;
  task_type: 'diagnosis' | 'disassembly' | 'repair' | 'replacement' | 'assembly' | 'testing' | 'quality_check';
  estimated_hours: number; // min 0
  actual_hours: number; // min 0
  labor_rate: number; // min 0
  task_status: 'pending' | 'in_progress' | 'completed' | 'cancelled' | 'on_hold';
  start_date?: string; // ISO 8601
  end_date?: string; // ISO 8601
  completion_percentage: number; // 0-100
  task_notes?: string;
  quality_check_passed?: boolean;
  assigned_mechanic_id?: number;
  verified_by?: number;
  verified_at?: string; // ISO 8601
  created_at: string; // ISO 8601
  updated_at: string; // ISO 8601
  
  // Related data
  work_order_number?: string;
  damage_description?: string;
  assigned_mechanic_name?: string;
  verified_by_name?: string;
}
```

### Repair Parts Usage

#### Data Schema
```typescript
interface RepairPartsUsage {
  usage_id: number;
  work_detail_id: number;
  product_id: number;
  quantity_used: number; // min 1
  unit_cost: number; // min 0
  total_cost: number; // min 0
  usage_date: string; // ISO 8601
  usage_type: 'new' | 'replacement' | 'additional' | 'warranty';
  part_condition: 'new' | 'refurbished' | 'used' | 'oem' | 'aftermarket';
  warranty_period_days: number; // min 0
  installation_notes?: string;
  issued_by: number;
  used_by: number;
  approved_by?: number;
  approved_at?: string; // ISO 8601
  created_at: string; // ISO 8601
  
  // Related data
  product_name?: string;
  product_code?: string;
  work_order_number?: string;
  task_description?: string;
  issued_by_name?: string;
  used_by_name?: string;
  approved_by_name?: string;
}
```

### Quality Inspections

#### Data Schema
```typescript
interface QualityInspection {
  inspection_id: number;
  work_order_id: number;
  inspection_type: 'pre_repair' | 'during_repair' | 'post_repair' | 'final_inspection';
  inspection_date: string; // ISO 8601
  inspector_id: number;
  overall_rating: number; // 1-10
  workmanship_rating: number; // 1-10
  safety_rating: number; // 1-10
  appearance_rating: number; // 1-10
  functionality_rating: number; // 1-10
  inspection_status: 'passed' | 'failed' | 'conditional_pass' | 'needs_rework';
  inspection_notes?: string;
  defects_found?: string;
  recommendations?: string;
  photos_json?: string;
  signed_off_by?: number;
  signed_off_at?: string; // ISO 8601
  rework_required: boolean;
  next_inspection_date?: string; // ISO 8601
  created_at: string; // ISO 8601
  updated_at: string; // ISO 8601
  
  // Related data
  work_order_number?: string;
  transaction_number?: string;
  vehicle_brand?: string;
  vehicle_model?: string;
  inspector_name?: string;
  signed_off_by_name?: string;
}
```

## Common Response Formats

### Success Response
```typescript
interface SuccessResponse<T> {
  success: true;
  data: T;
  message: string;
}
```

### Error Response
```typescript
interface ErrorResponse {
  success: false;
  error: {
    code: string;
    message: string;
    details?: any;
  };
  message: string;
}
```

### Paginated Response
```typescript
interface PaginatedResponse<T> {
  success: boolean;
  data: T[];
  pagination: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
    has_more: boolean;
  };
  message: string;
}
```

## HTTP Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict (duplicate data)
- `422` - Unprocessable Entity (business rule violations)
- `500` - Internal Server Error

## Authentication & Authorization

### Headers Required
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Role-Based Access
- **Admin**: Full access to all endpoints
- **Manager**: Read/write access, approval operations
- **Sales**: Vehicle purchase transactions
- **Mechanic**: Repair operations, parts usage
- **Cashier**: Payment processing
- **Inspector**: Quality inspections

## Rate Limiting
- 1000 requests per hour per user
- Burst limit: 50 requests per minute

## Validation Rules

### Common Validations
- All IDs must be positive integers
- Dates must be in ISO 8601 format
- Decimal values have 2 decimal places max
- Text fields have specified max lengths
- Enum fields must match predefined values

### Business Rule Validations
- VIN numbers must be unique
- Transaction numbers auto-generated and unique
- Payment amounts cannot exceed transaction value
- Quality ratings must be 1-10
- Condition ratings must be 1-10
- Priority levels must be 1-5

This API contract ensures consistent data exchange and proper validation across all Phase 4 and Phase 5 operations.