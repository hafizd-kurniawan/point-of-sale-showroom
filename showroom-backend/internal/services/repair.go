package services

import (
	"context"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/repair"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// RepairService handles business logic for repair operations
type RepairService struct {
	damageRepo            interfaces.VehicleDamageRepository
	workOrderRepo         interfaces.RepairWorkOrderRepository
	workDetailRepo        interfaces.RepairWorkDetailRepository
	partsUsageRepo        interfaces.RepairPartsUsageRepository
	qualityInspectionRepo interfaces.QualityInspectionRepository
	transactionRepo       interfaces.VehiclePurchaseTransactionRepository
	userRepo              interfaces.UserRepository
	productRepo           interfaces.ProductRepository
}

// NewRepairService creates a new repair service
func NewRepairService(
	damageRepo interfaces.VehicleDamageRepository,
	workOrderRepo interfaces.RepairWorkOrderRepository,
	workDetailRepo interfaces.RepairWorkDetailRepository,
	partsUsageRepo interfaces.RepairPartsUsageRepository,
	qualityInspectionRepo interfaces.QualityInspectionRepository,
	transactionRepo interfaces.VehiclePurchaseTransactionRepository,
	userRepo interfaces.UserRepository,
	productRepo interfaces.ProductRepository,
) *RepairService {
	return &RepairService{
		damageRepo:            damageRepo,
		workOrderRepo:         workOrderRepo,
		workDetailRepo:        workDetailRepo,
		partsUsageRepo:        partsUsageRepo,
		qualityInspectionRepo: qualityInspectionRepo,
		transactionRepo:       transactionRepo,
		userRepo:              userRepo,
		productRepo:           productRepo,
	}
}

// Vehicle Damage Methods

// CreateDamage creates a new vehicle damage record
func (s *RepairService) CreateDamage(ctx context.Context, req *repair.CreateVehicleDamageRequest, assessedBy int) (*repair.VehicleDamage, error) {
	// Validate transaction exists
	_, err := s.transactionRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	damage := &repair.VehicleDamage{
		TransactionID:      req.TransactionID,
		DamageCategory:     req.DamageCategory,
		DamageSeverity:     req.DamageSeverity,
		DamageDescription:  req.DamageDescription,
		DamageLocation:     req.DamageLocation,
		DamageCause:        req.DamageCause,
		EstimatedCost:      req.EstimatedCost,
		AssessmentNotes:    req.AssessmentNotes,
		PhotosJSON:         req.PhotosJSON,
		DamageStatus:       "reported",
		AssessedBy:         assessedBy,
	}

	return s.damageRepo.Create(ctx, damage)
}

// GetDamage retrieves a damage record by ID
func (s *RepairService) GetDamage(ctx context.Context, id int) (*repair.VehicleDamage, error) {
	return s.damageRepo.GetByID(ctx, id)
}

// UpdateDamage updates a damage record
func (s *RepairService) UpdateDamage(ctx context.Context, id int, req *repair.UpdateVehicleDamageRequest) (*repair.VehicleDamage, error) {
	damage := &repair.VehicleDamage{
		DamageCategory:     req.DamageCategory,
		DamageSeverity:     req.DamageSeverity,
		DamageDescription:  req.DamageDescription,
		DamageLocation:     req.DamageLocation,
		DamageCause:        req.DamageCause,
		EstimatedCost:      req.EstimatedCost,
		AssessmentNotes:    req.AssessmentNotes,
		PhotosJSON:         req.PhotosJSON,
		DamageStatus:       req.DamageStatus,
	}

	return s.damageRepo.Update(ctx, id, damage)
}

// AssessDamage assesses damage with specific criteria
func (s *RepairService) AssessDamage(ctx context.Context, id int, req *repair.DamageAssessmentRequest) error {
	return s.damageRepo.AssessDamage(ctx, id, req)
}

// GetDamageSummary gets damage summary for a transaction
func (s *RepairService) GetDamageSummary(ctx context.Context, transactionID int) (*repair.DamageSummary, error) {
	return s.damageRepo.GetDamageSummary(ctx, transactionID)
}

// Work Order Methods

// CreateWorkOrder creates a new repair work order
func (s *RepairService) CreateWorkOrder(ctx context.Context, req *repair.CreateRepairWorkOrderRequest, createdBy int) (*repair.RepairWorkOrder, error) {
	// Validate transaction exists
	_, err := s.transactionRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Generate work order number
	workOrderNumber, err := s.workOrderRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate work order number: %w", err)
	}

	workOrder := &repair.RepairWorkOrder{
		WorkOrderNumber:        workOrderNumber,
		TransactionID:          req.TransactionID,
		PriorityLevel:          req.PriorityLevel,
		EstimatedCost:          req.EstimatedCost,
		EstimatedDurationHours: req.EstimatedDurationHours,
		WorkDescription:        req.WorkDescription,
		SpecialInstructions:    req.SpecialInstructions,
		WorkOrderStatus:        "pending",
		CreatedBy:              createdBy,
	}

	return s.workOrderRepo.Create(ctx, workOrder)
}

