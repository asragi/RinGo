package game

import (
	"github.com/asragi/RinGo/core"
	"math"
)

type ConsumptionProb float32

type GainingPoint int

func (g GainingPoint) Multiply(num int) GainingPoint {
	value := int(g)
	return GainingPoint(value * num)
}

func (g GainingPoint) ApplyTo(exp core.SkillExp) core.SkillExp {
	return exp + core.SkillExp(g)
}

type ExploreId string

type StaminaReducibleRate float64

func ApplyReduction(s core.Stamina, reductionRate float64, reducibleRate StaminaReducibleRate) core.Stamina {
	constStamina := float64(s) * (1.0 - float64(reducibleRate))
	varyStamina := float64(s) * reductionRate * float64(reducibleRate)
	staminaRounded := int(math.Max(1, math.Round(constStamina+varyStamina)))
	return core.Stamina(staminaRounded)
}

type EarningProb float32
