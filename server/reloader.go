package server

import (
	"time"

	"github.com/Noofbiz/FeedMe/model"
)

func autoReloader() {
	model.TheFeedData.UpdateFeeds()
	model.TheFeedData.PruneFeeds()
	go func() {
		u, _ := model.TheFeedData.GetUpdatesEvery()
		updateTicker := time.NewTicker(u)
		r, _ := model.TheFeedData.GetRemoveAfter()
		removeTicker := time.NewTicker(r)
		for {
			select {
			case <-updateTicker.C:
				model.TheFeedData.UpdateFeeds()
			case <-removeTicker.C:
				model.TheFeedData.PruneFeeds()
			case <-model.ResetUpdateTicker:
				u, _ := model.TheFeedData.GetUpdatesEvery()
				updateTicker.Stop()
				updateTicker = time.NewTicker(u)
			case <-model.ResetRemoveTicker:
				r, _ := model.TheFeedData.GetRemoveAfter()
				removeTicker.Stop()
				removeTicker = time.NewTicker(r)
			}
		}
	}()
}
