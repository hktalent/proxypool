package getter

import (
	"github.com/hktalent/htmlquery"
	"github.com/hktalent/proxypool/pkg/models"
	"log"
	clog "unknwon.dev/clog/v2"
)

// PZZQZ get ip from http://pzzqz.com/
// https://pzzqz.com/?socks5=on&transparent=on&anonymous=on&ping=3000&country=all&ports=
// {"socks5":"on","transparent":"on","anonymous":"on","ping":"3000","country":"all","ports":""}
func PZZQZ() (result []*models.IP) {
	pollURL := "http://pzzqz.com/"
	if doc, err := htmlquery.LoadURLWithPost(pollURL, `{"socks5":"on","transparent":"on","anonymous":"on","ping":"3000","country":"all","ports":""}`, true); nil != err {
		log.Println("pzzqz", err)
	} else {
		trNode := htmlquery.Find(doc, "//table[@class='table table-hover']//tbody//tr")
		for i := 0; i < len(trNode); i++ {
			tdNode := htmlquery.Find(trNode[i], "//td")
			ip := htmlquery.InnerText(tdNode[0])
			port := htmlquery.InnerText(tdNode[1])
			Type := htmlquery.InnerText(tdNode[4])

			IP := models.NewIP()
			IP.Data = ip + ":" + port
			IP.Type1 = Type
			IP.Source = "pzzqz"
			result = append(result, IP)
		}
	}
	clog.Info("[pzzqz] done")
	return
}
