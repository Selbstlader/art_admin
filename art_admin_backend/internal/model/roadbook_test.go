package model

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: travel-planner, Property 1: 路书数据持久化 Round-Trip**
// **Validates: Requirements 3.2, 14.3**
// For any valid roadbook data, serializing to JSON and deserializing should produce equivalent data.
func TestRoadbookJSONRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for Roadbook
	roadbookGen := genRoadbook()

	properties.Property("Roadbook JSON round-trip preserves data", prop.ForAll(
		func(original Roadbook) bool {
			// Serialize to JSON
			jsonData, err := original.ToJSON()
			if err != nil {
				t.Logf("Serialization error: %v", err)
				return false
			}

			// Deserialize from JSON
			var restored Roadbook
			err = restored.FromJSON(jsonData)
			if err != nil {
				t.Logf("Deserialization error: %v", err)
				return false
			}

			// Compare essential fields (excluding time precision issues and associations)
			return compareRoadbooks(original, restored)
		},
		roadbookGen,
	))

	properties.TestingRun(t)
}

// **Feature: travel-planner, Property 1: 途经点数据持久化 Round-Trip**
// **Validates: Requirements 3.2, 14.3**
func TestWaypointJSONRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for Waypoint
	waypointGen := genWaypoint()

	properties.Property("Waypoint JSON round-trip preserves data", prop.ForAll(
		func(original RoadbookWaypoint) bool {
			// Serialize to JSON
			jsonData, err := original.ToJSON()
			if err != nil {
				t.Logf("Serialization error: %v", err)
				return false
			}

			// Deserialize from JSON
			var restored RoadbookWaypoint
			err = restored.FromJSON(jsonData)
			if err != nil {
				t.Logf("Deserialization error: %v", err)
				return false
			}

			// Compare essential fields
			return compareWaypoints(original, restored)
		},
		waypointGen,
	))

	properties.TestingRun(t)
}

// genRoadbook generates random Roadbook instances
func genRoadbook() gopter.Gen {
	return gopter.CombineGens(
		gen.Int64Range(1, 1000000), // ID
		gen.Int64Range(1, 1000000), // UserID
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 100 }), // Title
		gen.AlphaString(),  // Description
		gen.AlphaString(),  // CoverURL
		gen.IntRange(1, 3), // Visibility
		gen.OneConstOf("driving", "walking", "transit", "riding"), // TravelMode
		gen.IntRange(0, 1000000),                                  // TotalDistance
		gen.IntRange(0, 10000),                                    // TotalDuration
		gen.Float64Range(0, 100000),                               // TotalBudget
		gen.IntRange(0, 1000000),                                  // ViewCount
		gen.IntRange(0, 1000000),                                  // FavoriteCount
		gen.IntRange(0, 1000000),                                  // LikeCount
		gen.IntRange(0, 1000000),                                  // CommentCount
		gen.IntRange(0, 1000000),                                  // ShareCount
		gen.IntRange(1, 3),                                        // Status
	).Map(func(values []interface{}) Roadbook {
		now := time.Now().Truncate(time.Second)
		return Roadbook{
			ID:            values[0].(int64),
			UserID:        values[1].(int64),
			Title:         values[2].(string),
			Description:   values[3].(string),
			CoverURL:      values[4].(string),
			Visibility:    values[5].(int),
			TravelMode:    values[6].(string),
			TotalDistance: values[7].(int),
			TotalDuration: values[8].(int),
			TotalBudget:   values[9].(float64),
			ViewCount:     values[10].(int),
			FavoriteCount: values[11].(int),
			LikeCount:     values[12].(int),
			CommentCount:  values[13].(int),
			ShareCount:    values[14].(int),
			Status:        values[15].(int),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	})
}

