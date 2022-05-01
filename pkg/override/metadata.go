package override

import (
	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/mmcdole/gofeed"
)

type metadataOverride struct {
	spec *api.OverrideSpec
}

func (mo *metadataOverride) Apply(feed *gofeed.Feed) bool {
	switch mo.spec.Target {
	case "title":
		if feed.Title != mo.spec.Value {
			feed.Title = mo.spec.Value
			return true
		}
	case "description":
		if feed.Description != mo.spec.Value {
			feed.Description = mo.spec.Value
			return true
		}
	case "subtitle":
		if feed.ITunesExt.Subtitle != mo.spec.Value {
			feed.ITunesExt.Subtitle = mo.spec.Value
			return true
		}
	}
	return false
}

func (f *metadataOverride) GetSpec() *api.OverrideSpec {
	return f.spec
}
