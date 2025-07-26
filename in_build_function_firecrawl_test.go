package swarmgo

import (
	"fmt"
	"testing"
)

const (
	appkey  = ""
	baseurl = ""
)

func TestFireCrawlFun(t *testing.T) {

	// AppKey: fc-b4cc1afecab64a1290be172596498738
	// BaseUrl: https://api.firecrawl.dev
	fn := NewFirecrawlFunction(appkey, baseurl)
	if fn == nil {
		t.Error("Firecrawl function is nil")
	}

	args := make(map[string]interface{})

	args["url"] = "https://www.google.com"
	result := fn.Work(args, nil)

	fmt.Println(result)
}
