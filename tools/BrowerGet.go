package tools

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// BrowerGet 通过URL获取网页内容并返回文本
func BrowerGetDocument(url string) string {
	c := colly.NewCollector() // Colly OnHTML 回调 :contentReference[oaicite:6]{index=6}
	result := ""
	// // 2. 抓取"干员档案"内容
	c.OnHTML("span.mw-headline[id='干员档案']", func(e *colly.HTMLElement) {
		// e.DOM 是 *goquery.Selection
		e.DOM.Parent().
			NextUntil("h2"). // 取到下一个 h2 前的所有节点 :contentReference[oaicite:7]{index=7}
			Each(func(_ int, s *goquery.Selection) {
				text := strings.TrimSpace(s.Text())
				result += text
			})
	})

	// 4. 发起请求
	if err := c.Visit(url); err != nil {
		log.Fatal(err)
	}

	return result
}
