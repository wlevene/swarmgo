package swarmgo

import (
	"fmt"

	"github.com/mendableai/firecrawl-go"
)

type (
	FirecrawlFunction struct {
		BaseFunction
		app     *firecrawl.FirecrawlApp
		appkey  string
		baseurl string
	}
)

func NewFirecrawlFunction(appkey, baseurl string) *FirecrawlFunction {

	app, err := firecrawl.NewFirecrawlApp(appkey, baseurl)

	if err != nil {
		fmt.Println("Error creating FirecrawlApp:", err)
		return nil
	}

	fn := &FirecrawlFunction{
		appkey:  appkey,
		baseurl: baseurl,
		app:     app,
	}

	baseFn, err := NewCustomFunction(fn)
	if err != nil {
		return nil
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(fn.work)
	return fn
}

var _ AgentFunction = (*FirecrawlFunction)(nil)

func (fn *FirecrawlFunction) work(args map[string]interface{}, contextVariables map[string]interface{}) Result {

	result := Result{}

	var url string
	if val, ok := args["url"]; ok {
		url = val.(string)
	}

	scrapeResult, err := fn.app.ScrapeURL(url,
		&firecrawl.ScrapeParams{
			Formats: []string{"markdown"},
		})

	if err != nil {
		fmt.Println("Error crawling URL:", err)
		return result
	}

	result.Data = scrapeResult.Markdown
	result.Success = true
	return result
}

func (fn *FirecrawlFunction) GetName() string {
	return "scrape_websites_to_markdown"
}

func (fn *FirecrawlFunction) GetDescription() string {
	return "Crawl and scrape websites and return content in clean llm-ready markdown."
}

func (fn *FirecrawlFunction) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"url": map[string]interface{}{
				"type":        "string",
				"description": "The URL to scrape",
			},
		},
		"required": []interface{}{"url"},
	}
}
