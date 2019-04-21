package v4

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// 代理池
var proxyPool sync.Pool

func init() {
	proxyPool = sync.Pool{
		New: func() interface{} {
			return getProxyClient()
		},
	}
}

func proxyGet(url string) (resp *http.Response, err error) {
	var requestSuccess bool

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Connection", "keep-alive")
	rand.Seed(time.Now().UnixNano())
	request.Header.Set("User-Agent", userAgent[rand.Intn(len(userAgent))])
	client := proxyPool.Get().(*http.Client)

	for i := 0; i < 5; i++ {
		resp, err = client.Do(request)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("proxy may not work")
		}
		if resp.StatusCode == http.StatusOK {
			requestSuccess = true
			proxyPool.Put(client)
			break
		}
	}

	if !requestSuccess {
		fmt.Println("proxy may not work")
		return nil, errors.New("proxy may not work")
	}
	return
}

func getProxyClient() *http.Client {
	resp, err := http.Get(fmt.Sprintf("%s/get/", proxyPoolUrl))
	if err != nil {
		fmt.Println("get proxy url err")
		return nil
	}
	if resp == nil {
		fmt.Println("proxy url is nil")
		return nil
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	//var builder strings.Builder
	//builder.WriteString("http://")
	//builder.Write(data)
	var buf bytes.Buffer
	buf.WriteString("http://")
	buf.Write(data)
	fmt.Printf("proxy url: %s\n", buf.String())

	return &http.Client{
		Timeout: timeout, // 15s
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: timeout,
			}).DialContext,
			Proxy: func(request *http.Request) (i *url.URL, e error) {
				return url.Parse(buf.String())
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
