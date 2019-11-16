package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

type movie struct {
	id string
	name string
}

func main() {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)

	allMovie := []movie{}

	// 只抓取 網域為 "douban" 的, 並且併發數設定為5
	c.Limit(&colly.LimitRule{DomainGlob:  "*.douban.*", Parallelism: 5})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML(".hd", func (e *colly.HTMLElement) {
		var mId, mName string
		mId = strings.Split(e.ChildAttr("a", "href"), "/")[4]
		mName = strings.TrimSpace(e.DOM.Find("span.title").Eq(0).Text())
		log.Println(mId, mName)
		allMovie = append(allMovie, movie{
			id: mId,
			name: mName,
		})
	})

	c.OnHTML(".paginator a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	c.Visit("https://movie.douban.com/top250?start=0&filter=")

	c.Wait()
	fmt.Println(allMovie)
}
