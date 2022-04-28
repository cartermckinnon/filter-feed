package util

import (
	"fmt"
	"log"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

func ConvertFeedToString(feed *gofeed.Feed) (string, error) {
	gfeed := convertToGorillaFeed(feed)
	switch feed.FeedType {
	case "rss":
		return gfeed.ToRss()
	case "atom":
		return gfeed.ToAtom()
	default:
		return "", fmt.Errorf("unknown feedType=%s", feed.FeedType)
	}
}

func convertToGorillaFeed(feed *gofeed.Feed) *feeds.Feed {
	gfeed := &feeds.Feed{
		Title:       feed.Title,
		Link:        &feeds.Link{Href: feed.Link},
		Description: feed.Description,
	}
	if feed.Author != nil {
		author := feeds.Author{}
		author.Name = feed.Author.Name
		author.Email = feed.Author.Email
		gfeed.Author = &author
	}

	var items []*feeds.Item
	for _, item := range feed.Items {
		gitem := &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: item.Link},
			Description: item.Description,
			Created:     *item.PublishedParsed,
		}
		if item.Author != nil {
			gitem.Author = &feeds.Author{Name: item.Author.Name, Email: item.Author.Email}
		}
		if len(item.Enclosures) > 0 {
			e := item.Enclosures[0]
			gitem.Enclosure = &feeds.Enclosure{
				Url:    e.URL,
				Type:   e.Type,
				Length: e.Length,
			}
			if len(item.Enclosures) > 1 {
				log.Printf("multiple enclosures for item=%s", item.Title)
			}
		}
		items = append(items, gitem)
	}
	gfeed.Items = items

	return gfeed
}
