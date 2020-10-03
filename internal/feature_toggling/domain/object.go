package domain

const (
	Postpone FeatureName = `postpone_button`
)

type (
	FeatureName   string
	FeatureToggle interface {
		IsEnabled(userId int, feature FeatureName) bool
	}

	Feature struct {
		Name string
	}
)
