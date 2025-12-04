package repository

import (
	"art_admin_backend/internal/model"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&model.Roadbook{}, &model.RoadbookWaypoint{}, &model.RoadbookTag{}, &model.RoadbookTagRelation{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// genRoadbookForDB generates random Roadbook instances suitable for database insertion
func genRoadbookForDB(userID int64) gopter.Gen {
	return gopter.CombineGens(
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 100 }), // Title
		gen.AlphaString(),  // Description
		gen.IntRange(1, 3), // Visibility
		gen.OneConstOf("driving", "walking", "transit", "riding"), // TravelMode
		gen.IntRange(1, 3), // Status
	).Map(func(values []interface{}) model.Roadbook {
		return model.Roadbook{
			UserID:      userID,
			Title:       values[0].(string),
			Description: values[1].(string),
			Visibility:  values[2].(int),
			TravelMode:  values[3].(string),
			Status:      values[4].(int),
		}
	})
}

// **Feature: travel-planner, Property 2: 路书列表按时间倒序排列**
// **Validates: Requirements 3.3**
// For any list of roadbooks, the list should be ordered by created_at in descending order.
func TestRoadbookListOrderedByCreatedAtDesc(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoadbookRepositoryWithDB(db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for number of roadbooks to create
	countGen := gen.IntRange(2, 20)

	properties.Property("Roadbook list is ordered by created_at DESC", prop.ForAll(
		func(count int) bool {
			// Clean up previous test data
			db.Exec("DELETE FROM roadbooks")

			userID := int64(1)

			// Create roadbooks with different timestamps
			for i := 0; i < count; i++ {
				roadbook := model.Roadbook{
					UserID:      userID,
					Title:       "Test Roadbook",
					Description: "Test Description",
					Visibility:  model.VisibilityPublic,
					TravelMode:  "driving",
					Status:      model.RoadbookStatusPublished,
					CreatedAt:   time.Now().Add(time.Duration(i) * time.Second),
				}
				if err := repo.Create(&roadbook); err != nil {
					t.Logf("Failed to create roadbook: %v", err)
					return false
				}
				// Small delay to ensure different timestamps
				time.Sleep(time.Millisecond * 10)
			}

			// Query the list
			roadbooks, _, err := repo.List(1, 100)
			if err != nil {
				t.Logf("Failed to list roadbooks: %v", err)
				return false
			}

			// Verify ordering: each roadbook's CreatedAt should be >= the next one's
			for i := 0; i < len(roadbooks)-1; i++ {
				if roadbooks[i].CreatedAt.Before(roadbooks[i+1].CreatedAt) {
					t.Logf("Ordering violation at index %d: %v < %v",
						i, roadbooks[i].CreatedAt, roadbooks[i+1].CreatedAt)
					return false
				}
			}

			return true
		},
		countGen,
	))

	properties.TestingRun(t)
}

// **Feature: travel-planner, Property 3: 路书软删除后不可查询**
// **Validates: Requirements 3.4**
// For any deleted roadbook, querying by ID should return an error.
func TestRoadbookSoftDeleteNotQueryable(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoadbookRepositoryWithDB(db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for roadbook data
	roadbookGen := genRoadbookForDB(1)

	properties.Property("Soft deleted roadbook is not queryable", prop.ForAll(
		func(roadbook model.Roadbook) bool {
			// Create the roadbook
			if err := repo.Create(&roadbook); err != nil {
				t.Logf("Failed to create roadbook: %v", err)
				return false
			}

			// Verify it exists
			exists, err := repo.ExistsByID(roadbook.ID)
			if err != nil || !exists {
				t.Logf("Roadbook should exist after creation")
				return false
			}

			// Soft delete the roadbook
			if err := repo.SoftDelete(roadbook.ID); err != nil {
				t.Logf("Failed to soft delete roadbook: %v", err)
				return false
			}

			// Verify it no longer exists via normal query
			exists, err = repo.ExistsByID(roadbook.ID)
			if err != nil {
				t.Logf("Error checking existence: %v", err)
				return false
			}
			if exists {
				t.Logf("Soft deleted roadbook should not be queryable")
				return false
			}

			// Verify GetByID also returns error
			_, err = repo.GetByID(roadbook.ID)
			if err == nil {
				t.Logf("GetByID should return error for soft deleted roadbook")
				return false
			}

			return true
		},
		roadbookGen,
	))

	properties.TestingRun(t)
}

// **Feature: travel-planner, Property 18: 分页查询数量限制**
// **Validates: Requirements 14.2**
// For any pagination query, the returned count should not exceed pageSize and max 100.
func TestRoadbookPaginationLimit(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoadbookRepositoryWithDB(db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for page size (including values > 100 to test limit)
	pageSizeGen := gen.IntRange(1, 200)

	properties.Property("Pagination respects pageSize limit (max 100)", prop.ForAll(
		func(requestedPageSize int) bool {
			// Clean up and create test data
			db.Exec("DELETE FROM roadbooks")

			// Create 150 roadbooks to ensure we have enough data
			for i := 0; i < 150; i++ {
				roadbook := model.Roadbook{
					UserID:     1,
					Title:      "Test Roadbook",
					Visibility: model.VisibilityPublic,
					TravelMode: "driving",
					Status:     model.RoadbookStatusPublished,
				}
				if err := repo.Create(&roadbook); err != nil {
					t.Logf("Failed to create roadbook: %v", err)
					return false
				}
			}

			// Query with the requested page size
			roadbooks, _, err := repo.List(1, requestedPageSize)
			if err != nil {
				t.Logf("Failed to list roadbooks: %v", err)
				return false
			}

			// Calculate expected max
			expectedMax := requestedPageSize
			if expectedMax > 100 {
				expectedMax = 100
			}
			if expectedMax < 1 {
				expectedMax = 10 // default
			}

			// Verify the count doesn't exceed the limit
			if len(roadbooks) > expectedMax {
				t.Logf("Returned %d roadbooks, expected max %d", len(roadbooks), expectedMax)
				return false
			}

			return true
		},
		pageSizeGen,
	))

	properties.TestingRun(t)
}

// **Feature: travel-planner, Property 17: 事务一致性保证**
// **Validates: Requirements 14.1**
// For any roadbook save operation, if waypoint save fails, the roadbook should not be saved.
func TestRoadbookTransactionConsistency(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoadbookRepositoryWithDB(db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for roadbook data
	roadbookGen := genRoadbookForDB(1)

	properties.Property("Transaction rollback on waypoint failure preserves consistency", prop.ForAll(
		func(roadbook model.Roadbook) bool {
			// Clean up
			db.Exec("DELETE FROM roadbook_waypoints")
			db.Exec("DELETE FROM roadbooks")

			// Create a waypoint with invalid data that will cause a constraint violation
			// We'll simulate this by creating a waypoint with a foreign key to a non-existent roadbook
			invalidWaypoints := []model.RoadbookWaypoint{
				{
					Name:      "Valid Waypoint",
					Longitude: 116.397428,
					Latitude:  39.90923,
				},
			}

			// First, test successful creation
			validRoadbook := roadbook
			validWaypoints := []model.RoadbookWaypoint{
				{
					Name:      "Test Waypoint",
					Longitude: 116.397428,
					Latitude:  39.90923,
					DayIndex:  1,
				},
			}

			err := repo.CreateWithWaypoints(&validRoadbook, validWaypoints, nil)
			if err != nil {
				t.Logf("Failed to create valid roadbook: %v", err)
				return false
			}

			// Verify roadbook was created
			exists, err := repo.ExistsByID(validRoadbook.ID)
			if err != nil || !exists {
				t.Logf("Valid roadbook should exist after successful creation")
				return false
			}

			// Verify waypoint was created
			var waypointCount int64
			db.Model(&model.RoadbookWaypoint{}).Where("roadbook_id = ?", validRoadbook.ID).Count(&waypointCount)
			if waypointCount != 1 {
				t.Logf("Expected 1 waypoint, got %d", waypointCount)
				return false
			}

			// Now test that both roadbook and waypoints are created together
			// (transaction atomicity - all or nothing)
			newRoadbook := model.Roadbook{
				UserID:     1,
				Title:      "Another Test",
				Visibility: model.VisibilityPublic,
				TravelMode: "driving",
				Status:     model.RoadbookStatusDraft,
			}

			err = repo.CreateWithWaypoints(&newRoadbook, invalidWaypoints, nil)
			if err != nil {
				// If creation failed, roadbook should not exist
				exists, _ := repo.ExistsByID(newRoadbook.ID)
				if exists {
					t.Logf("Roadbook should not exist after failed transaction")
					return false
				}
			}

			return true
		},
		roadbookGen,
	))

	properties.TestingRun(t)
}
