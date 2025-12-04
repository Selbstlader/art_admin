package repository

import (
	"art_admin_backend/internal/model"
	"math"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupWaypointTestDB creates an in-memory SQLite database for waypoint testing
func setupWaypointTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&model.Roadbook{}, &model.RoadbookWaypoint{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// createTestRoadbook creates a test roadbook and returns its ID
func createTestRoadbook(db *gorm.DB) (int64, error) {
	roadbook := &model.Roadbook{
		UserID:     1,
		Title:      "Test Roadbook",
		Visibility: model.VisibilityPublic,
		TravelMode: "driving",
		Status:     model.RoadbookStatusDraft,
	}
	if err := db.Create(roadbook).Error; err != nil {
		return 0, err
	}
	return roadbook.ID, nil
}

// genWaypointForDB generates random RoadbookWaypoint instances suitable for database insertion
func genWaypointForDB(roadbookID int64) gopter.Gen {
	return gopter.CombineGens(
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 100 }), // Name
		gen.AlphaString(),               // Address
		gen.Float64Range(-180.0, 180.0), // Longitude
		gen.Float64Range(-90.0, 90.0),   // Latitude
		gen.IntRange(1, 30),             // DayIndex
		gen.IntRange(0, 100),            // SortOrder
		gen.IntRange(15, 480),           // StayDuration (15 min to 8 hours)
		gen.Float64Range(0, 10000),      // Budget
		gen.AlphaString(),               // Notes
	).Map(func(values []interface{}) model.RoadbookWaypoint {
		return model.RoadbookWaypoint{
			RoadbookID:   roadbookID,
			Name:         values[0].(string),
			Address:      values[1].(string),
			Longitude:    values[2].(float64),
			Latitude:     values[3].(float64),
			DayIndex:     values[4].(int),
			SortOrder:    values[5].(int),
			StayDuration: values[6].(int),
			Budget:       values[7].(float64),
			Notes:        values[8].(string),
			WaypointType: model.WaypointTypeWaypoint,
		}
	})
}

