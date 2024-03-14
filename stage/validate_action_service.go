package stage

import "github.com/asragi/RinGo/core"

type ValidateActionFunc func(CheckIsPossibleArgs) core.IsPossible

func validateAction(
	checkIsExplorePossibleFunc checkIsExplorePossibleFunc,
	checkIsPossibleArgs CheckIsPossibleArgs,
) core.IsPossible {
	isPossible := checkIsExplorePossibleFunc(checkIsPossibleArgs)
	if !isPossible[core.PossibleTypeAll] {
		return false
	}
	return true
}

func CreateValidateAction(
	checkIsExplorePossibleFunc checkIsExplorePossibleFunc,
) ValidateActionFunc {
	return func(checkIsPossibleArgs CheckIsPossibleArgs) core.IsPossible {
		return validateAction(checkIsExplorePossibleFunc, checkIsPossibleArgs)
	}
}
