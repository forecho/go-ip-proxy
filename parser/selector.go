package parser

import (
	"fmt"
	"github.com/gocolly/colly"
	"go-ip-proxy/logger"
	"go-ip-proxy/util"
	"go.uber.org/zap"
	"strings"
)

type Selector struct {
	configuration *Config
	currentUrl    string
	currentIndex  int
	urls          []string
	selectorMap   map[string][]string
}

var nameValue = make(map[string]string)
var proxy = make([]string, 0)

func NewSelector(config *Config) (err error) {

	parameters := strings.Split(config.UrlParameters, ",")
	urls := util.MakeUrls(config.UrlFormat, parameters)

	for _, url := range urls {
		itemMap := getPageBody(config, url)
		for _, value := range itemMap {
			logger.Infof("ip is %s", value)
		}
	}
	return
}

func getPageBody(config *Config, url string) []string {
	selectorMap := make(map[string][]string)
	for _, value := range config.ValueRuleMap {
		if value.Name == "" || value.Rule == "" {
			logger.Errorf("config name:%s contains valueRuleMap item with empty name or rule, this item will be ignored.", config.Name)
			continue
		}

		if value.Name == "table" {
			selectorMap[value.Name] = []string{value.Rule}
		} else if value.Attr != "" {
			selectorMap[value.Name] = []string{value.Rule, value.Attr}
		} else {
			selectorMap[value.Name] = []string{value.Rule}
		}
	}

	c := colly.NewCollector()

	c.OnHTML(selectorMap["table"][0], func(element *colly.HTMLElement) {
		//fmt.Printf("%s\n", element.DOM.Find(selectorMap["ip"][0]).Text())

		for key, value := range selectorMap {
			if key != "table" {
				//fmt.Printf(element.DOM.Find(value[0]).Text() + "\n")
				nameValue[key] = element.DOM.Find(value[0]).Text()
				logger.Info(nameValue[key])
				//fmt.Printf("%s\n", element.DOM.Find(value[0]).Text())
			}
		}
		proxy = append(
			proxy,
			fmt.Sprintf("%s://%s:%s", strings.ToLower(nameValue["type"]), nameValue["ip"], nameValue["port"]))
	})

	c.UserAgent = util.RandomUA()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
	})
	//c.OnResponse(func(resp *colly.Response) {
	//	body = string(resp.Body)
	//})

	c.OnError(func(resp *colly.Response, errHttp error) {
		logger.Error("response error", zap.Error(errHttp))
	})

	err := c.Visit(url)

	if err != nil {
		logger.Error("visit error", zap.Error(err))
	}
	logger.Info("visit " + url)

	return proxy
}
