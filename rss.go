package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var rssFeed RSSFeed
	rssFeed.Channel.Items = make([]RSSItem, 0)
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return nil, err
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Link = html.UnescapeString(rssFeed.Channel.Link)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, ri := range rssFeed.Channel.Items {
		rssFeed.Channel.Items[i].Title = html.UnescapeString(ri.Title)
		rssFeed.Channel.Items[i].Link = html.UnescapeString(ri.Link)
		rssFeed.Channel.Items[i].Description = html.UnescapeString(ri.Description)
		rssFeed.Channel.Items[i].PubDate = html.UnescapeString(ri.PubDate)
	}
	return &rssFeed, nil
}
