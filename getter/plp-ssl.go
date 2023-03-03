package getter

import (
	"fmt"
	clog "unknwon.dev/clog/v2"

	"github.com/henson/proxypool/pkg/models"
	"github.com/hktalent/htmlquery"
)

// PLPSSL get ip from proxylistplus.com
func PLPSSL() (result []*models.IP) {
	for i := 1; i < 6; i++ {
		pollURL := fmt.Sprintf("https://list.proxylistplus.com/SSL-List-%d", i)
		doc, _ := htmlquery.LoadURL(pollURL)
		trNode := htmlquery.Find(doc, "//div[@class='hfeed site']//table[@class='bg']//tbody//tr")

		for i := 3; i < len(trNode); i++ {
			tdNode := htmlquery.Find(trNode[i], "//td")
			ip := htmlquery.InnerText(tdNode[1])
			port := htmlquery.InnerText(tdNode[2])
			Type := htmlquery.InnerText(tdNode[6])

			IP := models.NewIP()
			IP.Data = ip + ":" + port

			if Type == "yes" {
				IP.Type1 = "https"
				IP.Type2 = ""

			} else if Type == "no" {
				IP.Type1 = "http"
			}
			IP.Source = "plp-ssl"

			clog.Info("[PLP SSL] ip.Data = %s,ip.Type = %s,%s", IP.Data, IP.Type1, IP.Type2)

			result = append(result, IP)
		}
	}
	clog.Info("PLP SSL done.")
	return
}
