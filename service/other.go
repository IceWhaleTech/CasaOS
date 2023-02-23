package service

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type OtherService interface {
	Search(key string) ([]model.SearchEngine, error)
	AgentSearch(url string) ([]byte, error)
}

type otherService struct{}

func (s *otherService) Search(key string) ([]model.SearchEngine, error) {

	engines := []model.SearchEngine{}
	engines = append(engines, model.SearchEngine{
		Name:      "bing",
		Icon:      "https://files.codelife.cc/itab/search/bing.svg",
		SearchUrl: "https://www.bing.com/search?q=",
		RecoUrl:   "https://www.bing.com/osjson.aspx?query=", // + keyword
	}, model.SearchEngine{
		Name:      "google",
		Icon:      "https://files.codelife.cc/itab/search/google.svg",
		SearchUrl: "https://www.google.com/search?q=",
		RecoUrl:   "https://www.google.com/complete/search?client=gws-wiz&xssi=t&hl=en-US&authuser=0&dpr=1&q=", // + keyword
	}, model.SearchEngine{
		Name:      "baidu",
		Icon:      "https://files.codelife.cc/itab/search/baidu.svg",
		SearchUrl: "https://www.baidu.com/s?wd=",
		RecoUrl:   "https://www.baidu.com/sugrec?json=1&prod=pc&wd=", // + keyword
	}, model.SearchEngine{
		Name:      "duckduckgo",
		Icon:      "https://files.codelife.cc/itab/search/duckduckgo.svg",
		SearchUrl: "https://duckduckgo.com/?q=",
		RecoUrl:   "https://duckduckgo.com/ac/?type=list&q=", // + keyword
	}, model.SearchEngine{
		Name:      "startpage",
		Icon:      "https://www.startpage.com/sp/cdn/favicons/apple-touch-icon-60x60--default.png",
		SearchUrl: "https://www.startpage.com/do/search?q=",
		RecoUrl:   "https://www.startpage.com/suggestions?segment=startpage.udog&lui=english&q=", // + keyword
	})

	client := resty.New()
	client.SetTimeout(3 * time.Second) // 设置全局超时时间
	var wg sync.WaitGroup
	for i := 0; i < len(engines); i++ {
		wg.Add(1)
		go func(i int, k string) {
			name := engines[i].Name
			url := engines[i].RecoUrl + url.QueryEscape(k)
			defer wg.Done()
			resp, err := client.R().Get(url)
			if err != nil {
				logger.Error("Then get search result error: %v", zap.Error(err), zap.String("name", name), zap.String("url", url))
				return
			}
			res := []string{}
			if name == "bing" {
				r := gjson.Get(resp.String(), "1")
				for _, v := range r.Array() {
					res = append(res, v.String())
				}
			} else if name == "google" {
				r := gjson.Get(strings.Replace(resp.String(), ")]}'", "", 1), "0.#.0")
				for _, v := range r.Array() {
					res = append(res, strings.ReplaceAll(strings.ReplaceAll(v.String(), "<b>", " "), "</b>", ""))
				}
			} else if name == "baidu" {
				r := gjson.Get(resp.String(), "g.#.q")
				for _, v := range r.Array() {
					res = append(res, v.String())
				}
			} else if name == "duckduckgo" {
				r := gjson.Get(resp.String(), "1")
				for _, v := range r.Array() {
					res = append(res, v.String())
				}
			} else if name == "startpage" {
				r := gjson.Get(resp.String(), "suggestions.#.text")
				for _, v := range r.Array() {
					res = append(res, v.String())
				}
			}
			engines[i].Data = res
		}(i, key)
	}
	wg.Wait()

	return engines, nil

}

func (s *otherService) AgentSearch(url string) ([]byte, error) {
	client := resty.New()
	client.SetTimeout(3 * time.Second) // 设置全局超时时间
	resp, err := client.R().Get(url)
	if err != nil {
		logger.Error("Then get search result error: %v", zap.Error(err), zap.String("url", url))
		return nil, err
	}
	return resp.Body(), nil
}

func NewOtherService() OtherService {
	return &otherService{}
}
