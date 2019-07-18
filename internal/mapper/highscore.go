package mapper

import (
	"github.com/RuneHistory/collector/internal/application/domain"
	"time"
)

type HighScoreTransport struct {
	Nickname  string                  `json:"nickname"`
	CreatedAt time.Time               `json:"created_at"`
	Skills    map[string]*SkillHttpV1 `json:"skills"`
}

func HighScoreFromTransport(hs *HighScoreTransport) *domain.HighScore {
	mapped := &domain.HighScore{
		Nickname:  hs.Nickname,
		CreatedAt: hs.CreatedAt,
		Skills:    make(map[string]*domain.Skill),
	}

	for name, skill := range hs.Skills {
		mapped.Skills[name] = SkillFromHttpV1(skill)
	}

	return mapped
}
