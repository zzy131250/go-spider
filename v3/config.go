package v3

import "time"

var (
	minDelay       = 50
	maxDelay       = 1500
	timeout        = time.Second * 5
	detailUrlLabel = map[string]string{
		"label": "article link",
		"attr":  "href",
	}
	contentLabel = map[string]string{
		"label": "article header h2",
	}
)
