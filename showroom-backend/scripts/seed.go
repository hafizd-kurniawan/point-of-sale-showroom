package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/database"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/implementations"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	userRepo := implementations.NewUserRepository(database.GetDB())

	// Seed data
	if err := seedUsers(userRepo); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	log.Println("Database seeding completed successfully")
}

func seedUsers(userRepo interfaces.UserRepository) error {
	ctx := context.Background()

	// Check if admin user already exists
	_, err := userRepo.GetByUsername(ctx, "admin")
	if err == nil {
		log.Println("Admin user already exists, skipping seeding")
		return nil
	}

	log.Println("Creating seed users...")

	// Create admin user first (as it will create other users)
	adminPassword, _ := utils.HashPassword("admin123")
	hireDate := time.Now().Truncate(24 * time.Hour)

	adminUser := &user.User{
		Username:     "admin",
		Email:        "admin@showroom.com",
		PasswordHash: adminPassword,
		FullName:     "System Administrator",
		Phone:        "081234567890",
		Address:      stringPtr("Jl. Admin Street No. 1, Jakarta"),
		Role:         common.RoleSuperAdmin,
		Salary:       float64Ptr(15000000),
		HireDate:     &hireDate,
		CreatedBy:    1, // Self-reference, will be updated after creation
		IsActive:     true,
		ProfileImage: nil,
		Notes:        stringPtr("System administrator account"),
	}

	createdAdmin, err := userRepo.Create(ctx, adminUser)
	if err != nil {
		return err
	}

	log.Printf("Created admin user with ID: %d", createdAdmin.UserID)

	// Update admin user's created_by to reference itself
	_, err = userRepo.Update(ctx, createdAdmin.UserID, &user.User{
		UserID:       createdAdmin.UserID,
		Username:     createdAdmin.Username,
		Email:        createdAdmin.Email,
		PasswordHash: createdAdmin.PasswordHash,
		FullName:     createdAdmin.FullName,
		Phone:        createdAdmin.Phone,
		Address:      createdAdmin.Address,
		Role:         createdAdmin.Role,
		Salary:       createdAdmin.Salary,
		HireDate:     createdAdmin.HireDate,
		CreatedBy:    createdAdmin.UserID, // Self-reference
		IsActive:     createdAdmin.IsActive,
		ProfileImage: createdAdmin.ProfileImage,
		Notes:        createdAdmin.Notes,
	})
	if err != nil {
		return err
	}

	// Create other test users
	testUsers := []struct {
		username string
		email    string
		password string
		fullName string
		phone    string
		role     common.UserRole
		salary   float64
	}{
		{"kasir1", "kasir1@showroom.com", "kasir123", "Siti Kasir", "081234567891", common.RoleKasir, 5000000},
		{"mekanik1", "mekanik1@showroom.com", "mekanik123", "Budi Mekanik", "081234567892", common.RoleMekanik, 6000000},
	}

	for _, testUser := range testUsers {
		hashedPassword, err := utils.HashPassword(testUser.password)
		if err != nil {
			return err
		}

		newUser := &user.User{
			Username:     testUser.username,
			Email:        testUser.email,
			PasswordHash: hashedPassword,
			FullName:     testUser.fullName,
			Phone:        testUser.phone,
			Address:      stringPtr("Jl. Test Street No. 123, Jakarta"),
			Role:         testUser.role,
			Salary:       &testUser.salary,
			HireDate:     &hireDate,
			CreatedBy:    createdAdmin.UserID,
			IsActive:     true,
			ProfileImage: nil,
			Notes:        stringPtr("Test user account"),
		}

		createdUser, err := userRepo.Create(ctx, newUser)
		if err != nil {
			return err
		}

		log.Printf("Created %s user: %s (ID: %d)", testUser.role, testUser.username, createdUser.UserID)
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}