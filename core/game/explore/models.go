package explore

type StageId string

func (id StageId) ToString() string {
	return string(id)
}
