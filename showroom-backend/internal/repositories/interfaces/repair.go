package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/repair"
)

// VehicleDamageRepository defines the interface for vehicle damage data operations
type VehicleDamageRepository interface {
	Create(ctx context.Context, damage *repair.VehicleDamage) (*repair.VehicleDamage, error)
	GetByID(ctx context.Context, id int) (*repair.VehicleDamage, error)
	Update(ctx context.Context, id int, damage *repair.VehicleDamage) (*repair.VehicleDamage, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error)
	GetByTransactionID(ctx context.Context, transactionID int, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error)
	GetByCategory(ctx context.Context, category string, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error)
	GetBySeverity(ctx context.Context, severity string, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error)
	GetByStatus(ctx context.Context, status string, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error)
	GetHighPriority(ctx context.Context, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error)
	AssessDamage(ctx context.Context, id int, request *repair.DamageAssessmentRequest) error
	GetDamageSummary(ctx context.Context, transactionID int) (*repair.DamageSummary, error)
	Search(ctx context.Context, query string, params *repair.VehicleDamageFilterParams) (*common.PaginatedResponse, error)
}

// RepairWorkOrderRepository defines the interface for repair work order data operations
type RepairWorkOrderRepository interface {
	Create(ctx context.Context, workOrder *repair.RepairWorkOrder) (*repair.RepairWorkOrder, error)
	GetByID(ctx context.Context, id int) (*repair.RepairWorkOrder, error)
	GetByNumber(ctx context.Context, number string) (*repair.RepairWorkOrder, error)
	Update(ctx context.Context, id int, workOrder *repair.RepairWorkOrder) (*repair.RepairWorkOrder, error)
	UpdateStatus(ctx context.Context, id int, request *repair.WorkOrderStatusRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error)
	GetByTransactionID(ctx context.Context, transactionID int, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error)
	GetByStatus(ctx context.Context, status string, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error)
	GetByMechanic(ctx context.Context, mechanicID int, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error)
	GetPendingApproval(ctx context.Context, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error)
	AssignMechanic(ctx context.Context, id int, request *repair.WorkOrderAssignmentRequest) error
	ProcessApproval(ctx context.Context, id int, request *repair.WorkOrderApprovalRequest, approvedBy int) error
	GenerateNumber(ctx context.Context) (string, error)
	IsNumberExists(ctx context.Context, number string) (bool, error)
	GetWorkOrderSummary(ctx context.Context, transactionID int) (*repair.WorkOrderSummary, error)
	Search(ctx context.Context, query string, params *repair.RepairWorkOrderFilterParams) (*common.PaginatedResponse, error)
}

// RepairWorkDetailRepository defines the interface for repair work detail data operations
type RepairWorkDetailRepository interface {
	Create(ctx context.Context, workDetail *repair.RepairWorkDetail) (*repair.RepairWorkDetail, error)
	GetByID(ctx context.Context, id int) (*repair.RepairWorkDetail, error)
	Update(ctx context.Context, id int, workDetail *repair.RepairWorkDetail) (*repair.RepairWorkDetail, error)
	UpdateProgress(ctx context.Context, id int, request *repair.WorkDetailProgressRequest) error
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error)
	GetByWorkOrderID(ctx context.Context, workOrderID int, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error)
	GetByDamageID(ctx context.Context, damageID int, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error)
	GetByMechanic(ctx context.Context, mechanicID int, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error)
	GetByStatus(ctx context.Context, status string, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error)
	AssignMechanic(ctx context.Context, id int, request *repair.WorkDetailAssignmentRequest) error
	PerformQualityCheck(ctx context.Context, id int, request *repair.WorkDetailQualityCheckRequest, verifiedBy int) error
	GetWorkDetailSummary(ctx context.Context, workOrderID int) (*repair.WorkDetailSummary, error)
	Search(ctx context.Context, query string, params *repair.RepairWorkDetailFilterParams) (*common.PaginatedResponse, error)
}

// RepairPartsUsageRepository defines the interface for repair parts usage data operations
type RepairPartsUsageRepository interface {
	Create(ctx context.Context, usage *repair.RepairPartsUsage) (*repair.RepairPartsUsage, error)
	GetByID(ctx context.Context, id int) (*repair.RepairPartsUsage, error)
	Update(ctx context.Context, id int, usage *repair.RepairPartsUsage) (*repair.RepairPartsUsage, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error)
	GetByWorkDetailID(ctx context.Context, workDetailID int, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error)
	GetByProductID(ctx context.Context, productID int, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error)
	GetByUsageType(ctx context.Context, usageType string, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error)
	GetPendingApproval(ctx context.Context, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error)
	ProcessApproval(ctx context.Context, id int, request *repair.PartsUsageApprovalRequest, approvedBy int) error
	IssuePartsForRepair(ctx context.Context, workDetailID int, request *repair.PartsUsageIssueRequest, issuedBy int) error
	GetPartsUsageSummary(ctx context.Context, workOrderID int) (*repair.PartsUsageSummary, error)
	GetInventoryImpact(ctx context.Context, workOrderID int) ([]*repair.PartsInventoryImpact, error)
	Search(ctx context.Context, query string, params *repair.RepairPartsUsageFilterParams) (*common.PaginatedResponse, error)
}

// QualityInspectionRepository defines the interface for quality inspection data operations
type QualityInspectionRepository interface {
	Create(ctx context.Context, inspection *repair.QualityInspection) (*repair.QualityInspection, error)
	GetByID(ctx context.Context, id int) (*repair.QualityInspection, error)
	Update(ctx context.Context, id int, inspection *repair.QualityInspection) (*repair.QualityInspection, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error)
	GetByWorkOrderID(ctx context.Context, workOrderID int, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error)
	GetByInspector(ctx context.Context, inspectorID int, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error)
	GetByStatus(ctx context.Context, status string, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error)
	GetByType(ctx context.Context, inspectionType string, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error)
	GetFailedInspections(ctx context.Context, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error)
	GetReworkRequired(ctx context.Context, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error)
	SignOffInspection(ctx context.Context, id int, request *repair.InspectionSignOffRequest, signedOffBy int) error
	ScheduleRework(ctx context.Context, id int, request *repair.InspectionReworkRequest) error
	ScheduleInspection(ctx context.Context, request *repair.InspectionScheduleRequest) (*repair.QualityInspection, error)
	GetQualityMetrics(ctx context.Context, workOrderID int) (*repair.QualityMetrics, error)
	GetInspectionDashboard(ctx context.Context) (*repair.InspectionDashboard, error)
	Search(ctx context.Context, query string, params *repair.QualityInspectionFilterParams) (*common.PaginatedResponse, error)
}