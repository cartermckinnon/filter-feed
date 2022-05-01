package override

import (
	"errors"
	"fmt"
	"strings"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
)

func buildOverride(spec *api.OverrideSpec) (Override, error) {
	target := strings.ToLower(spec.Target)
	if strings.HasPrefix(target, "metadata.") {
		return &metadataOverride{
			spec: spec,
		}, nil
	}
	return nil, ErrUnknownOverrideTarget
}

var (
	ErrUnknownOverrideTarget = errors.New("unknown override target")
)

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
