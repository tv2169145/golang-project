package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"sort"
	"strings"
	"time"
)

type item struct {
	StoryURL  string
	Source    string
	comments  string
	CrawledAt time.Time
	Comments  string
	Title     string
}

type wordInfo struct {
	word string
	count int
}

func main() {
	start := time.Now()
	items := []item{}
	c := colly.NewCollector(
			colly.AllowedDomains("old.reddit.com"),
			colly.Async(true),
		)

	c.OnHTML(".top-matter", func(e *colly.HTMLElement) {
		storyUrl := e.ChildAttr("a[data-event-action=title]", "href")
		source := "https://old.reddit.com/r/programming/"
		title := e.ChildText("a[data-event-action=title]")
		commitUrl := e.ChildAttr("a[data-event-action=commits]", "href")
		crawlAt := time.Now()
		items = append(items, item{
			StoryURL: storyUrl,
			Source: source,
			Title: title,
			Comments: commitUrl,
			CrawledAt: crawlAt,
		})
	})

	c.OnHTML("span.next-button", func(e *colly.HTMLElement) {
		nextUrl := e.ChildAttr("a", "href")
		c.Visit(nextUrl)
	})

	c.Limit(&colly.LimitRule{
		Parallelism: 8,
		RandomDelay: 2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	reddits := os.Args[1:]
	for _, reddit := range reddits {
		c.Visit(reddit)
	}
	c.Wait()
	//fmt.Println(items)
	fmt.Println(len(items))
	wordCount(items)
	fmt.Println("time spend:", time.Since(start).Seconds())
}

func wordCount(allItems []item) {
	wordMap := map[string]int{}

	for _, v := range allItems {
		titleStrings := strings.Split(strings.ToLower(v.Title), " ")
		for _, string := range titleStrings {
			wordMap[string]++
		}
	}
	allWordCount := []wordInfo{}
	for w, c := range wordMap {
		allWordCount = append(allWordCount, wordInfo{w, c})
	}

	sort.Slice(allWordCount, func(i, j int) bool {
		return allWordCount[i].count > allWordCount[j].count
	})

	for rank, count := range allWordCount {
		fmt.Println("No", rank+1, " word: ", count.word, " count: ", count.count)
		if rank >= 10 {
			break
		}
	}
}
