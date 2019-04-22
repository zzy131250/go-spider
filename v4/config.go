package v4

import (
	"errors"
	"time"
)

const (
	minDelay = 50
	maxDelay = 1500
	chanSize = 100
)

var (
	timeout = time.Second * 15

	// 自定义error变量
	ErrProxyMayNotWork = errors.New("proxy may not work")

	userAgent = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	}
	// 参考https://github.com/jhao104/proxy_pool，仅供测试
	proxyPoolUrl = "http://118.24.52.95:5010"

	detailUrlLabel = map[string]string{
		"label": "article link",
		"attr":  "href",
	}
	contentLabel = map[string]string{
		"label": "article header h2",
	}
)
