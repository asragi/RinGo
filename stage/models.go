package stage

import "github.com/asragi/RinGo/core"

// item
type ConsumptionProb float32

// skill
type GainingPoint int

func (g GainingPoint) Multiply(num int) GainingPoint {
	value := int(g)
	return GainingPoint(value * num)
}

func (g GainingPoint) ApplyTo(exp core.SkillExp) core.SkillExp {
	return exp + core.SkillExp(g)
}

// explore
type ExploreId string

// stage
type StageId string
