package domain

import (
	"testing"
	"time"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()
	uuid := "my-uuid"
	bucketUuid := "bucket-uuid"
	nickname := "Jim"
	now := time.Now()
	acc := NewAccount(uuid, bucketUuid, nickname, now)
	if uuid != acc.ID {
		t.Errorf("expected account ID to equal %s, got %s", uuid, acc.ID)
	}
	if bucketUuid != acc.BucketID {
		t.Errorf("expected account BucketID to equal %s, got %s", bucketUuid, acc.BucketID)
	}
	if nickname != acc.Nickname {
		t.Errorf("expected account nickname to equal %s, got %s", nickname, acc.Nickname)
	}
	if now != acc.CreatedAt {
		t.Errorf("expected account created at to equal %s, got %s", now, acc.CreatedAt)
	}
}
