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

type ActionId string

func NewActionId(id string) (ActionId, error) {
	return ActionId(id), nil
}

func (id ActionId) String() string {
	return string(id)
}

type StaminaReducibleRate float64

func ApplyReduction(s core.StaminaCost, reductionRate float64, reducibleRate StaminaReducibleRate) core.StaminaCost {
	constStamina := float64(s) * (1.0 - float64(reducibleRate))
	varyStamina := float64(s) * reductionRate * float64(reducibleRate)
	staminaRounded := int(math.Max(1, math.Round(constStamina+varyStamina)))
	return core.StaminaCost(staminaRounded)
}

type EarningProb float32

type PricePenalty float32

func NewPricePenalty(basePrice core.Price) PricePenalty {
	// 100 -> 1, 10000 -> 2, 1000000 -> 3
	return PricePenalty(math.Log10(float64(basePrice)) / 2)
}