// genWaypoint generates random RoadbookWaypoint instances
func genWaypoint() gopter.Gen {
	return gopter.CombineGens(
		gen.Int64Range(1, 1000000), // ID
		gen.Int64Range(1, 1000000), // RoadbookID
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 100 }), // Name
		gen.AlphaString(),           // Address
		gen.Float64Range(-180, 180), // Longitude
		gen.Float64Range(-90, 90),   // Latitude
		gen.AlphaString(),           // PoiID
		gen.AlphaString(),           // PoiType
		gen.IntRange(1, 30),         // DayIndex
		gen.IntRange(0, 100),        // SortOrder
		gen.IntRange(0, 1440),       // StayDuration (minutes in a day)
		gen.Float64Range(0, 10000),  // Budget
		gen.AlphaString(),           // Notes
		gen.IntRange(1, 3),          // WaypointType
	).Map(func(values []interface{}) RoadbookWaypoint {
		now := time.Now().Truncate(time.Second)
		return RoadbookWaypoint{
			ID:           values[0].(int64),
			RoadbookID:   values[1].(int64),
			Name:         values[2].(string),
			Address:      values[3].(string),
			Longitude:    values[4].(float64),
			Latitude:     values[5].(float64),
			PoiID:        values[6].(string),
			PoiType:      values[7].(string),
			DayIndex:     values[8].(int),
			SortOrder:    values[9].(int),
			StayDuration: values[10].(int),
			Budget:       values[11].(float64),
			Notes:        values[12].(string),
			WaypointType: values[13].(int),
			Images:       []string{},
			CreatedAt:    now,
			UpdatedAt:    now,
		}
	})
}

// compareRoadbooks compares two Roadbook instances for equality
func compareRoadbooks(a, b Roadbook) bool {
	if a.ID != b.ID {
		return false
	}
	if a.UserID != b.UserID {
		return false
	}
	if a.Title != b.Title {
		return false
	}
	if a.Description != b.Description {
		return false
	}
	if a.CoverURL != b.CoverURL {
		return false
	}
	if a.Visibility != b.Visibility {
		return false
	}
	if a.TravelMode != b.TravelMode {
		return false
	}
	if a.TotalDistance != b.TotalDistance {
		return false
	}
	if a.TotalDuration != b.TotalDuration {
		return false
	}
	if !floatEquals(a.TotalBudget, b.TotalBudget, 0.01) {
		return false
	}
	if a.ViewCount != b.ViewCount {
		return false
	}
	if a.FavoriteCount != b.FavoriteCount {
		return false
	}
	if a.LikeCount != b.LikeCount {
		return false
	}
	if a.CommentCount != b.CommentCount {
		return false
	}
	if a.ShareCount != b.ShareCount {
		return false
	}
	if a.Status != b.Status {
		return false
	}
	return true
}

