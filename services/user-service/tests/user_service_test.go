package tests

import (
	"context"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/user-service/internal/model"
	"enterprise-microservice-system/services/user-service/internal/repository"
	"enterprise-microservice-system/services/user-service/internal/service"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate schema
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, nil)

	ctx := context.Background()

	tests := []struct {
		name    string
		req     *model.CreateUserRequest
		wantErr bool
	}{
		{
			name: "valid user",
			req: &model.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req: &model.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Another User",
				Age:   30,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
	user, err := svc.CreateUser(ctx, tt.req, "tester")
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && user == nil {
				t.Error("Expected user to be created, got nil")
			}

			if !tt.wantErr {
				if user.Email != tt.req.Email {
					t.Errorf("Expected email %s, got %s", tt.req.Email, user.Email)
				}
				if user.Name != tt.req.Name {
					t.Errorf("Expected name %s, got %s", tt.req.Name, user.Name)
				}
				if user.Age != tt.req.Age {
					t.Errorf("Expected age %d, got %d", tt.req.Age, user.Age)
				}
				if user.Status != model.UserStatusActive {
					t.Errorf("Expected status %s, got %s", model.UserStatusActive, user.Status)
				}
				if user.CreatedBy == "" || user.UpdatedBy == "" {
					t.Error("Expected created_by and updated_by to be set")
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, nil)

	ctx := context.Background()

	// Create a test user
	createReq := &model.CreateUserRequest{
		Email: "get@example.com",
		Name:  "Get User",
		Age:   28,
	}
	createdUser, err := svc.CreateUser(ctx, createReq, "tester")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      createdUser.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      9999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.GetUser(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && user == nil {
				t.Error("Expected user to be retrieved, got nil")
			}

			if !tt.wantErr && user.ID != tt.id {
				t.Errorf("Expected user ID %d, got %d", tt.id, user.ID)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, nil)

	ctx := context.Background()

	// Create a test user
	createReq := &model.CreateUserRequest{
		Email: "update@example.com",
		Name:  "Update User",
		Age:   30,
	}
	createdUser, err := svc.CreateUser(ctx, createReq, "tester")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	newName := "Updated Name"
	newAge := 35
	updateReq := &model.UpdateUserRequest{
		Name: &newName,
		Age:  &newAge,
	}

	updatedUser, err := svc.UpdateUser(ctx, createdUser.ID, updateReq, "tester")
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if updatedUser.Name != newName {
		t.Errorf("Expected name %s, got %s", newName, updatedUser.Name)
	}

	if updatedUser.Age != newAge {
		t.Errorf("Expected age %d, got %d", newAge, updatedUser.Age)
	}
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, nil)

	ctx := context.Background()

	// Create a test user
	createReq := &model.CreateUserRequest{
		Email: "delete@example.com",
		Name:  "Delete User",
		Age:   32,
	}
	createdUser, err := svc.CreateUser(ctx, createReq, "tester")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Delete the user
	err = svc.DeleteUser(ctx, createdUser.ID, "tester")
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Try to get the deleted user (should fail)
	_, err = svc.GetUser(ctx, createdUser.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user, got nil")
	}

	appErr, ok := err.(*errors.AppError)
	if !ok || appErr.Code != errors.ErrCodeNotFound {
		t.Errorf("Expected NotFound error, got %v", err)
	}
}

func TestListUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, nil)

	ctx := context.Background()

	// Create multiple test users
	for i := 1; i <= 5; i++ {
		createReq := &model.CreateUserRequest{
			Email: testEmail(i),
			Name:  testName(i),
			Age:   20 + i,
		}
		_, err := svc.CreateUser(ctx, createReq, "tester")
		if err != nil {
			t.Fatalf("Failed to create test user %d: %v", i, err)
		}
	}

	query := &model.ListUsersQuery{
		Page:     1,
		PageSize: 10,
	}

	users, total, err := svc.ListUsers(ctx, query)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(users))
	}

	if total != 5 {
		t.Errorf("Expected total count 5, got %d", total)
	}
}

func testEmail(i int) string {
	return "user" + string(rune('0'+i)) + "@example.com"
}

func testName(i int) string {
	return "User " + string(rune('0'+i))
}