// GetWorkOrder retrieves a work order by ID
func (s *RepairService) GetWorkOrder(ctx context.Context, id int) (*repair.RepairWorkOrder, error) {
	return s.workOrderRepo.GetByID(ctx, id)
}

// GetWorkOrderByNumber retrieves a work order by number
func (s *RepairService) GetWorkOrderByNumber(ctx context.Context, number string) (*repair.RepairWorkOrder, error) {
	return s.workOrderRepo.GetByNumber(ctx, number)
}

// UpdateWorkOrder updates a work order
func (s *RepairService) UpdateWorkOrder(ctx context.Context, id int, req *repair.UpdateRepairWorkOrderRequest) (*repair.RepairWorkOrder, error) {
	workOrder := &repair.RepairWorkOrder{
		PriorityLevel:          req.PriorityLevel,
		EstimatedCost:          req.EstimatedCost,
		ActualCost:             req.ActualCost,
		EstimatedDurationHours: req.EstimatedDurationHours,
		ActualDurationHours:    req.ActualDurationHours,
		WorkDescription:        req.WorkDescription,
		SpecialInstructions:    req.SpecialInstructions,
		WorkOrderStatus:        req.WorkOrderStatus,
		StartDate:              req.StartDate,
		TargetCompletionDate:   req.TargetCompletionDate,
		ActualCompletionDate:   req.ActualCompletionDate,
	}

	return s.workOrderRepo.Update(ctx, id, workOrder)
}

// AssignMechanic assigns a mechanic to a work order
func (s *RepairService) AssignMechanic(ctx context.Context, id int, req *repair.WorkOrderAssignmentRequest) error {
	// Validate mechanic exists
	_, err := s.userRepo.GetByID(ctx, req.MechanicID)
	if err != nil {
		return fmt.Errorf("mechanic not found: %w", err)
	}

	return s.workOrderRepo.AssignMechanic(ctx, id, req)
}

// ProcessWorkOrderApproval processes work order approval
func (s *RepairService) ProcessWorkOrderApproval(ctx context.Context, id int, req *repair.WorkOrderApprovalRequest, approvedBy int) error {
	return s.workOrderRepo.ProcessApproval(ctx, id, req, approvedBy)
}

// GetWorkOrderSummary gets work order summary for a transaction
func (s *RepairService) GetWorkOrderSummary(ctx context.Context, transactionID int) (*repair.WorkOrderSummary, error) {
	return s.workOrderRepo.GetWorkOrderSummary(ctx, transactionID)
}

// Work Detail Methods

// CreateWorkDetail creates a new work detail
func (s *RepairService) CreateWorkDetail(ctx context.Context, req *repair.CreateRepairWorkDetailRequest) (*repair.RepairWorkDetail, error) {
	// Validate work order exists
	_, err := s.workOrderRepo.GetByID(ctx, req.WorkOrderID)
	if err != nil {
		return nil, fmt.Errorf("work order not found: %w", err)
	}

	// Validate damage exists if provided
	if req.DamageID > 0 {
		_, err := s.damageRepo.GetByID(ctx, req.DamageID)
		if err != nil {
			return nil, fmt.Errorf("damage record not found: %w", err)
		}
	}

	workDetail := &repair.RepairWorkDetail{
		WorkOrderID:      req.WorkOrderID,
		DamageID:         req.DamageID,
		TaskDescription:  req.TaskDescription,
		EstimatedHours:   req.EstimatedHours,
		TaskPriority:     req.TaskPriority,
		TaskCategory:     req.TaskCategory,
		WorkInstructions: req.WorkInstructions,
		TaskStatus:       "pending",
	}

	return s.workDetailRepo.Create(ctx, workDetail)
}

// GetWorkDetail retrieves a work detail by ID
func (s *RepairService) GetWorkDetail(ctx context.Context, id int) (*repair.RepairWorkDetail, error) {
	return s.workDetailRepo.GetByID(ctx, id)
}

