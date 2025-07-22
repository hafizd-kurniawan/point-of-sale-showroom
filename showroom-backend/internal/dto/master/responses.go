package master

import (
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
)

// CustomerListResponse represents the response for customer list
type CustomerListResponse struct {
	Customers []master.CustomerListItem `json:"customers"`
	common.PaginationMeta
}

// SupplierListResponse represents the response for supplier list
type SupplierListResponse struct {
	Suppliers []master.SupplierListItem `json:"suppliers"`
	common.PaginationMeta
}

// VehicleModelListResponse represents the response for vehicle model list
type VehicleModelListResponse struct {
	Models []master.VehicleModelListItem `json:"models"`
	common.PaginationMeta
}