package mapper

import (
	"github.com/RuneHistory/collector/internal/application/domain"
	"testing"
	"time"
)

func TestAccountToHttpV1(t *testing.T) {
	t.Parallel()
	acc := domain.NewHighScore("Jim", time.Now())
	mapped := HighScoreToHttpV1(acc)
	if acc.Nickname != mapped.Nickname {
		t.Errorf("Expecting %s, got %s", acc.Nickname, mapped.Nickname)
	}
}
