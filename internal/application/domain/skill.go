package domain

type Skill struct {
	Name       string
	Rank       int
	Level      int
	Experience int
}

func NewSkill(name string, rank int, level int, experience int) *Skill {
	return &Skill{
		Name:       name,
		Rank:       rank,
		Level:      level,
		Experience: experience,
	}
}
