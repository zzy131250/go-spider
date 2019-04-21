package v4

import "time"

var (
	minDelay = 50
	maxDelay = 1500
	timeout  = time.Second * 15

	userAgent = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	}
	proxyPoolUrl = "http://118.24.52.95:5010"

	detailUrlLabel = map[string]string{
		"label": "article link",
		"attr":  "href",
	}
	contentLabel = map[string]string{
		"label": "article header h2",
	}
)
