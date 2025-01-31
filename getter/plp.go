package getter

import (
	clog "unknwon.dev/clog/v2"

	"github.com/hktalent/htmlquery"
	"github.com/hktalent/proxypool/pkg/models"
)

//PLP get ip from proxylistplus.com
func PLP() (result []*models.IP) {
	pollURL := "https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-1"
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

		IP.Source = "plp"
		clog.Info("[PLP] ip.Data = %s,ip.Type = %s,%s", IP.Data, IP.Type1, IP.Type2)

		result = append(result, IP)
	}

	clog.Info("PLP done.")
	return
}
