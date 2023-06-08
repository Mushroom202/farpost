package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type good struct {
	date  string
	name  string
	price string
	user  string
	about string
	url   string
	view  string
	time  string
}

var (
	mainurl       string = "https://www.farpost.ru/vladivostok/sport/"
	scraperUrlApi        = "http://api.scraperapi.com?api_key=72b06c509284c5689e1cd3bab9a7a1a7&url="
)

// api + main doman
var url string = fmt.Sprintf("%v%v", scraperUrlApi, mainurl)

type FakeBrowserHeadersResponse struct {
	Result []map[string]string `json:"result"`
}

func RandomHeader(headersList []map[string]string) map[string]string {
	randomIndex := rand.Intn(len(headersList))
	return headersList[randomIndex]
}

func GetHeadersList() []map[string]string {

	// ScrapeOps Browser Headers API Endpint
	scrapeopsAPIKey := "3e3c6bc6-0b4e-40c4-a6b0-a54e3df02823"
	scrapeopsAPIEndpoint := "http://headers.scrapeops.io/v1/browser-headers?api_key=" + scrapeopsAPIKey

	req, _ := http.NewRequest("GET", scrapeopsAPIEndpoint, nil)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make Request
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()

		// Convert Body To JSON
		var fakeBrowserHeadersResponse FakeBrowserHeadersResponse
		json.NewDecoder(resp.Body).Decode(&fakeBrowserHeadersResponse)
		return fakeBrowserHeadersResponse.Result
	}

	var emptySlice []map[string]string
	return emptySlice
}
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
func farpostGetPages() string {
	var total string
	//var total string
	headersList := GetHeadersList()
	//initialision of collector
	c := colly.NewCollector()
	//Navigate on site pages
	c.OnRequest(func(r *colly.Request) {
		randomHeader := RandomHeader(headersList)
		for key, value := range randomHeader {
			r.Headers.Set(key, value)
		}
	})
	//code of connect to the site
	c.OnHTML("td#bulletins", func(e *colly.HTMLElement) {
		total = strings.Split(strings.Replace(strings.Split(strings.Replace(strings.TrimSpace(e.ChildText("div.pager")), "	", "", -1), "\n")[13], " ", "", 1), " ")[0]
	})
	c.Visit(fmt.Sprintf("%v", url))
	return total
}
func farpostGetUrls(total int) map[string]string {
	var listofviews map[string]string = map[string]string{}
	//var total string
	headersList := GetHeadersList()
	//initialision of collector
	c := colly.NewCollector()
	//Navigate on site pages
	c.OnRequest(func(r *colly.Request) {
		randomHeader := RandomHeader(headersList)
		for key, value := range randomHeader {
			r.Headers.Set(key, value)
		}
		fmt.Println("Scraping:", r.URL)
	})

	//code of connect to the site
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status:", r.StatusCode)

	})

	c.OnHTML("td#bulletins", func(e *colly.HTMLElement) {
		//
		e.ForEach("div.descriptionCell", func(_ int, e *colly.HTMLElement) {
			listofviews[fmt.Sprintf("%v%v", strings.TrimSuffix(url, "/vladivostok/sport/"), e.ChildAttr("a.bulletinLink", "href"))] = e.ChildText("span.views")
		})

	})
	for i := 1; i <= 1 && i <= total; i++ {
		c.Visit(fmt.Sprintf("%v?page=%v", url, i))
	}
	return listofviews //listofviews
}

// не получает данные со страницы
func farpostGetCard(url, view string) {
	var data good
	now := time.Now()
	headersList := GetHeadersList()
	//initialision of collector
	c := colly.NewCollector()
	//Navigate on site pages
	c.OnRequest(func(r *colly.Request) {
		randomHeader := RandomHeader(headersList)
		for key, value := range randomHeader {
			r.Headers.Set(key, value)
		}
		fmt.Println("Scraping:", r.URL)
	})
	//code of connect to the site
	c.OnHTML("td#bulletins", func(e *colly.HTMLElement) {
		data.about = e.ChildText("p.inplace")
		data.date = e.ChildText("div.viewbull-actual-date")
		data.name = e.ChildText("span.inplace")
		data.price = e.ChildAttr("span.data-bulletin-price", "data-bulletin-price")
		data.url = url
		data.user = e.ChildAttr("span.userNick > a", "href")
		data.view = view
		data.time = now.Format("2006/01/02 15:04")
		fmt.Println(data)
	})
	c.Visit(fmt.Sprintf("%v", url))

}
func main() {
	total := farpostGetPages()
	pages, err := strconv.Atoi(total)
	if err != nil {
		log.Fatal(err)
	}
	urlView := farpostGetUrls(pages)
	for url, view := range urlView {
		farpostGetCard(url, view)
	}

}
