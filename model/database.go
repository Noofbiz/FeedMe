package model

import (
	"database/sql"
	"errors"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
)

// Database Schema for assets/feeds.db
// *****************************************************************************
// all_feeds
// | name | url | updated |
// *****************************************************************************
// settings
// | updates_every | remove_after |
// *****************************************************************************
// feed_items
// | origin | title | author_name | author_email | published | content | description | link |

var (
	ResetUpdateTicker = make(chan struct{})
	ResetRemoveTicker = make(chan struct{})
)

type feedDatabase struct {
	db *sql.DB
}

var TheFeedData feedDatabase

func (f *feedDatabase) Open() error {
	var err error
	var genStartTable bool
	if _, err := os.Stat("assets/feeds.db"); os.IsNotExist(err) {
		genStartTable = true
	}

	f.db, err = sql.Open("sqlite3", "assets/feeds.db")
	if err != nil {
		return err
	}
	f.db.SetMaxOpenConns(1)

	if genStartTable {
		f.db.Exec(`
      CREATE TABLE all_feeds (name TEXT, url TEXT, updated TEXT);
      DELETE FROM all_feeds;
      `)
		f.db.Exec(`
      CREATE TABLE settings (updates_every INTEGER, remove_after INTEGER);
      DELETE FROM settings;
      INSERT INTO settings(updates_every, remove_after) values(12, 30);
      `)
		f.db.Exec(`
      CREATE TABLE feed_items (origin TEXT, title TEXT, author_name TEXT, author_email TEXT, published TEXT, content TEXT, description TEXT, link TEXT);
      DELETE FROM feed_items;
      `)
	}

	return nil
}

func (f *feedDatabase) Close() error {
	err := f.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (f *feedDatabase) AddFeed(url string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	tx, err := f.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO all_feeds(name, url, updated) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	stmt.Exec(feed.Title, url, time.Time{}.String())
	tx.Commit()

	if err = f.UpdateFeeds(); err != nil {
		return err
	}

	return nil
}

func (f *feedDatabase) UpdateFeeds() error {
	allFeedsRows, err := f.db.Query("select * from all_feeds")
	defer allFeedsRows.Close()
	if err != nil {
		return err
	}
	allFeeds := []struct {
		name, url, updated string
	}{}
	for allFeedsRows.Next() {
		row := struct {
			name, url, updated string
		}{}
		err = allFeedsRows.Scan(&row.name, &row.url, &row.updated)
		if err != nil {
			return err
		}
		allFeeds = append(allFeeds, row)
	}

	fp := gofeed.NewParser()
	feedTx, err := f.db.Begin()
	if err != nil {
		return err
	}
	feedUpdateSTMT, err := feedTx.Prepare("UPDATE all_feeds SET updated=? WHERE name=?")
	if err != nil {
		return err
	}
	itemSTMT, err := feedTx.Prepare("INSERT INTO feed_items(origin, title, author_name, author_email, published, content, description, link) values(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer itemSTMT.Close()
	for _, fd := range allFeeds {
		if _, err = feedUpdateSTMT.Exec(time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST"), fd.name); err != nil {
			return err
		}
		feed, err := fp.ParseURL(fd.url)
		if err != nil {
			return err
		}
		for _, i := range feed.Items {
			t1, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", fd.updated)
			if err != nil {
				return err
			}
			t2 := time.Time{}
			if i.PublishedParsed != nil {
				t2 = *i.PublishedParsed
			} else if i.UpdatedParsed != nil {
				t2 = *i.UpdatedParsed
			}
			var author_name, author_email string
			if i.Author != nil {
				author_name = i.Author.Name
				author_email = i.Author.Email
			}
			if t2.After(t1) {
				if _, err = itemSTMT.Exec(fd.name, i.Title, author_name, author_email, t2.Format("2006-01-02 15:04:05.999999999 -0700 MST"), i.Content, i.Description, i.Link); err != nil {
					return err
				}
			}
		}
	}
	if err = feedTx.Commit(); err != nil {
		return err
	}

	return nil
}

func (f *feedDatabase) RemoveFeed(url string) error {
	var name string
	r, err := f.db.Query("SELECT name FROM all_feeds WHERE url=?", url)
	if err != nil {
		return err
	}
	if !r.Next() {
		return errors.New("No rows returned with url: " + url)
	}
	err = r.Scan(&name)
	if err != nil {
		return err
	}
	r.Close()
	_, err = f.db.Exec("DELETE FROM all_feeds WHERE url=?", url)
	if err != nil {
		return err
	}
	_, err = f.db.Exec("DELETE FROM feed_items WHERE origin=?", name)
	if err != nil {
		return err
	}
	return nil
}

func (f *feedDatabase) SetUpdatesEvery(n int) error {
	_, err := f.db.Exec("UPDATE settings SET updates_every=?", n)
	ResetUpdateTicker <- struct{}{}
	return err
}

func (f *feedDatabase) GetUpdatesEvery() (time.Duration, error) {
	r, err := f.db.Query("SELECT updates_every FROM settings")
	if err != nil {
		return time.Duration(0), err
	}
	r.Next()
	var t int
	err = r.Scan(&t)
	r.Close()
	if err != nil {
		return time.Duration(0), err
	}
	return time.Hour * time.Duration(t), nil
}

func (f *feedDatabase) SetRemoveAfter(n int) error {
	_, err := f.db.Exec("UPDATE settings SET remove_after=?", n)
	ResetRemoveTicker <- struct{}{}
	return err
}

func (f *feedDatabase) GetRemoveAfter() (time.Duration, error) {
	r, err := f.db.Query("SELECT remove_after FROM settings")
	if err != nil {
		return time.Duration(0), err
	}
	r.Next()
	var t int
	err = r.Scan(&t)
	r.Close()
	if err != nil {
		return time.Duration(0), err
	}
	return time.Hour * 24 * time.Duration(t), nil
}

func (f *feedDatabase) PruneFeeds() error {
	titles := make([]string, 0)
	dur, err := f.GetRemoveAfter()
	r, err := f.db.Query("SELECT title, published FROM feed_items")
	if err != nil {
		return err
	}
	for r.Next() {
		var title string
		var pub string
		r.Scan(&title, &pub)
		pubTime, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", pub)
		if time.Now().Sub(pubTime).Nanoseconds() > dur.Nanoseconds() {
			titles = append(titles, title)
		}
	}
	r.Close()

	for _, t := range titles {
		_, err := f.db.Exec("DELETE FROM feed_items WHERE title=?", t)
		if err != nil {
			return err
		}
	}
	return nil
}
