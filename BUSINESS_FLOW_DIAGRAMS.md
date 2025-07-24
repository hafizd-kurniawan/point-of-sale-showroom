# Business Flow Diagrams - Phase 4 & 5: Vehicle Purchase and Repair Management

## Overview

This document provides comprehensive business flow diagrams for the Point of Sale Showroom system, specifically covering Phase 4 (Vehicle Purchase Flow) and Phase 5 (Repair Management System).

## Table of Contents

1. [Phase 4: Vehicle Purchase Flow](#phase-4-vehicle-purchase-flow)
2. [Phase 5: Repair Management Flow](#phase-5-repair-management-flow)
3. [Integrated End-to-End Process](#integrated-end-to-end-process)
4. [API Endpoint Mapping](#api-endpoint-mapping)

---

## Phase 4: Vehicle Purchase Flow

### High-Level Vehicle Purchase Process

```mermaid
graph TD
    A[Customer brings vehicle] --> B[Initial Vehicle Evaluation]
    B --> C[Create Purchase Transaction]
    C --> D[Vehicle Inspection]
    D --> E{Inspection Results}
    E -->|Pass| F[Transaction Approval]
    E -->|Issues Found| G[Update Condition Rating]
    G --> H[Damage Assessment]
    H --> I[Adjust Purchase Price]
    I --> F
    F --> J[Create Payment Records]
    J --> K[Process Payment]
    K --> L{Payment Status}
    L -->|Completed| M[Transaction Complete]
    L -->|Pending| N[Wait for Payment]
    N --> K
    M --> O[Vehicle Ready for Repair]
    O --> P[Phase 5: Repair Process]

    style A fill:#e1f5fe
    style M fill:#c8e6c9
    style P fill:#fff3e0
```

### Detailed Transaction Management Flow

```mermaid
sequenceDiagram
    participant C as Customer
    participant S as Sales Staff
    participant I as Inspector
    participant M as Manager
    participant SYS as System

    C->>S: Brings vehicle for sale
    S->>SYS: Create transaction (POST /transactions)
    SYS-->>S: Transaction number generated
    
    S->>I: Request vehicle inspection
    I->>SYS: Complete inspection (POST /transactions/{id}/inspect)
    SYS-->>I: Inspection recorded
    
    I->>M: Submit for approval
    M->>SYS: Process approval (POST /transactions/{id}/approve)
    SYS-->>M: Transaction approved
    
    M->>S: Create payment record
    S->>SYS: Create payment (POST /payments)
    SYS-->>S: Payment record created
    
    S->>SYS: Process payment (POST /payments/{id}/process)
    SYS-->>C: Payment confirmation
    
    Note over SYS: Transaction complete, vehicle ready for repair
```

### Payment Processing Flow

```mermaid
stateDiagram-v2
    [*] --> Pending: Create Payment
    Pending --> Processed: Process Payment
    Pending --> Cancelled: Cancel Payment
    Processed --> Approved: Approve Payment
    Processed --> Rejected: Reject Payment
    Approved --> Completed: Finalize
    Rejected --> Pending: Reprocess
    Cancelled --> [*]
    Completed --> [*]
    
    note right of Approved: Payment verified and confirmed
    note right of Completed: Funds transferred successfully
```

---

## Phase 5: Repair Management Flow

### Complete Repair Management Process

```mermaid
graph TD
    A[Vehicle from Purchase] --> B[Damage Assessment]
    B --> C[Create Damage Records]
    C --> D[Estimate Repair Costs]
    D --> E[Create Work Order]
    E --> F[Work Order Approval]
    F --> G[Assign Mechanic]
    G --> H[Break Down into Tasks]
    H --> I[Create Work Details]
    I --> J[Assign Tasks to Mechanics]
    J --> K[Issue Required Parts]
    K --> L[Start Repair Work]
    L --> M[Update Progress]
    M --> N{Task Complete?}
    N -->|No| L
    N -->|Yes| O[Quality Check]
    O --> P{Quality Pass?}
    P -->|No| Q[Schedule Rework]
    Q --> L
    P -->|Yes| R[Mark Task Complete]
    R --> S{All Tasks Done?}
    S -->|No| I
    S -->|Yes| T[Final Inspection]
    T --> U{Final Quality Check}
    U -->|Pass| V[Complete Work Order]
    U -->|Fail| W[Schedule Major Rework]
    W --> I
    V --> X[Vehicle Ready for Sale]

    style A fill:#e1f5fe
    style V fill:#c8e6c9
    style X fill:#4caf50
```

### Detailed Repair Workflow

```mermaid
sequenceDiagram
    participant T as Transaction System
    participant D as Damage Assessor
    participant M as Repair Manager
    participant W as Work Order System
    participant ME as Mechanic
    participant Q as Quality Inspector
    participant P as Parts System

    T->>D: Vehicle ready for repair
    D->>W: Create damage record (POST /damages)
    D->>W: Assess damage (POST /damages/{id}/assess)
    
    M->>W: Create work order (POST /work-orders)
    M->>W: Assign mechanic (POST /work-orders/{id}/assign)
    
    M->>W: Create work details (POST /work-details)
    M->>W: Assign tasks (POST /work-details/{id}/assign)
    
    ME->>P: Request parts (POST /work-details/{id}/issue-parts)
    P-->>ME: Parts issued
    
    loop Repair Process
        ME->>W: Update progress (POST /work-details/{id}/progress)
        ME->>Q: Request quality check
        Q->>W: Perform quality check (POST /work-details/{id}/quality-check)
        
        alt Quality Pass
            Q-->>ME: Task approved
        else Quality Fail
            Q->>W: Schedule rework
            W-->>ME: Rework required
        end
    end
    
    Q->>W: Final inspection (POST /inspections)
    Q->>W: Sign off inspection (POST /inspections/{id}/sign-off)
    W-->>T: Repair complete
```

### Parts Management Flow

```mermaid
graph LR
    A[Work Detail Created] --> B[Identify Required Parts]
    B --> C[Check Inventory]
    C --> D{Parts Available?}
    D -->|Yes| E[Issue Parts]
    D -->|No| F[Order Parts]
    F --> G[Receive Parts]
    G --> E
    E --> H[Update Stock]
    H --> I[Record Usage]
    I --> J[Track Costs]
    J --> K[Complete Task]

    style A fill:#e3f2fd
    style K fill:#e8f5e8
```

### Quality Inspection Workflow

```mermaid
stateDiagram-v2
    [*] --> Scheduled: Create Inspection
    Scheduled --> InProgress: Start Inspection
    InProgress --> Completed: Finish Inspection
    Completed --> Passed: Quality Met
    Completed --> Failed: Quality Not Met
    Failed --> ReworkRequired: Schedule Rework
    ReworkRequired --> Scheduled: New Inspection
    Passed --> SignedOff: Manager Approval
    SignedOff --> [*]
    
    note right of Passed: All quality standards met
    note right of Failed: Issues found, rework needed
    note right of SignedOff: Final approval completed
```

---

## Integrated End-to-End Process

### Complete Vehicle Lifecycle Flow

```mermaid
gantt
    title Vehicle Purchase and Repair Timeline
    dateFormat  YYYY-MM-DD
    section Purchase Phase
    Customer Evaluation    :done, eval, 2024-01-15, 1d
    Vehicle Inspection     :done, inspect, after eval, 1d
    Transaction Approval   :done, approve, after inspect, 1d
    Payment Processing     :done, payment, after approve, 2d
    
    section Repair Phase
    Damage Assessment     :active, damage, after payment, 1d
    Work Order Creation   :work-order, after damage, 1d
    Parts Procurement     :parts, after work-order, 2d
    Repair Execution      :repair, after parts, 5d
    Quality Inspection    :quality, after repair, 1d
    Final Approval        :final, after quality, 1d
    
    section Completion
    Vehicle Ready         :milestone, ready, after final, 0d
```

### Data Flow Architecture

```mermaid
graph TB
    subgraph "Phase 4: Purchase"
        VPT[Vehicle Purchase Transactions]
        VPP[Vehicle Purchase Payments]
        VPT --> VPP
    end
    
    subgraph "Phase 5: Repair"
        VD[Vehicle Damages]
        RWO[Repair Work Orders]
        RWD[Repair Work Details]
        RPU[Repair Parts Usage]
        QI[Quality Inspections]
        
        VD --> RWO
        RWO --> RWD
        RWD --> RPU
        RWD --> QI
    end
    
    subgraph "Integration Points"
        CUST[Customers]
        USERS[Users]
        PROD[Products/Inventory]
    end
    
    VPT --> VD
    CUST --> VPT
    USERS --> VPT
    USERS --> RWO
    PROD --> RPU

    style VPT fill:#e1f5fe
    style VD fill:#fff3e0
    style QI fill:#e8f5e8
```

---

## API Endpoint Mapping

### Phase 4: Vehicle Purchase API Flow

```mermaid
graph LR
    subgraph "Transaction Management"
        A[POST /transactions] --> B[GET /transactions/{id}]
        B --> C[POST /transactions/{id}/inspect]
        C --> D[POST /transactions/{id}/approve]
    end
    
    subgraph "Payment Processing"
        E[POST /payments] --> F[POST /payments/{id}/process]
        F --> G[POST /payments/{id}/approve]
    end
    
    subgraph "Monitoring"
        H[GET /transactions/pending-inspection]
        I[GET /transactions/pending-approval]
        J[GET /payments/overdue]
        K[GET /dashboard]
    end
    
    D --> E
    G --> K
```

### Phase 5: Repair Management API Flow

```mermaid
graph TD
    subgraph "Damage Management"
        A[POST /damages] --> B[POST /damages/{id}/assess]
    end
    
    subgraph "Work Order Flow"
        C[POST /work-orders] --> D[POST /work-orders/{id}/assign]
        D --> E[POST /work-orders/{id}/approve]
    end
    
    subgraph "Task Execution"
        F[POST /work-details] --> G[POST /work-details/{id}/assign]
        G --> H[POST /work-details/{id}/progress]
        H --> I[POST /work-details/{id}/quality-check]
    end
    
    subgraph "Parts & Quality"
        J[POST /work-details/{id}/issue-parts]
        K[POST /inspections] --> L[POST /inspections/{id}/sign-off]
    end
    
    B --> C
    E --> F
    F --> J
    I --> K
```

---

## Business Rules Summary

### Phase 4 Rules
1. **VIN Uniqueness**: Each VIN can only exist once in the system
2. **Sequential Approval**: Inspection must be completed before approval
3. **Payment Validation**: Payments cannot exceed agreed vehicle value
4. **Status Progression**: Transactions follow: Pending → Inspection → Approved → Completed

### Phase 5 Rules
1. **Damage-First**: Work orders can only be created after damage assessment
2. **Parts Availability**: Tasks cannot start without required parts
3. **Quality Gates**: Each task must pass quality check before completion
4. **Progressive Completion**: Work details must be completed before work order completion
5. **Final Inspection**: All work orders require final quality inspection

### Integration Rules
1. **User Roles**: Different endpoints require different user roles (admin, sales, mechanic, inspector)
2. **Transaction Linkage**: Repair records must link to valid purchase transactions
3. **Audit Trail**: All critical operations are logged with user and timestamp
4. **Cost Tracking**: All costs are tracked and aggregated for profitability analysis

---

This comprehensive flow documentation provides clear guidance for implementation, testing, and operational understanding of the complete vehicle purchase and repair management system.