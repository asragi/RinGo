package stage

func UserExploreToIdArray(userExplores []UserExplore) []ExploreId {
	result := make([]ExploreId, len(userExplores))
	for i, explore := range userExplores {
		result[i] = explore.ExploreId
	}
	return result
}
