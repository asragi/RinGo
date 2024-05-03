package explore

type StageId string

func CreateStageId(id string) (StageId, error) {
	return StageId(id), nil
}

func (id StageId) String() string {
	return string(id)
}
