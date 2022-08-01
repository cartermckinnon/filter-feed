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
		gfeed.Author = &feeds.Author{
			Name:  feed.Author.Name,
			Email: feed.Author.Email,
		}
	}
	if feed.Image != nil {
		gfeed.Image = &feeds.Image{
			Title: feed.Image.Title,
			Url:   feed.Image.URL,
		}
	}

	var items []*feeds.Item
	for _, item := range feed.Items {
		gitem := &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: item.Link},
			Description: item.Description,
			Created:     *item.PublishedParsed,
			Content:     item.Content,
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