// compareWaypoints compares two RoadbookWaypoint instances for equality
func compareWaypoints(a, b RoadbookWaypoint) bool {
	if a.ID != b.ID {
		return false
	}
	if a.RoadbookID != b.RoadbookID {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.Address != b.Address {
		return false
	}
	if !floatEquals(a.Longitude, b.Longitude, 0.0000001) {
		return false
	}
	if !floatEquals(a.Latitude, b.Latitude, 0.0000001) {
		return false
	}
	if a.PoiID != b.PoiID {
		return false
	}
	if a.PoiType != b.PoiType {
		return false
	}
	if a.DayIndex != b.DayIndex {
		return false
	}
	if a.SortOrder != b.SortOrder {
		return false
	}
	if a.StayDuration != b.StayDuration {
		return false
	}
	if !floatEquals(a.Budget, b.Budget, 0.01) {
		return false
	}
	if a.Notes != b.Notes {
		return false
	}
	if a.WaypointType != b.WaypointType {
		return false
	}
	if !reflect.DeepEqual(a.Images, b.Images) {
		return false
	}
	return true
}

// floatEquals compares two floats with a tolerance
func floatEquals(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

// TestRoadbookJSONRoundTripWithWaypoints tests round-trip with nested waypoints
func TestRoadbookJSONRoundTripWithWaypoints(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	// Generator for Roadbook with Waypoints
	roadbookWithWaypointsGen := genRoadbookWithWaypoints()

	properties.Property("Roadbook with waypoints JSON round-trip preserves data", prop.ForAll(
		func(original Roadbook) bool {
			// Serialize to JSON
			jsonData, err := json.Marshal(original)
			if err != nil {
				t.Logf("Serialization error: %v", err)
				return false
			}

			// Deserialize from JSON
			var restored Roadbook
			err = json.Unmarshal(jsonData, &restored)
			if err != nil {
				t.Logf("Deserialization error: %v", err)
				return false
			}

			// Compare roadbook
			if !compareRoadbooks(original, restored) {
				return false
			}

			// Compare waypoints count
			if len(original.Waypoints) != len(restored.Waypoints) {
				return false
			}

			// Compare each waypoint
			for i := range original.Waypoints {
				if !compareWaypoints(original.Waypoints[i], restored.Waypoints[i]) {
					return false
				}
			}

			return true
		},
		roadbookWithWaypointsGen,
	))

	properties.TestingRun(t)
}

// genSimpleWaypoint generates simple RoadbookWaypoint instances without SuchThat constraints
func genSimpleWaypoint() gopter.Gen {
	return gopter.CombineGens(
		gen.Int64Range(1, 1000000),  // ID
		gen.Int64Range(1, 1000000),  // RoadbookID
		gen.Identifier(),            // Name (always non-empty)
		gen.Identifier(),            // Address
		gen.Float64Range(-180, 180), // Longitude
		gen.Float64Range(-90, 90),   // Latitude
		gen.Identifier(),            // PoiID
		gen.Identifier(),            // PoiType
		gen.IntRange(1, 30),         // DayIndex
		gen.IntRange(0, 100),        // SortOrder
		gen.IntRange(0, 1440),       // StayDuration (minutes in a day)
		gen.Float64Range(0, 10000),  // Budget
		gen.Identifier(),            // Notes
		gen.IntRange(1, 3),          // WaypointType
	).Map(func(values []interface{}) RoadbookWaypoint {
		now := time.Now().Truncate(time.Second)
		return RoadbookWaypoint{
			ID:           values[0].(int64),
			RoadbookID:   values[1].(int64),
			Name:         values[2].(string),
			Address:      values[3].(string),
			Longitude:    values[4].(float64),
			Latitude:     values[5].(float64),
			PoiID:        values[6].(string),
			PoiType:      values[7].(string),
			DayIndex:     values[8].(int),
			SortOrder:    values[9].(int),
			StayDuration: values[10].(int),
			Budget:       values[11].(float64),
			Notes:        values[12].(string),
			WaypointType: values[13].(int),
			Images:       []string{},
			CreatedAt:    now,
			UpdatedAt:    now,
		}
	})
}

// genSimpleRoadbook generates simple Roadbook instances without SuchThat constraints
func genSimpleRoadbook() gopter.Gen {
	return gopter.CombineGens(
		gen.Int64Range(1, 1000000), // ID
		gen.Int64Range(1, 1000000), // UserID
		gen.Identifier(),           // Title (always non-empty)
		gen.Identifier(),           // Description
		gen.Identifier(),           // CoverURL
		gen.IntRange(1, 3),         // Visibility
		gen.OneConstOf("driving", "walking", "transit", "riding"), // TravelMode
		gen.IntRange(0, 1000000),                                  // TotalDistance
		gen.IntRange(0, 10000),                                    // TotalDuration
		gen.Float64Range(0, 100000),                               // TotalBudget
		gen.IntRange(0, 1000000),                                  // ViewCount
		gen.IntRange(0, 1000000),                                  // FavoriteCount
		gen.IntRange(0, 1000000),                                  // LikeCount
		gen.IntRange(0, 1000000),                                  // CommentCount
		gen.IntRange(0, 1000000),                                  // ShareCount
		gen.IntRange(1, 3),                                        // Status
	).Map(func(values []interface{}) Roadbook {
		now := time.Now().Truncate(time.Second)
		return Roadbook{
			ID:            values[0].(int64),
			UserID:        values[1].(int64),
			Title:         values[2].(string),
			Description:   values[3].(string),
			CoverURL:      values[4].(string),
			Visibility:    values[5].(int),
			TravelMode:    values[6].(string),
			TotalDistance: values[7].(int),
			TotalDuration: values[8].(int),
			TotalBudget:   values[9].(float64),
			ViewCount:     values[10].(int),
			FavoriteCount: values[11].(int),
			LikeCount:     values[12].(int),
			CommentCount:  values[13].(int),
			ShareCount:    values[14].(int),
			Status:        values[15].(int),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	})
}

// genRoadbookWithWaypoints generates Roadbook with random waypoints
func genRoadbookWithWaypoints() gopter.Gen {
	return gopter.CombineGens(
		genSimpleRoadbook(),
		gen.SliceOfN(5, genSimpleWaypoint()),
	).Map(func(values []interface{}) Roadbook {
		roadbook := values[0].(Roadbook)
		waypoints := values[1].([]RoadbookWaypoint)

		// Assign roadbook ID to waypoints
		for i := range waypoints {
			waypoints[i].RoadbookID = roadbook.ID
		}
		roadbook.Waypoints = waypoints
		return roadbook
	})
}
