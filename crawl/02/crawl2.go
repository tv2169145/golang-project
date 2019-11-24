package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"))
	allMove := []string{}

	c.Limit(&colly.LimitRule{DomainGlob: "*.edwardmovieclub.*", Parallelism: 5})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML("h2", func(e *colly.HTMLElement) {
		mName := e.ChildText("span")
		//fmt.Println(mName)
		allMove = append(allMove, mName)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit("https://edwardmovieclub.com/imdb-top-10-2019/")

	c.Wait()
	fmt.Println(allMove)
}
