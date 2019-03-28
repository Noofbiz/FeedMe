package model

import (
	"database/sql"
	"time"
)

type ItemInfo struct {
	Title, AuthorName, AuthorEmail, Content, Description, Link string
	Read                                                       bool
	Published                                                  time.Time
}

type FeedData struct {
	FeedURLs   []string
	FeedTitles []string
	Items      []ItemInfo
}

func (f *feedDatabase) GetFeedData(origin, search string) (fd FeedData, err error) {
	allFeedsRows, err := f.db.Query("SELECT name, url FROM all_feeds")
	if err != nil {
		return fd, err
	}

	for allFeedsRows.Next() {
		var name, url string
		allFeedsRows.Scan(&name, &url)
		fd.FeedTitles = append(fd.FeedTitles, name)
		fd.FeedURLs = append(fd.FeedURLs, url)
	}
	allFeedsRows.Close()

	var feedItems *sql.Rows
	if search == "" && origin == "" {
		feedItems, err = f.db.Query("SELECT title, author_name, author_email, published, content, description, link, read FROM feed_items")
		if err != nil {
			return fd, err
		}
	} else if search == "" {
		feedItems, err = f.db.Query("SELECT title, author_name, author_email, published, content, description, link, read FROM feed_items WHERE origin=?", origin)
		if err != nil {
			return fd, err
		}
	} else if origin == "" {
		feedItems, err = f.db.Query("SELECT title, author_name, author_email, published, content, description, link, read FROM feed_items WHERE feed_items MATCH ?", search)
		if err != nil {
			return fd, err
		}
	} else {
		feedItems, err = f.db.Query("SELECT title, author_name, author_email, published, content, description, link, read FROM feed_items WHERE origin=? AND feed_items MATCH ?", origin, search)
		if err != nil {
			return fd, err
		}
	}

	for feedItems.Next() {
		item := ItemInfo{}
		var pub string
		var read int
		feedItems.Scan(&item.Title, &item.AuthorName, &item.AuthorEmail, &pub, &item.Content, &item.Description, &item.Link, &read)
		item.Published, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", pub)
		if err != nil {
			return fd, err
		}
		if read == 1 {
			item.Read = true
		}
		fd.Items = append(fd.Items, item)
	}
	feedItems.Close()

	return fd, err
}

func (f *feedDatabase) MarkAsRead(title string) error {
	_, err := f.db.Exec("UPDATE feed_items SET read=1 WHERE title=?", title)
	return err
}
