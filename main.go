package main

import (
	"go-ip-proxy/logger"
	"go-ip-proxy/parser"
	"go.uber.org/zap"
	"sync"
	"time"
)

func main() {
	log := logger.Config("./logs/all.log", "info", true)
	log.Info("test log", zap.Int("line", 47))

	configs := parser.NewConfig("./config/proxyWebsiteConfig.json")
	run(configs)

}

func run(configs *parser.Configs) {
	for {
		var wg sync.WaitGroup

		for _, configuration := range configs.Configs {
			parser.NewSelector(configuration)
		}

		wg.Wait()
		logger.Debug("finish once, sleep 10 minutes.")
		time.Sleep(time.Minute * 10)
	}
}

// func main() {
//     // 实例化
//     c := colly.NewCollector(
//         // 限定域名
//         colly.AllowedDomains("blog.phpha.com"),
//         // 最大深度
//         colly.MaxDepth(1),
//     )

//     // On every a element which has href attribute call callback
//     c.OnHTML("a[href]", func(e *colly.HTMLElement) {
//         link := e.Attr("href")
//         // Print link
//         fmt.Printf("Link found: %q -> %s\n", e.Text, link)
//         // Visit link found on page
//         // Only those links are visited which are in AllowedDomains
//         c.Visit(e.Request.AbsoluteURL(link))
//     })

//     // Before making a request print "Visiting ..."
//     c.OnRequest(func(r *colly.Request) {
//         fmt.Println("Visiting", r.URL.String())
//     })

//     // Start scraping on http://blog.phpha.com/
//     c.Visit("http://blog.phpha.com/")
// }