// UpdateWorkDetail updates a work detail
func (s *RepairService) UpdateWorkDetail(ctx context.Context, id int, req *repair.UpdateRepairWorkDetailRequest) (*repair.RepairWorkDetail, error) {
	workDetail := &repair.RepairWorkDetail{
		TaskDescription:    req.TaskDescription,
		EstimatedHours:     req.EstimatedHours,
		ActualHours:        req.ActualHours,
		TaskPriority:       req.TaskPriority,
		TaskCategory:       req.TaskCategory,
		WorkInstructions:   req.WorkInstructions,
		TaskStatus:         req.TaskStatus,
		ProgressPercentage: req.ProgressPercentage,
		QualityCheckStatus: req.QualityCheckStatus,
		CompletionNotes:    req.CompletionNotes,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
	}

	return s.workDetailRepo.Update(ctx, id, workDetail)
}

// UpdateWorkDetailProgress updates work detail progress
func (s *RepairService) UpdateWorkDetailProgress(ctx context.Context, id int, req *repair.WorkDetailProgressRequest) error {
	return s.workDetailRepo.UpdateProgress(ctx, id, req)
}

// AssignWorkDetailMechanic assigns a mechanic to a work detail
func (s *RepairService) AssignWorkDetailMechanic(ctx context.Context, id int, req *repair.WorkDetailAssignmentRequest) error {
	// Validate mechanic exists
	_, err := s.userRepo.GetByID(ctx, req.MechanicID)
	if err != nil {
		return fmt.Errorf("mechanic not found: %w", err)
	}

	return s.workDetailRepo.AssignMechanic(ctx, id, req)
}

// PerformQualityCheck performs quality check on work detail
func (s *RepairService) PerformQualityCheck(ctx context.Context, id int, req *repair.WorkDetailQualityCheckRequest, verifiedBy int) error {
	return s.workDetailRepo.PerformQualityCheck(ctx, id, req, verifiedBy)
}

// GetWorkDetailSummary gets work detail summary for a work order
func (s *RepairService) GetWorkDetailSummary(ctx context.Context, workOrderID int) (*repair.WorkDetailSummary, error) {
	return s.workDetailRepo.GetWorkDetailSummary(ctx, workOrderID)
}

// Parts Usage Methods

// CreatePartsUsage creates a new parts usage record
func (s *RepairService) CreatePartsUsage(ctx context.Context, req *repair.CreateRepairPartsUsageRequest, issuedBy int) (*repair.RepairPartsUsage, error) {
	// Validate work detail exists
	_, err := s.workDetailRepo.GetByID(ctx, req.WorkDetailID)
	if err != nil {
		return nil, fmt.Errorf("work detail not found: %w", err)
	}

	// Validate product exists
	_, err = s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	partsUsage := &repair.RepairPartsUsage{
		WorkDetailID:  req.WorkDetailID,
		ProductID:     req.ProductID,
		QuantityUsed:  req.QuantityUsed,
		UnitCost:      req.UnitCost,
		TotalCost:     req.TotalCost,
		UsageType:     req.UsageType,
		UsageNotes:    req.UsageNotes,
		UsageStatus:   "pending",
		IssuedBy:      issuedBy,
	}

	return s.partsUsageRepo.Create(ctx, partsUsage)
}

// GetPartsUsage retrieves a parts usage by ID
func (s *RepairService) GetPartsUsage(ctx context.Context, id int) (*repair.RepairPartsUsage, error) {
	return s.partsUsageRepo.GetByID(ctx, id)
}

// ProcessPartsUsageApproval processes parts usage approval
func (s *RepairService) ProcessPartsUsageApproval(ctx context.Context, id int, req *repair.PartsUsageApprovalRequest, approvedBy int) error {
	return s.partsUsageRepo.ProcessApproval(ctx, id, req, approvedBy)
}

// IssuePartsForRepair issues parts for repair
func (s *RepairService) IssuePartsForRepair(ctx context.Context, workDetailID int, req *repair.PartsUsageIssueRequest, issuedBy int) error {
	return s.partsUsageRepo.IssuePartsForRepair(ctx, workDetailID, req, issuedBy)
}

// GetPartsUsageSummary gets parts usage summary for a work order
func (s *RepairService) GetPartsUsageSummary(ctx context.Context, workOrderID int) (*repair.PartsUsageSummary, error) {
	return s.partsUsageRepo.GetPartsUsageSummary(ctx, workOrderID)
}

// GetInventoryImpact gets inventory impact for a work order
func (s *RepairService) GetInventoryImpact(ctx context.Context, workOrderID int) ([]*repair.PartsInventoryImpact, error) {
	return s.partsUsageRepo.GetInventoryImpact(ctx, workOrderID)
}

// Quality Inspection Methods

