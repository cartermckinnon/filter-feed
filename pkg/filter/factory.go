package filter

import (
	"errors"
	"fmt"
	"regexp"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
)

func buildFilter(spec *api.FilterSpec) (Filter, error) {
	switch spec.Type {
	case api.FilterSpec_REGEX:
		regex, err := regexp.Compile(spec.Expression)
		if err != nil {
			return nil, err
		}
		return &regexFilter{
			regex: regex,
			spec:  spec,
		}, nil
	}
	return nil, ErrUnimplementedFilterType
}

var (
	ErrUnimplementedFilterType = errors.New("unimplemented filter type")
)

func buildFilters(specs []*api.FilterSpec) ([]Filter, error) {
	filters := make([]Filter, len(specs))
	for i, spec := range specs {
		filter, err := buildFilter(spec)
		if err != nil {
			return nil, fmt.Errorf("build filterSpec=%v: %w", spec, err)
		}
		filters[i] = filter
	}
	return filters, nil
}
