package toggle

import (
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/feature_toggling/domain"
)

type Toggle struct {
	cfg *config.Config
}

func NewToggle(cfg *config.Config) *Toggle {
	return &Toggle{
		cfg: cfg,
	}
}

func (t *Toggle) IsEnabled(userId int, feature domain.FeatureName) bool {
	if feature == domain.Postpone {
		for _, id := range t.cfg.Features.Postpone.Ids {
			if id == userId {
				return true
			}
		}
	}

	return false
}
