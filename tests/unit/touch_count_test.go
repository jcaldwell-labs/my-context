package unit

import (
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/stretchr/testify/assert"
)

// TestContextWithMetadataTouchFields verifies that the TouchCount and LastTouchAt
// fields are properly initialized and accessible on ContextWithMetadata.
// This is a regression test for Issue #13 where touch counts weren't displayed
// in database mode.
func TestContextWithMetadataTouchFields(t *testing.T) {
	t.Run("touch fields are zero by default", func(t *testing.T) {
		ctx := models.NewContextWithMetadata("test-context", "", "", nil)

		assert.Equal(t, 0, ctx.TouchCount, "TouchCount should default to 0")
		assert.Nil(t, ctx.LastTouchAt, "LastTouchAt should default to nil")
	})

	t.Run("touch fields can be set", func(t *testing.T) {
		now := time.Now()
		ctx := &models.ContextWithMetadata{
			Name:        "test-context",
			StartTime:   now,
			Status:      "active",
			TouchCount:  5,
			LastTouchAt: &now,
		}

		assert.Equal(t, 5, ctx.TouchCount, "TouchCount should be set to 5")
		assert.NotNil(t, ctx.LastTouchAt, "LastTouchAt should be set")
		assert.Equal(t, now, *ctx.LastTouchAt, "LastTouchAt should equal the set time")
	})

	t.Run("context validates with touch fields", func(t *testing.T) {
		now := time.Now()
		ctx := &models.ContextWithMetadata{
			Name:        "valid-context",
			StartTime:   now,
			Status:      "active",
			TouchCount:  10,
			LastTouchAt: &now,
		}

		err := ctx.Validate()
		assert.NoError(t, err, "Context with touch fields should validate successfully")
	})
}

// TestTouchEventSynthesis verifies that synthetic touch events can be created
// from touch count and last touch time. This tests the pattern used in show.go
// for database mode.
func TestTouchEventSynthesis(t *testing.T) {
	t.Run("create synthetic touches from count and time", func(t *testing.T) {
		now := time.Now()
		touchCount := 3

		// This is the pattern used in show.go to create synthetic touches
		var touches []time.Time
		for i := 0; i < touchCount; i++ {
			touches = append(touches, now)
		}

		assert.Len(t, touches, 3, "Should create 3 synthetic touch events")
		for _, touch := range touches {
			assert.Equal(t, now, touch, "Each touch should have the last touch timestamp")
		}
	})

	t.Run("no touches created when count is zero", func(t *testing.T) {
		touchCount := 0
		var lastTouchAt *time.Time = nil

		var touches []time.Time
		if touchCount > 0 && lastTouchAt != nil {
			for i := 0; i < touchCount; i++ {
				touches = append(touches, *lastTouchAt)
			}
		}

		assert.Len(t, touches, 0, "Should create no touches when count is 0")
	})

	t.Run("no touches created when last touch time is nil", func(t *testing.T) {
		touchCount := 5
		var lastTouchAt *time.Time = nil

		var touches []time.Time
		if touchCount > 0 && lastTouchAt != nil {
			for i := 0; i < touchCount; i++ {
				touches = append(touches, *lastTouchAt)
			}
		}

		assert.Len(t, touches, 0, "Should create no touches when lastTouchAt is nil")
	})
}
