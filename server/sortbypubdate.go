package server

import "github.com/Noofbiz/FeedMe/model"

type sortByPubDate []model.ItemInfo

func (s sortByPubDate) Len() int {
	return len(s)
}

func (s sortByPubDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByPubDate) Less(i, j int) bool {
	return s[i].Published.After(s[j].Published)
}
