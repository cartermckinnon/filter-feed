package override

import (
	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/mmcdole/gofeed"
)

type overrideImpl struct {
	spec *api.OverrideSpec
}

func (mo *overrideImpl) Apply(feed *gofeed.Feed) bool {
	switch mo.spec.Target {
	case api.OverrideSpec_TITLE:
		if feed.Title != mo.spec.Value {
			feed.Title = mo.spec.Value
			return true
		}
	case api.OverrideSpec_SUBTITLE:
		if feed.ITunesExt.Subtitle != mo.spec.Value {
			feed.ITunesExt.Subtitle = mo.spec.Value
			return true
		}
	case api.OverrideSpec_DESCRIPTION:
		if feed.Description != mo.spec.Value {
			feed.Description = mo.spec.Value
			return true
		}
	}
	return false
}

func (f *overrideImpl) GetSpec() *api.OverrideSpec {
	return f.spec
}
