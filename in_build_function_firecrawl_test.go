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

	fn := NewFirecrawlFunction(appkey, baseurl)
	if fn == nil {
		t.Error("Firecrawl function is nil")
	}

	args := make(map[string]interface{})

	args["url"] = "https://www.google.com"
	result := fn.Work(args, nil)

	fmt.Println(result)
}
