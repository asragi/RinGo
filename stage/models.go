package stage

// item
type ConsumptionProb float32

// skill
type GainingPoint int

func (g GainingPoint) Multiply(num int) GainingPoint {
	value := int(g)
	return GainingPoint(value * num)
}

// explore
type ExploreId string

// condition
type ConditionId string
type ConditionType string

const (
	ConditionTypeItem  = ConditionType("item")
	ConditionTypeSkill = ConditionType("skill")
)

type ConditionTargetId string
type ConditionTargetValue int

// stage
type StageId string
