package v1

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
	"strings"
	"time"
)

// parse web page with xpath
func parse(body io.ReadCloser, xpath string, results chan string) {
	defer body.Close()
	doc, err := htmlquery.Parse(body)
	if err != nil {
		fmt.Println("parse body error")
		return
	}
	for _, result := range htmlquery.Find(doc, xpath) {
		results <- strings.TrimSpace(htmlquery.InnerText(result))
	}
}

// get web page via http
func crawl(xpath string, ch chan string) chan string {
	results := make(chan string)
	go func() {
		for u := range ch {
			go func(url string) {
				fmt.Printf("start crawl url: %s\n", url)
				resp, err := http.Get(url)
				if err != nil {
					fmt.Printf("crawl url error: %s\n", url)
					return
				}
				go parse(resp.Body, xpath, results)
			}(u)
		}
	}()
	return results
}

// 爬虫v1版：
// 1. 指定列表页的详情url xpath，以及详情页的内容xpath进行爬取
// 2. 通过select超时进行任务的停止
func Crawler() {
	ch := make(chan string)
	// add page urls
	go func() {
		ch <- "http://www.ziazhou.com/"
		ch <- "http://www.ziazhou.com/page/2/"
	}()
	// crawl detail urls
	urls := crawl(detailUrlXpath, ch)
	// crawl contents
	contents := crawl(contentXpath, urls)
	for {
		select {
		case content := <-contents:
			fmt.Printf("find content: %s\n", content)
		case <-time.After(time.Second * 10):
			return
		}
	}
}
