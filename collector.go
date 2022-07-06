package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"log"
	"strconv"
	"strings"
)

func SetLocalProxy(c *colly.Collector) {
	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:9050")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)
}

func CollectPageDBInfo(i int) {
	c := colly.NewCollector(colly.AllowURLRevisit())
	c.SetRequestTimeout(0)
	SetLocalProxy(c)
	var leakPage string = "http://aby6efzmp7jzbwgidgqc6ghxi2vwpo6d7eaood5xuoxutrfofsmzcjqd.onion/page.php?pid="
	rec.Source = leakPage + strconv.Itoa(i)
	c.OnHTML(".single", func(e *colly.HTMLElement) {
		rec.Name = e.ChildText("h1")
	})

	c.OnHTML(".content", func(e *colly.HTMLElement) {
		contentString := e.DOM.Find("p:first-child").Text()
		data := strings.Split(contentString, " / ")[0]
		priceT := strings.Split(contentString, " / ")[2]
		price := strings.Split(priceT, "P")[1]
		cost := price[7:]
		rec.Date = data
		rec.Price = cost + "$"
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	err := c.Visit(leakPage + strconv.Itoa(i))
	if err != nil {
		return
	}
}

func CollectBuyDBInfo(i int) {
	c := colly.NewCollector(colly.AllowURLRevisit())
	c.SetRequestTimeout(0)
	SetLocalProxy(c)
	var leakBuy string = "http://aby6efzmp7jzbwgidgqc6ghxi2vwpo6d7eaood5xuoxutrfofsmzcjqd.onion/buy.php?db="
	rec.Buy = leakBuy + strconv.Itoa(i)
	c.OnHTML(".content", func(e *colly.HTMLElement) {
		contentString := e.Text
		fileSizeT := strings.Split(contentString, ": ")[2]
		fileSize := strings.Split(fileSizeT, "B")[0]
		rec.Size = fileSize + "B"
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(leakBuy + strconv.Itoa(i))
	if err != nil {
		return
	}
}

func CheckMaxNumber() int {
	c := colly.NewCollector(colly.AllowURLRevisit())
	c.SetRequestTimeout(0)
	SetLocalProxy(c)
	key := 0
	maxNumb := 0
	var leakDB string = "http://aby6efzmp7jzbwgidgqc6ghxi2vwpo6d7eaood5xuoxutrfofsmzcjqd.onion/"
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if key == 3 {
			maxNumb, _ = strconv.Atoi(strings.Split(e.Attr("href"), "=")[1])
		}
		key++
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(leakDB)
	if err != nil {
		return -1
	}
	return maxNumb
}
