package filter

import (
	"log"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/mmcdole/gofeed"
)

func shouldInclude(matches bool, effect api.FilterSpec_FilterEffect) bool {
	switch effect {
	case api.FilterSpec_INCLUDE:
		return matches
	case api.FilterSpec_EXCLUDE:
		return !matches
	default:
		log.Printf("unknown filterEffect=%d, known=%v", effect, api.FilterSpec_FilterEffect_name)
		return true
	}
}

func FilterItems(feed *gofeed.Feed, filterSpecs []*api.FilterSpec) ([]*gofeed.Item, error) {
	filters, err := buildFilters(filterSpecs)
	if err != nil {
		return nil, err
	}
	var items []*gofeed.Item
	var include bool
	for _, item := range feed.Items {
		include = true
		for _, filter := range filters {
			include = shouldInclude(filter.Matches(item), filter.GetSpec().Effect)
			if !include {
				break
			}
		}
		if include {
			items = append(items, item)
		}
	}
	return items, nil
}
