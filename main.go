package main

import (
	"fmt"
	"go-ip-proxy/logger"
	"go-ip-proxy/parser"
	"go-ip-proxy/storage"
	"go.uber.org/zap"
	"sync"
	"time"
)

func main() {
	log := logger.Config("./logs/all.log", "info", true)
	log.Info("test log", zap.Int("line", 47))

	configs := parser.NewConfig("./config/proxyWebsiteConfig.json")

	// Load database.
	database, err := storage.NewStorage()
	defer database.Close()
	if err != nil {
		logger.Error("db error", zap.Error(err))
		panic(err)
	}

	run(database, configs)

}

func run(storage storage.Storage, configs *parser.Configs) {

	for {
		var wg sync.WaitGroup

		for _, configuration := range configs.Configs {
			items := parser.NewSelector(configuration)
			for _, item := range items {
				err := storage.Create(item, "1")
				if err != nil {
					logger.Error("db error", zap.Error(err))
				}
			}
		}
		for _, item := range storage.GetAll() {
			fmt.Printf("%s\n", string(item))
		}

		wg.Wait()
		logger.Debug("finish once, sleep 10 minutes.")
		time.Sleep(time.Second * 5)
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
