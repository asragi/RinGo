package core

type DecideInitialName func() Name
type DecideInitialShopName func() Name

// TODO: initial name must be decided depending on locale
func CreateDecideInitialName() DecideInitialName {
	return func() Name {
		return "夢追い人"
	}
}

func CreateDecideInitialShopName() DecideInitialShopName {
	return func() Name {
		return "夢追い人の店"
	}
}
