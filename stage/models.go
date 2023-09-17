package stage

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
