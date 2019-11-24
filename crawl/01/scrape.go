package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"sort"
	"strconv"
)

type movie struct {
	No string	`json:"no"`
	Name string	`json:"name"`
}

type wordInfo struct {
	word string
	count int
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

	c.OnHTML(".item", func (e *colly.HTMLElement) {
		var no string
		var name string
		no = e.ChildText("em")
		name = e.DOM.Find(".hd .title").Eq(0).Text()
		allMovie = append(allMovie, movie{
			no, name,
		})
		//var mId, mName string
		//mId = strings.Split(e.ChildAttr("a", "href"), "/")[4]
		//mName = strings.TrimSpace(e.DOM.Find("span.title").Eq(0).Text())
		//log.Println(mId, mName)
		//allMovie = append(allMovie, movie{
		//	id: mId,
		//	name: mName,
		//})
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
	sortByNo(allMovie)
	//dumpAllMovies(allMovie)
	//wordCount(allMovie)

}

func dumpAllMovies(movies []movie) {
	bss, err := json.MarshalIndent(movies, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bss))
}

func wordCount(movies []movie) {
	wordMap := map[string]int{}

	for _, v := range movies {
		wordMap[v.Name]++
	}
	wis := []wordInfo{}
	for w, c := range wordMap {
		wis = append(wis, wordInfo{
			word: w,
			count: c,
		})
	}
	fmt.Println(wis)
}

func sortByNo(allMovies []movie) {
	sort.Slice(allMovies, func(i, j int) bool {
		is, err := strconv.Atoi(allMovies[i].No)
		if err != nil {
			panic(err)
		}
		js, err := strconv.Atoi(allMovies[j].No)
		if err != nil {
			panic(err)
		}
		return is < js
		//return strconv.Atoi(allMovies[i].No) < strconv.Atoi(allMovies[j].No)
	})
	fmt.Println(allMovies)
	fmt.Println(len(allMovies))
}
