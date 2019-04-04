package v2

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"time"
)

// parse web page with goquery
func parse(ctx context.Context, body io.ReadCloser, label map[string]string, results chan string) {
	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println("parse body error")
		return
	}
	doc.Find(label["label"]).Each(func(i int, selection *goquery.Selection) {
		if attr, ok := label["attr"]; ok {
			// select attribute
			result, _ := selection.Attr(attr)
			select {
			case <-ctx.Done():
				fmt.Println("request timeout, stop parse")
			default:
				results <- strings.TrimSpace(result)
			}
		} else {
			// select text
			select {
			case <-ctx.Done():
				fmt.Println("request timeout, stop parse")
			default:
				results <- strings.TrimSpace(selection.Text())
			}
		}
	})
}

// get web page via http
func crawl(parent context.Context, label map[string]string, ch chan string) chan string {
	results := make(chan string)
	// set timeout context
	ctx, _ := context.WithTimeout(parent, timeout)
	go func() {
		for u := range ch {
			go func(url string) {
				select {
				case <-ctx.Done():
					fmt.Println("request timeout, stop http request")
					return
				default:
				}
				fmt.Printf("start crawl url: %s\n", url)
				resp, err := http.Get(url)
				if err != nil {
					fmt.Printf("crawl url error: %s\n", url)
					return
				}
				go parse(ctx, resp.Body, label, results)
			}(u)
		}
	}()
	return results
}

// 爬虫v2版：
// 1. html解析器改为goquery，并使用label筛选，需要指出取attr还是text
// 2. 添加context，限制每个请求的处理时间
func Crawler() {
	ch := make(chan string)
	ctx := context.Background()
	// add page urls
	go func() {
		ch <- "http://www.ziazhou.com/"
		ch <- "http://www.ziazhou.com/page/2/"
	}()
	// crawl detail urls
	urls := crawl(ctx, detailUrlLabel, ch)
	// crawl contents
	contents := crawl(ctx, contentLabel, urls)
	for {
		select {
		case content := <-contents:
			fmt.Printf("find content: %s\n", content)
		case <-time.After(time.Second * 10):
			return
		}
	}
}
