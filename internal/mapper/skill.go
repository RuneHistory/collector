package mapper

import "github.com/RuneHistory/collector/internal/application/domain"

type SkillHttpV1 struct {
	Name       string `json:"name"`
	Rank       int    `json:"rank"`
	Level      int    `json:"level"`
	Experience int    `json:"experience"`
}

func SkillFromHttpV1(s *SkillHttpV1) *domain.Skill {
	return domain.NewSkill(s.Name, s.Rank, s.Level, s.Experience)
}
