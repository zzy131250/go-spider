package v2

import "time"

var (
	timeout        = time.Second * 3
	detailUrlLabel = map[string]string{
		"label": "article link",
		"attr":  "href",
	}
	contentLabel = map[string]string{
		"label": "article header h2",
	}
)
