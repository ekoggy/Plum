package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"strconv"
	"strings"
)

func CollectInfoFromDarknet() error {
	//creating collector
	c := colly.NewCollector(colly.AllowURLRevisit())

	//setting proxy (we need the TOR-proxy)
	prx, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:9050")
	if err != nil {
		return err
	}
	c.SetProxyFunc(prx)

	//getting number of last record
	update, err := checkUpdate(c)
	if err != nil {
		return err
	}
	fmt.Println("last record")
	//getting last N = 10 record
	err = getSomeRecords(c, update, 10)
	if err != nil {
		return err
	}
	fmt.Println("init table")
	lastUpdate := update

	//checking updates
	for {
		update, err = checkUpdate(c)
		if lastUpdate != update {
			lastUpdate = update
			err = getSomeRecords(c, update, 1)
		}
		if err != nil {
			return err
		}
	}
}

func getSomeRecords(c *colly.Collector, startRecord int, amount int) error {
	var err error
	for i := 0; i < amount; i++ {
		// may be gorutine
		var rec Record
		err = scrubMainPage(c, &rec, startRecord-i)
		if err != nil {
			return err
		}
		err = scrubBuyPage(c, &rec, startRecord-i)
		if err != nil {
			return err
		}

		_, err = insert(rec.Name, rec.Size, rec.Date, rec.Price, rec.Buy, rec.Source)
		fmt.Println(rec)
	}
	return nil
}

func scrubMainPage(c *colly.Collector, rec *Record, i int) error {
	var leakPage = "http://aby6efzmp7jzbwgidgqc6ghxi2vwpo6d7eaood5xuoxutrfofsmzcjqd.onion/page.php?pid="
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
		return err
	}
	return nil
}

func scrubBuyPage(c *colly.Collector, rec *Record, i int) error {
	var leakBuy = "http://aby6efzmp7jzbwgidgqc6ghxi2vwpo6d7eaood5xuoxutrfofsmzcjqd.onion/buy.php?db="
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
		return err
	}
	return nil
}

func checkUpdate(c *colly.Collector) (int, error) {
	lastPost := 0
	key := 0
	var leakDB = "http://aby6efzmp7jzbwgidgqc6ghxi2vwpo6d7eaood5xuoxutrfofsmzcjqd.onion/"
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if key == 3 {
			lastPost, _ = strconv.Atoi(strings.Split(e.Attr("href"), "=")[1])
		}
		key++
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(leakDB)
	if err != nil {
		return -1, err
	}
	fmt.Println("hello workd")
	return lastPost, nil
}
