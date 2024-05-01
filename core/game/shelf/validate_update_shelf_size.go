package shelf

import (
	"fmt"
	"github.com/asragi/RinGo/core/game"
)

type ValidateUpdateShelfSizeFunc func(targetSize Size, currentSize Size) error

func validateUpdateShelfSize(size Size, currentSize Size) error {
	if !size.ValidSize() {
		return fmt.Errorf("invalid shelf size: %d :%w", size, game.InvalidActionError)
	}
	if size.Equals(currentSize) {
		return fmt.Errorf("shelf size is already %d: %w", size, game.InvalidActionError)
	}
	return nil
}
