package getter

import (
	"fmt"
	"log"
	"strconv"

	clog "unknwon.dev/clog/v2"

	"github.com/PuerkitoBio/goquery"
	"github.com/hktalent/proxypool/pkg/models"
	"github.com/parnurzeal/gorequest"
)

//Data5u is not work now
// Data5u get ip from data5u.com
func Data5u() (result []*models.IP) {
	pollURL := "http://www.data5u.com/free/index.shtml"
	resp, _, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}
	if resp.StatusCode != 200 {
		log.Println(errs)
		return
	}
	fmt.Println(resp.Body)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}

	doc.Find("body > div.wlist > ul > li:nth-child(2) > ul").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		ss := s.Find("ul:nth-child(" + node + ") > span:nth-child(1) > li").Text()
		sss := s.Find("ul:nth-child(" + node + ") > span:nth-child(2) > li").Text()
		ssss := s.Find("ul:nth-child(" + node + ") > span:nth-child(4) > li").Text()
		ip := models.NewIP()
		ip.Data = ss + ":" + sss
		ip.Type1 = ssss
		ip.Source = "data5u"
		fmt.Printf("ip.Data = %s, ip.Type = %s", ip.Data, ip.Type1)
		result = append(result, ip)
	})
	clog.Info("Data5u done.")
	return
}
