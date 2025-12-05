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

// setupCommentTestDB creates an in-memory SQLite database for comment testing
func setupCommentTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&model.Roadbook{}, &model.RoadbookComment{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// genCommentContent generates valid comment content (non-empty, <= 500 chars)
func genCommentContent() gopter.Gen {
	return gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0 && len(s) <= 500
	})
}

// genCommentForDB generates random RoadbookComment instances suitable for database insertion
func genCommentForDB(roadbookID, userID int64) gopter.Gen {
	return gopter.CombineGens(
		genCommentContent(), // Content
	).Map(func(values []interface{}) model.RoadbookComment {
		return model.RoadbookComment{
			RoadbookID: roadbookID,
			UserID:     userID,
			ParentID:   0,
			Content:    values[0].(string),
		}
	})
}

// **Feature: travel-planner, Property 7: 评论数据持久化 Round-Trip**
// **Validates: Requirements 5.1**
// For any valid comment, saving to database and querying by ID should return equivalent data.
func TestCommentRoundTrip(t *testing.T) {
	db := setupCommentTestDB(t)
	repo := NewCommentRepositoryWithDB(db)

	// Create a test roadbook first
	roadbook := &model.Roadbook{
		UserID:     1,
		Title:      "Test Roadbook",
		Visibility: model.VisibilityPublic,
		TravelMode: "driving",
		Status:     model.RoadbookStatusPublished,
	}
	if err := db.Create(roadbook).Error; err != nil {
		t.Fatalf("Failed to create test roadbook: %v", err)
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for comment data
	commentGen := genCommentForDB(roadbook.ID, 1)

	properties.Property("Comment round-trip preserves data", prop.ForAll(
		func(comment model.RoadbookComment) bool {
			// Create the comment
			if err := repo.Create(&comment); err != nil {
				t.Logf("Failed to create comment: %v", err)
				return false
			}

			// Query by ID
			retrieved, err := repo.GetByID(comment.ID)
			if err != nil {
				t.Logf("Failed to retrieve comment: %v", err)
				return false
			}

			// Verify data consistency
			if retrieved.RoadbookID != comment.RoadbookID {
				t.Logf("RoadbookID mismatch: expected %d, got %d", comment.RoadbookID, retrieved.RoadbookID)
				return false
			}
			if retrieved.UserID != comment.UserID {
				t.Logf("UserID mismatch: expected %d, got %d", comment.UserID, retrieved.UserID)
				return false
			}
			if retrieved.Content != comment.Content {
				t.Logf("Content mismatch: expected %s, got %s", comment.Content, retrieved.Content)
				return false
			}
			if retrieved.ParentID != comment.ParentID {
				t.Logf("ParentID mismatch: expected %d, got %d", comment.ParentID, retrieved.ParentID)
				return false
			}

			return true
		},
		commentGen,
	))

	properties.TestingRun(t)
}

// **Feature: travel-planner, Property 8: 评论列表按时间倒序排列**
// **Validates: Requirements 5.2**
// For any list of comments, the list should be ordered by created_at in descending order.
func TestCommentListOrderedByCreatedAtDesc(t *testing.T) {
	db := setupCommentTestDB(t)
	repo := NewCommentRepositoryWithDB(db)

	// Create a test roadbook first
	roadbook := &model.Roadbook{
		UserID:     1,
		Title:      "Test Roadbook",
		Visibility: model.VisibilityPublic,
		TravelMode: "driving",
		Status:     model.RoadbookStatusPublished,
	}
	if err := db.Create(roadbook).Error; err != nil {
		t.Fatalf("Failed to create test roadbook: %v", err)
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for number of comments to create
	countGen := gen.IntRange(2, 20)

	properties.Property("Comment list is ordered by created_at DESC", prop.ForAll(
		func(count int) bool {
			// Clean up previous test data
			db.Exec("DELETE FROM roadbook_comments")

			// Create comments with different timestamps
			for i := 0; i < count; i++ {
				comment := model.RoadbookComment{
					RoadbookID: roadbook.ID,
					UserID:     1,
					Content:    "Test Comment",
					CreatedAt:  time.Now().Add(time.Duration(i) * time.Second),
				}
				if err := repo.Create(&comment); err != nil {
					t.Logf("Failed to create comment: %v", err)
					return false
				}
				// Small delay to ensure different timestamps
				time.Sleep(time.Millisecond * 10)
			}

			// Query the list
			comments, _, err := repo.ListByRoadbookID(roadbook.ID, 1, 100)
			if err != nil {
				t.Logf("Failed to list comments: %v", err)
				return false
			}

			// Verify ordering: each comment's CreatedAt should be >= the next one's
			for i := 0; i < len(comments)-1; i++ {
				if comments[i].CreatedAt.Before(comments[i+1].CreatedAt) {
					t.Logf("Ordering violation at index %d: %v < %v",
						i, comments[i].CreatedAt, comments[i+1].CreatedAt)
					return false
				}
			}

			return true
		},
		countGen,
	))

	properties.TestingRun(t)
}

// **Feature: travel-planner, Property 9: 评论父子关系有效性**
// **Validates: Requirements 5.3**
// For any reply comment (parentId > 0), its parentId should point to an existing comment in the same roadbook.
func TestCommentParentChildRelationship(t *testing.T) {
	db := setupCommentTestDB(t)
	repo := NewCommentRepositoryWithDB(db)

	// Create a test roadbook first
	roadbook := &model.Roadbook{
		UserID:     1,
		Title:      "Test Roadbook",
		Visibility: model.VisibilityPublic,
		TravelMode: "driving",
		Status:     model.RoadbookStatusPublished,
	}
	if err := db.Create(roadbook).Error; err != nil {
		t.Fatalf("Failed to create test roadbook: %v", err)
	}

	// Create another roadbook for cross-roadbook testing
	otherRoadbook := &model.Roadbook{
		UserID:     1,
		Title:      "Other Roadbook",
		Visibility: model.VisibilityPublic,
		TravelMode: "driving",
		Status:     model.RoadbookStatusPublished,
	}
	if err := db.Create(otherRoadbook).Error; err != nil {
		t.Fatalf("Failed to create other roadbook: %v", err)
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for number of parent comments
	countGen := gen.IntRange(1, 10)

	properties.Property("Reply comment's parent exists in same roadbook", prop.ForAll(
		func(parentCount int) bool {
			// Clean up previous test data
			db.Exec("DELETE FROM roadbook_comments")

			// Create parent comments
			var parentIDs []int64
			for i := 0; i < parentCount; i++ {
				parent := model.RoadbookComment{
					RoadbookID: roadbook.ID,
					UserID:     1,
					Content:    "Parent Comment",
					ParentID:   0,
				}
				if err := repo.Create(&parent); err != nil {
					t.Logf("Failed to create parent comment: %v", err)
					return false
				}
				parentIDs = append(parentIDs, parent.ID)
			}

			// Create reply comments for each parent
			for _, parentID := range parentIDs {
				reply := model.RoadbookComment{
					RoadbookID: roadbook.ID,
					UserID:     1,
					Content:    "Reply Comment",
					ParentID:   parentID,
				}
				if err := repo.Create(&reply); err != nil {
					t.Logf("Failed to create reply comment: %v", err)
					return false
				}

				// Verify parent exists and belongs to same roadbook
				exists, err := repo.ExistsByIDAndRoadbookID(parentID, roadbook.ID)
				if err != nil {
					t.Logf("Error checking parent existence: %v", err)
					return false
				}
				if !exists {
					t.Logf("Parent comment should exist in same roadbook")
					return false
				}

				// Verify parent does NOT exist in other roadbook
				existsInOther, err := repo.ExistsByIDAndRoadbookID(parentID, otherRoadbook.ID)
				if err != nil {
					t.Logf("Error checking parent in other roadbook: %v", err)
					return false
				}
				if existsInOther {
					t.Logf("Parent comment should NOT exist in other roadbook")
					return false
				}
			}

			// Verify tree structure
			tree, _, err := repo.GetCommentTree(roadbook.ID, 1, 100)
			if err != nil {
				t.Logf("Failed to get comment tree: %v", err)
				return false
			}

			// Each top-level comment should have replies
			for _, parent := range tree {
				if parent.ParentID != 0 {
					t.Logf("Top-level comment should have ParentID=0")
					return false
				}
				// Verify all replies belong to this parent
				for _, reply := range parent.Replies {
					if reply.ParentID != parent.ID {
						t.Logf("Reply's ParentID should match parent's ID")
						return false
					}
					if reply.RoadbookID != parent.RoadbookID {
						t.Logf("Reply should belong to same roadbook as parent")
						return false
					}
				}
			}

			return true
		},
		countGen,
	))

	properties.TestingRun(t)
}

// **Feature: travel-planner, Property 10: 评论内容验证**
// **Validates: Requirements 5.5**
// For any comment content, if length > 500 chars or only whitespace, validation should fail.
func TestCommentContentValidation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for valid content (non-empty, <= 500 chars)
	validContentGen := gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0 && len(s) <= 500
	})

	// Generator for whitespace-only content
	whitespaceGen := gen.OneConstOf("", " ", "  ", "\t", "\n", "   \t\n   ")

	// Generator for content > 500 chars
	longContentGen := gen.SliceOfN(501, gen.AlphaChar()).Map(func(chars []rune) string {
		return string(chars)
	})

	properties.Property("Valid content passes validation", prop.ForAll(
		func(content string) bool {
			return model.ValidateCommentContent(content)
		},
		validContentGen,
	))

	properties.Property("Whitespace-only content fails validation", prop.ForAll(
		func(content string) bool {
			return !model.ValidateCommentContent(content)
		},
		whitespaceGen,
	))

	properties.Property("Content > 500 chars fails validation", prop.ForAll(
		func(content string) bool {
			if len([]rune(content)) <= 500 {
				// Skip if generator didn't produce long enough content
				return true
			}
			return !model.ValidateCommentContent(content)
		},
		longContentGen,
	))

	properties.TestingRun(t)
}
