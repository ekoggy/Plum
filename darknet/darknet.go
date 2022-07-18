package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"strconv"
	"strings"
	"time"
)

func CollectInfoFromDarknet() error {

	//getting number of last record
	update, err := checkUpdate()
	if err != nil {
		return err
	}
	//getting last N = 10 record
	err = getSomeRecords(update, 10)
	if err != nil {
		return err
	}

	//checking updates
	lastUpdate := update
	for {
		update, err = checkUpdate()
		if lastUpdate != update {
			lastUpdate = update
			err = getSomeRecords(update, 1)
		}
		if err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
	}
}

func getSomeRecords(startRecord int, amount int) error {
	var err error
	for i := 0; i < amount; i++ {
		// may be gorutine
		var rec Record
		err = scrubMainPage(&rec, startRecord-i)
		if err != nil {
			return err
		}
		err = scrubBuyPage(&rec, startRecord-i)
		if err != nil {
			return err
		}
		if rec.Name == "" {
			amount++
			continue
		}
		_, err = insert(rec.Name, rec.Size, rec.Date, rec.Price, rec.Buy, rec.Source)
	}
	return nil
}

func scrubMainPage(rec *Record, i int) error {
	//creating collector
	c := colly.NewCollector(colly.AllowURLRevisit())

	//setting proxy (we need the TOR-proxy)
	prx, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:9050")
	if err != nil {
		return err
	}
	c.SetProxyFunc(prx)

	c.SetRequestTimeout(0)

	var leakPage = "http://aby6efzmp7jzbwgidgqc6ghxi2vwpo6d7eaood5xuoxutrfofsmzcjqd.onion/page.php?pid="
	rec.Source = leakPage + strconv.Itoa(i)

	c.OnHTML(".single", func(e *colly.HTMLElement) {
		rec.Name = e.ChildText("h1")
	})

	c.OnHTML(".content", func(e *colly.HTMLElement) {
		contentString := e.DOM.Find("p:first-child").Text()
		data := strings.Split(contentString, " / ")[0]
		price := strings.Split(strings.Split(contentString, " / ")[2], "P")[1]
		cost := price[7:]
		rec.Date = data
		rec.Price = cost + "$"
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	err = c.Visit(rec.Source)
	if err != nil {
		return err
	}
	return nil
}

func scrubBuyPage(rec *Record, i int) error {
	//creating collector
	c := colly.NewCollector(colly.AllowURLRevisit())

	//setting proxy (we need the TOR-proxy)
	prx, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:9050")
	if err != nil {
		return err
	}
	c.SetProxyFunc(prx)

	c.SetRequestTimeout(0)

	var leakBuy = "http://aby6efzmp7jzbwgidgqc6ghxi2vwpo6d7eaood5xuoxutrfofsmzcjqd.onion/buy.php?db="
	rec.Buy = leakBuy + strconv.Itoa(i)

	c.OnHTML(".content", func(e *colly.HTMLElement) {
		rec.Size = strings.Split(strings.Split(e.Text, ": ")[2], "B")[0] + "B"
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err = c.Visit(leakBuy + strconv.Itoa(i))

	if err != nil {
		return err
	}
	return nil
}

func checkUpdate() (int, error) {
	//creating collector
	c := colly.NewCollector(colly.AllowURLRevisit())

	//setting proxy (we need the TOR-proxy)
	prx, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:9050")
	if err != nil {
		return -1, err
	}
	c.SetProxyFunc(prx)

	c.SetRequestTimeout(0)

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

	err = c.Visit(leakDB)
	if err != nil {
		return -1, err
	}
	return lastPost, nil
}
