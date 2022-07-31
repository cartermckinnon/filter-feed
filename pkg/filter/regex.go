package filter

import (
	"regexp"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/mmcdole/gofeed"
)

type regexFilter struct {
	regex *regexp.Regexp
	spec  *api.FilterSpec
}

func (f *regexFilter) Matches(item *gofeed.Item) bool {
	value := getTargetValue(f.spec.Target, item)
	if value == "" {
		return false
	}
	return f.regex.MatchString(value)
}

func (f *regexFilter) GetSpec() *api.FilterSpec {
	return f.spec
}

func getTargetValue(target api.FilterSpec_FilterTarget, item *gofeed.Item) string {
	switch target {
	case api.FilterSpec_TITLE:
		return item.Title
	case api.FilterSpec_DESCRIPTION:
		return item.Description
	}
	return ""
}