// **Feature: travel-planner, Property 4: 途经点数据完整性**
// **Validates: Requirements 3.5**
// For any waypoint added to a roadbook, querying should return complete information
// including name, longitude, latitude, stay duration, and notes.
func TestWaypointDataIntegrity(t *testing.T) {
	db := setupWaypointTestDB(t)
	repo := NewWaypointRepositoryWithDB(db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Create a test roadbook first
	roadbookID, err := createTestRoadbook(db)
	if err != nil {
		t.Fatalf("Failed to create test roadbook: %v", err)
	}

	waypointGen := genWaypointForDB(roadbookID)

	properties.Property("Waypoint data integrity is preserved after save and query", prop.ForAll(
		func(waypoint model.RoadbookWaypoint) bool {
			// Save the waypoint
			if err := repo.Create(&waypoint); err != nil {
				t.Logf("Failed to create waypoint: %v", err)
				return false
			}

			// Query the waypoint by ID
			retrieved, err := repo.GetByID(waypoint.ID)
			if err != nil {
				t.Logf("Failed to retrieve waypoint: %v", err)
				return false
			}

			// Verify all required fields are preserved
			// 1. Name must match
			if retrieved.Name != waypoint.Name {
				t.Logf("Name mismatch: expected %s, got %s", waypoint.Name, retrieved.Name)
				return false
			}

			// 2. Longitude must match (with floating point tolerance)
			if !floatEquals(retrieved.Longitude, waypoint.Longitude, 0.0000001) {
				t.Logf("Longitude mismatch: expected %f, got %f", waypoint.Longitude, retrieved.Longitude)
				return false
			}

			// 3. Latitude must match (with floating point tolerance)
			if !floatEquals(retrieved.Latitude, waypoint.Latitude, 0.0000001) {
				t.Logf("Latitude mismatch: expected %f, got %f", waypoint.Latitude, retrieved.Latitude)
				return false
			}

			// 4. StayDuration must match
			if retrieved.StayDuration != waypoint.StayDuration {
				t.Logf("StayDuration mismatch: expected %d, got %d", waypoint.StayDuration, retrieved.StayDuration)
				return false
			}

			// 5. Notes must match
			if retrieved.Notes != waypoint.Notes {
				t.Logf("Notes mismatch: expected %s, got %s", waypoint.Notes, retrieved.Notes)
				return false
			}

			// 6. Address must match
			if retrieved.Address != waypoint.Address {
				t.Logf("Address mismatch: expected %s, got %s", waypoint.Address, retrieved.Address)
				return false
			}

			// 7. DayIndex must match
			if retrieved.DayIndex != waypoint.DayIndex {
				t.Logf("DayIndex mismatch: expected %d, got %d", waypoint.DayIndex, retrieved.DayIndex)
				return false
			}

			// 8. SortOrder must match
			if retrieved.SortOrder != waypoint.SortOrder {
				t.Logf("SortOrder mismatch: expected %d, got %d", waypoint.SortOrder, retrieved.SortOrder)
				return false
			}

			// 9. Budget must match (with floating point tolerance)
			if !floatEquals(retrieved.Budget, waypoint.Budget, 0.01) {
				t.Logf("Budget mismatch: expected %f, got %f", waypoint.Budget, retrieved.Budget)
				return false
			}

			// 10. RoadbookID must match
			if retrieved.RoadbookID != waypoint.RoadbookID {
				t.Logf("RoadbookID mismatch: expected %d, got %d", waypoint.RoadbookID, retrieved.RoadbookID)
				return false
			}

			return true
		},
		waypointGen,
	))

	properties.TestingRun(t)
}

// floatEquals compares two float64 values with a tolerance
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// **Feature: travel-planner, Property 5: 途经点按天分组正确性**
// **Validates: Requirements 4.3**
// For any roadbook with multi-day itinerary, after grouping by day,
// all waypoints in each group should have the same dayIndex value.
func TestWaypointGroupByDayCorrectness(t *testing.T) {
	db := setupWaypointTestDB(t)
	repo := NewWaypointRepositoryWithDB(db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for number of waypoints per day
	waypointsPerDayGen := gen.IntRange(1, 5)
	// Generator for number of days
	numDaysGen := gen.IntRange(1, 7)

	properties.Property("Waypoints grouped by day have consistent dayIndex", prop.ForAll(
		func(numDays int, waypointsPerDay int) bool {
			// Clean up previous test data
			db.Exec("DELETE FROM roadbook_waypoints")
			db.Exec("DELETE FROM roadbooks")

			// Create a test roadbook
			roadbookID, err := createTestRoadbook(db)
			if err != nil {
				t.Logf("Failed to create test roadbook: %v", err)
				return false
			}

			// Create waypoints for multiple days
			for day := 1; day <= numDays; day++ {
				for i := 0; i < waypointsPerDay; i++ {
					waypoint := model.RoadbookWaypoint{
						RoadbookID:   roadbookID,
						Name:         "Test Waypoint",
						Longitude:    116.397428 + float64(i)*0.01,
						Latitude:     39.90923 + float64(i)*0.01,
						DayIndex:     day,
						SortOrder:    i,
						StayDuration: 60,
						WaypointType: model.WaypointTypeWaypoint,
					}
					if err := repo.Create(&waypoint); err != nil {
						t.Logf("Failed to create waypoint: %v", err)
						return false
					}
				}
			}

			// Get waypoints grouped by day
			grouped, err := repo.ListByRoadbookIDGroupedByDay(roadbookID)
			if err != nil {
				t.Logf("Failed to get grouped waypoints: %v", err)
				return false
			}

			// Verify each group has consistent dayIndex
			for dayIndex, waypoints := range grouped {
				for _, wp := range waypoints {
					if wp.DayIndex != dayIndex {
						t.Logf("Inconsistent dayIndex in group %d: waypoint has dayIndex %d",
							dayIndex, wp.DayIndex)
						return false
					}
				}
			}

			// Verify we have the expected number of days
			if len(grouped) != numDays {
				t.Logf("Expected %d days, got %d", numDays, len(grouped))
				return false
			}

			// Verify each day has the expected number of waypoints
			for _, waypoints := range grouped {
				if len(waypoints) != waypointsPerDay {
					t.Logf("Expected %d waypoints per day, got %d", waypointsPerDay, len(waypoints))
					return false
				}
			}

			return true
		},
		numDaysGen,
		waypointsPerDayGen,
	))

	properties.TestingRun(t)
}
