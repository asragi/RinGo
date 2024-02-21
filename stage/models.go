package stage

import (
	"math"

	"github.com/asragi/RinGo/core"
)

// item
type ConsumptionProb float32
type EarningProb float32

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

func (id ExploreId) ToString() string {
	return string(id)
}

type StaminaReducibleRate float64

func ApplyReduction(s core.Stamina, reductionRate float64, reducibleRate StaminaReducibleRate) core.Stamina {
	constStamina := float64(s) * (1.0 - float64(reducibleRate))
	varyStamina := float64(s) * reductionRate * float64(reducibleRate)
	staminaRounded := int(math.Max(1, math.Round(constStamina+varyStamina)))
	return core.Stamina(staminaRounded)
}

// stage
type StageId string

func (id StageId) ToString() string {
	return string(id)
}
