package getter

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hktalent/proxypool/pkg/models"
	"github.com/parnurzeal/gorequest"
	clog "unknwon.dev/clog/v2"
)

// GBJ get ip from goubanjia.com
func GBJ() (result []*models.IP) {
	pollURL := "http://www.goubanjia.com/"

	resp, _, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}
	fmt.Println(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(errs)
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}

	doc.Find("body > div.container > div.section-header > div.row > div.container-fluid > div.row-fluid > div.span12 > tbody > tr").Each(func(_ int, s *goquery.Selection) {
		sf, _ := s.Find(".ip").Html()
		tee := regexp.MustCompile("<pstyle=\"display:none;\">.?.?</p>").ReplaceAllString(strings.Replace(sf, " ", "", -1), "")
		re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
		ip := models.NewIP()
		ip.Data = re.ReplaceAllString(tee, "")
		ip.Type1 = s.Find("td:nth-child(3) > a").Text()
		ip.Source = "goubanjia"
		clog.Info("ip.Data = %s , ip.Type = %s\n", ip.Data, ip.Type1)
		result = append(result, ip)
	})

	clog.Info("GBJ done.")
	return
}
