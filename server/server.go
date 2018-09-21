package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"strconv"

	"github.com/Noofbiz/FeedMe/model"
)

var tmpl = template.Must(template.ParseGlob("assets/gohtml/*"))

func StartServer() string {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/save", saveSettingsHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/static/"))))

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	go http.Serve(l, nil)
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	feedData, err := model.TheFeedData.GetFeedData(r.FormValue("show-feed"), "")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sort.Sort(sortByPubDate(feedData.Items))

	err = tmpl.ExecuteTemplate(w, "index", feedData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveSettingsHandler(w http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	settingsChanges := struct {
		RemovedFeeds []string
		AddedFeeds   []string
		UpdatesEvery string
		ExpiresAfter string
	}{}
	err = json.Unmarshal(d, &settingsChanges)
	if err != nil {
		log.Println(err)
		http.Error(w, "malformed json data sent to server", http.StatusInternalServerError)
		return
	}
	for _, r := range settingsChanges.RemovedFeeds {
		err = model.TheFeedData.RemoveFeed(r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	for _, a := range settingsChanges.AddedFeeds {
		err = model.TheFeedData.AddFeed(a)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if i, err := strconv.Atoi(settingsChanges.UpdatesEvery); err == nil && i > 0 {
		err = model.TheFeedData.SetUpdatesEvery(i)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if i, err := strconv.Atoi(settingsChanges.ExpiresAfter); err == nil && i > 0 {
		err = model.TheFeedData.SetRemoveAfter(i)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
