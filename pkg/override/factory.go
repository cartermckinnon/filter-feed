package override

import (
	"fmt"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
)

func buildOverride(spec *api.OverrideSpec) (Override, error) {
	return &overrideImpl{
		spec: spec,
	}, nil
}

func buildOverrides(specs []*api.OverrideSpec) ([]Override, error) {
	overrides := make([]Override, len(specs))
	for i, spec := range specs {
		o, err := buildOverride(spec)
		if err != nil {
			return nil, fmt.Errorf("build overrideSpec=%v: %w", spec, err)
		}
		overrides[i] = o
	}
	return overrides, nil
}
