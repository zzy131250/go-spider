package v4

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"io"
	"math/rand"
	"strings"
	"time"
)

// parse web page with goquery
func parse(ctx context.Context, body io.ReadCloser, label map[string]string, results chan string) {
	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println(errors.Wrap(err, "parse body error"))
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

// 限速器，按照随机时延访问
func limiter(limitFunc func()) {
	rand.Seed(time.Now().UnixNano())
	delay := minDelay + rand.Intn(maxDelay-minDelay)
	select {
	case <-time.After(time.Millisecond * time.Duration(delay)):
		limitFunc()
	}
}

// get web page via http
func crawl(ctx context.Context, url string, label map[string]string, sources chan string, results chan string) {
	select {
	case <-ctx.Done():
		fmt.Println("request timeout, stop http request")
		return
	default:
	}
	limiter(func() {
		fmt.Printf("start crawl url: %s\n", url)
		resp, err := proxyGet(url)
		if err != nil && err == ErrProxyMayNotWork {
			// 重新放入channel
			fmt.Printf("proxy not work, put url back into channel: %s\n", url)
			time.Sleep(time.Second)
			sources <- url
			return
		}
		if err != nil {
			fmt.Println(errors.Wrap(err, "crawl url error"))
			return
		}
		go parse(ctx, resp.Body, label, results)
	})
}

func crawlWithContext(parent context.Context, label map[string]string, ch chan string) chan string {
	results := make(chan string, chanSize)
	go func() {
		for url := range ch {
			// set timeout context for each request, 30s
			ctx, _ := context.WithTimeout(parent, timeout*2)
			go crawl(ctx, url, label, ch, results)
		}
	}()
	return results
}

// 爬虫v4版：
// 1. 添加代理和http头
// 2. 改进错误处理
// 3. url爬取失败重试
func Crawler() {
	ch := make(chan string, chanSize)
	ctx := context.Background()
	// add page urls
	go func() {
		ch <- "http://www.ziazhou.com/"
		ch <- "http://www.ziazhou.com/page/2/"
	}()
	// crawl detail urls
	urls := crawlWithContext(ctx, detailUrlLabel, ch)
	// crawl contents
	contents := crawlWithContext(ctx, contentLabel, urls)
	for {
		select {
		case content := <-contents:
			fmt.Printf("find content: %s\n", content)
		case <-time.After(timeout * 4): // 60s
			return
		}
	}
}
