package models_test

import (
	"testing"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/stretchr/testify/assert"
)

func TestUserRole_IsValid(t *testing.T) {
	tests := []struct {
		role  common.UserRole
		valid bool
	}{
		{common.RoleSuperAdmin, true},
		{common.RoleKasir, true},
		{common.RoleMekanik, true},
		{common.UserRole("invalid"), false},
		{common.UserRole(""), false},
	}

	for _, test := range tests {
		assert.Equal(t, test.valid, test.role.IsValid(), "Role %s validity should be %v", test.role, test.valid)
	}
}

func TestUserRole_String(t *testing.T) {
	assert.Equal(t, "super_admin", common.RoleSuperAdmin.String())
	assert.Equal(t, "kasir", common.RoleKasir.String())
	assert.Equal(t, "mekanik", common.RoleMekanik.String())
}

func TestPaginationParams_Validate(t *testing.T) {
	tests := []struct {
		input    common.PaginationParams
		expected common.PaginationParams
	}{
		{common.PaginationParams{Page: 0, Limit: 0}, common.PaginationParams{Page: 1, Limit: 10}},
		{common.PaginationParams{Page: -1, Limit: -1}, common.PaginationParams{Page: 1, Limit: 10}},
		{common.PaginationParams{Page: 2, Limit: 20}, common.PaginationParams{Page: 2, Limit: 20}},
		{common.PaginationParams{Page: 1, Limit: 200}, common.PaginationParams{Page: 1, Limit: 100}},
	}

	for _, test := range tests {
		test.input.Validate()
		assert.Equal(t, test.expected, test.input)
	}
}

func TestPaginationParams_GetOffset(t *testing.T) {
	tests := []struct {
		params common.PaginationParams
		offset int
	}{
		{common.PaginationParams{Page: 1, Limit: 10}, 0},
		{common.PaginationParams{Page: 2, Limit: 10}, 10},
		{common.PaginationParams{Page: 3, Limit: 20}, 40},
	}

	for _, test := range tests {
		assert.Equal(t, test.offset, test.params.GetOffset())
	}
}

func TestPaginationParams_GetTotalPages(t *testing.T) {
	tests := []struct {
		params common.PaginationParams
		total  int
		pages  int
	}{
		{common.PaginationParams{Page: 1, Limit: 10}, 25, 3},
		{common.PaginationParams{Page: 1, Limit: 10}, 30, 3},
		{common.PaginationParams{Page: 1, Limit: 10}, 35, 4},
		{common.PaginationParams{Page: 1, Limit: 10}, 0, 0},
	}

	for _, test := range tests {
		assert.Equal(t, test.pages, test.params.GetTotalPages(test.total))
	}
}

func TestPaginationParams_GetHasMore(t *testing.T) {
	tests := []struct {
		params  common.PaginationParams
		total   int
		hasMore bool
	}{
		{common.PaginationParams{Page: 1, Limit: 10}, 25, true},
		{common.PaginationParams{Page: 3, Limit: 10}, 25, false},
		{common.PaginationParams{Page: 2, Limit: 10}, 15, false},
		{common.PaginationParams{Page: 1, Limit: 10}, 0, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.hasMore, test.params.GetHasMore(test.total))
	}
}