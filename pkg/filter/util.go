package filter

import (
	"log"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/mmcdole/gofeed"
)

func shouldInclude(matches bool, effect api.FilterEffect) bool {
	switch effect {
	case api.FilterEffect_INCLUDE:
		return matches
	case api.FilterEffect_EXCLUDE:
		return !matches
	default:
		log.Printf("unknown filterEffect=%d, known=%v", effect, api.FilterEffect_name)
		return true
	}
}

func FilterItems(feed gofeed.Feed, filterSpecs []*api.FilterSpec) ([]*gofeed.Item, error) {
	filters, err := buildFilters(filterSpecs)
	if err != nil {
		return nil, err
	}
	var items []*gofeed.Item
	for _, item := range feed.Items {
		for _, filter := range filters {
			if shouldInclude(filter.Matches(item), filter.GetSpec().Effect) {
				items = append(items, item)
				break
			}
		}
	}
	return items, nil
}
