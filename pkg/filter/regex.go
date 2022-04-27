package filter

import (
	"regexp"
	"strings"

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

func getTargetValue(target string, item *gofeed.Item) string {
	switch strings.ToLower(target) {
	case "title":
		return item.Title
	case "description":
		return item.Description
	}
	return ""
}