// CreateInspection creates a new quality inspection
func (s *RepairService) CreateInspection(ctx context.Context, req *repair.CreateQualityInspectionRequest) (*repair.QualityInspection, error) {
	// Validate work order exists
	_, err := s.workOrderRepo.GetByID(ctx, req.WorkOrderID)
	if err != nil {
		return nil, fmt.Errorf("work order not found: %w", err)
	}

	// Validate inspector exists
	if req.InspectorID > 0 {
		_, err := s.userRepo.GetByID(ctx, req.InspectorID)
		if err != nil {
			return nil, fmt.Errorf("inspector not found: %w", err)
		}
	}

	inspection := &repair.QualityInspection{
		WorkOrderID:             req.WorkOrderID,
		InspectionType:         req.InspectionType,
		InspectionChecklistJSON: req.InspectionChecklistJSON,
		QualityStandardsJSON:   req.QualityStandardsJSON,
		InspectionStatus:       "scheduled",
		ScheduledDate:          req.ScheduledDate,
		InspectorID:            req.InspectorID,
	}

	return s.qualityInspectionRepo.Create(ctx, inspection)
}

// GetInspection retrieves an inspection by ID
func (s *RepairService) GetInspection(ctx context.Context, id int) (*repair.QualityInspection, error) {
	return s.qualityInspectionRepo.GetByID(ctx, id)
}

// SignOffInspection signs off an inspection
func (s *RepairService) SignOffInspection(ctx context.Context, id int, req *repair.InspectionSignOffRequest, signedOffBy int) error {
	return s.qualityInspectionRepo.SignOffInspection(ctx, id, req, signedOffBy)
}

// ScheduleRework schedules rework for failed inspection
func (s *RepairService) ScheduleRework(ctx context.Context, id int, req *repair.InspectionReworkRequest) error {
	return s.qualityInspectionRepo.ScheduleRework(ctx, id, req)
}

// ScheduleInspection schedules a new inspection
func (s *RepairService) ScheduleInspection(ctx context.Context, req *repair.InspectionScheduleRequest) (*repair.QualityInspection, error) {
	return s.qualityInspectionRepo.ScheduleInspection(ctx, req)
}

// GetQualityMetrics gets quality metrics for a work order
func (s *RepairService) GetQualityMetrics(ctx context.Context, workOrderID int) (*repair.QualityMetrics, error) {
	return s.qualityInspectionRepo.GetQualityMetrics(ctx, workOrderID)
}

// GetInspectionDashboard gets inspection dashboard data
func (s *RepairService) GetInspectionDashboard(ctx context.Context) (*repair.InspectionDashboard, error) {
	return s.qualityInspectionRepo.GetInspectionDashboard(ctx)
}

// List Methods (stubs for now)
func (s *RepairService) ListDamages(ctx context.Context, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error) {
	return s.damageRepo.List(ctx, params)
}

func (s *RepairService) ListWorkOrders(ctx context.Context, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error) {
	return s.workOrderRepo.List(ctx, params)
}

func (s *RepairService) ListWorkDetails(ctx context.Context, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error) {
	return s.workDetailRepo.List(ctx, params)
}

func (s *RepairService) ListPartsUsage(ctx context.Context, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error) {
	return s.partsUsageRepo.List(ctx, params)
}

func (s *RepairService) ListInspections(ctx context.Context, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return s.qualityInspectionRepo.List(ctx, params)
}

func (s *RepairService) GetHighPriorityDamages(ctx context.Context, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error) {
	return s.damageRepo.GetHighPriority(ctx, params)
}

func (s *RepairService) GetPendingWorkOrderApprovals(ctx context.Context, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error) {
	return s.workOrderRepo.GetPendingApproval(ctx, params)
}

func (s *RepairService) GetWorkOrdersByMechanic(ctx context.Context, mechanicID int, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error) {
	return s.workOrderRepo.GetByMechanic(ctx, mechanicID, params)
}

func (s *RepairService) GetWorkDetailsByWorkOrder(ctx context.Context, workOrderID int, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error) {
	return s.workDetailRepo.GetByWorkOrderID(ctx, workOrderID, params)
}

func (s *RepairService) GetPartsUsageByWorkDetail(ctx context.Context, workDetailID int, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error) {
	return s.partsUsageRepo.GetByWorkDetailID(ctx, workDetailID, params)
}

func (s *RepairService) GetInspectionsByWorkOrder(ctx context.Context, workOrderID int, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return s.qualityInspectionRepo.GetByWorkOrderID(ctx, workOrderID, params)
}

func (s *RepairService) GetFailedInspections(ctx context.Context, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return s.qualityInspectionRepo.GetFailedInspections(ctx, params)
}

func (s *RepairService) GetReworkRequired(ctx context.Context, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error) {
	return s.qualityInspectionRepo.GetReworkRequired(ctx, params)
}