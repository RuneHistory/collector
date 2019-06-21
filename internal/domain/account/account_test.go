package account

import (
	"testing"
	"time"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()
	id := "my-uuid"
	bucketId := "my-bucket-uuid"
	collectedAt := time.Now()
	acc := NewAccount(id, bucketId, collectedAt)
	if id != acc.ID {
		t.Errorf("expected account ID to equal %s, got %s", id, acc.ID)
	}
	if bucketId != acc.BucketID {
		t.Errorf("expected account bucketId to equal %s, got %s", bucketId, acc.BucketID)
	}
	if collectedAt != acc.CollectedAt {
		t.Errorf("expected account collectedAt to equal %s, got %s", collectedAt, acc.CollectedAt)
	}
}
