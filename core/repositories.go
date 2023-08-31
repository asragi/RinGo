package core

type GetItemMasterRes struct {
	ItemId      ItemId
	Price       Price
	DisplayName DisplayName
	Description Description
	MaxStock    MaxStock
}

type ItemMasterRepo interface {
	Get(ItemId) (GetItemMasterRes, error)
}

type GetItemStorageRes struct {
	UserId UserId
	Stock  Stock
}

type ItemData struct {
	UserId UserId
	ItemId ItemId
	Stock  Stock
}

type BatchGetStorageRes struct {
	UserId   UserId
	ItemData []ItemData
}

type ItemStorageRepo interface {
	Get(UserId, ItemId, AccessToken) (GetItemStorageRes, error)
	BatchGet(UserId, []ItemId, AccessToken) (BatchGetStorageRes, error)
}

// Skill
type SkillMaster struct {
	SkillId     SkillId
	DisplayName DisplayName
}

type BatchGetSkillMasterRes struct {
	Skills []SkillMaster
}

type SkillMasterRepo interface {
	BatchGet([]SkillId) (BatchGetSkillMasterRes, error)
}

type UserSkillRes struct {
	UserId  UserId
	SkillId SkillId
	SkillLv SkillLv
}

type BatchGetUserSkillRes struct {
	UserId UserId
	Skills []UserSkillRes
}

type UserSkillRepo interface {
	BatchGet(UserId, []SkillId, AccessToken) (BatchGetUserSkillRes, error)
}

// Explore
type GetAllExploreMasterRes struct {
	ExploreId   ExploreId
	DisplayName DisplayName
	Description Description
}

type ExploreMasterRepo interface {
	GetAllExploreMaster(ItemId) ([]GetAllExploreMasterRes, error)
}

type ExploreUserData struct {
	ExploreId ExploreId
	IsKnown   IsKnown
}

type GetActionsRes struct {
	UserId   UserId
	Explores []ExploreUserData
}

type UserExploreRepo interface {
	GetActions(UserId, []ExploreId, AccessToken) (GetActionsRes, error)
}

type Condition struct {
	ConditionId          ConditionId
	ConditionType        ConditionType
	ConditionTargetId    ConditionTargetId
	ConditionTargetValue ConditionTargetValue
}

type ExploreConditions struct {
	ExploreId  ExploreId
	Conditions []Condition
}

type GetAllConditionsRes struct {
	Explores []ExploreConditions
}

type ConditionRepo interface {
	GetAllConditions([]ExploreId) (GetAllConditionsRes, error)
}
