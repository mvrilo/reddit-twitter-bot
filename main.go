package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var twitter *anaconda.TwitterApi
var titles []string
var subreddit *string
var hot *bool

type datum struct {
	ID        string
	Title     string
	URL       string
	Permalink string
}

type response struct {
	Data struct {
		Children []struct {
			Data datum
		}
	}
}

func alreadyStored(d datum) bool {
	for _, title := range titles {
		if title == d.Title {
			return true
		}
	}
	return false
}

func store(d datum) {
	if len(titles) > 100 {
		titles = append(titles[1:], d.Title)
		return
	}
	titles = append(titles, d.Title)
}

func get(subreddit string) (data []datum) {
	kind := "new"
	if *hot {
		kind = "hot"
	}
	res, err := http.Get("https://www.reddit.com/r/" + subreddit + "/" + kind + ".json?limit=100")
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var resp response
	json.Unmarshal(body, &resp)
	for _, d := range resp.Data.Children {
		if alreadyStored(d.Data) {
			continue
		}
		store(d.Data)
		data = append(data, d.Data)
	}
	return
}

func fetch(populate bool, subreddit string) {
	for _, d := range get(subreddit) {
		if populate {
			return
		}
		t := fmt.Sprintf("%s %s", d.Title, d.URL)
		if _, err := twitter.PostTweet(t, nil); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	key := flag.String("key", "", "Twitter API consumer key")
	secret := flag.String("secret", "", "Twitter API consumer secret")
	accessToken := flag.String("access_token", "", "Twitter access key")
	accessSecret := flag.String("access_secret", "", "Twitter access secret")
	subreddit = flag.String("subreddit", "", "Subreddit to watch")
	timer := flag.Int("time", 30, "Time in seconds to fetch for posts")
	hot = flag.Bool("hot", false, "Lookup for hot posts instead of new")
	flag.Parse()

	if *subreddit == "" {
		log.Fatal("The subreddit is required")
	}

	if *key != "" && *secret != "" {
		anaconda.SetConsumerKey(*key)
		anaconda.SetConsumerSecret(*secret)
	} else {
		log.Fatal("You need to authorize Twitter API by passing the consumer key/secret")
	}

	if *accessToken != "" && *accessSecret != "" {
		twitter = anaconda.NewTwitterApi(*accessToken, *accessSecret)
	} else {
		log.Fatal("You need to authorize Twitter API by passing the access key/secret")
	}

	go fetch(true, *subreddit)
	log.Printf("reddit-twitter-bot started, watching subreddit: %s every %d seconds\n", *subreddit, *timer)
	for {
		select {
		case <-time.Tick(time.Duration(*timer) * time.Second):
			go fetch(false, *subreddit)
		}
	}
}
