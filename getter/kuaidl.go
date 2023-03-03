package getter

import (
	"fmt"
	"github.com/henson/proxypool/pkg/models"
	"github.com/hktalent/htmlquery"
	"log"
	clog "unknwon.dev/clog/v2"
)

// KDL get ip from kuaidaili.com
func KDL() (result []*models.IP) {
	for i := 1; i < 5044; i++ {
		// https://www.kuaidaili.com/free/inha/5044/
		pollURL := fmt.Sprintf("http://www.kuaidaili.com/free/inha/%d/", i)
		if doc, err := htmlquery.LoadURL(pollURL); nil != err {
			log.Println("kuaidaili", err)
		} else {
			trNode := htmlquery.Find(doc, "//table[@class='table table-bordered table-striped']//tbody//tr")
			for i := 0; i < len(trNode); i++ {
				tdNode := htmlquery.Find(trNode[i], "//td")
				ip := htmlquery.InnerText(tdNode[0])
				port := htmlquery.InnerText(tdNode[1])
				Type := htmlquery.InnerText(tdNode[3])
				speed := htmlquery.InnerText(tdNode[5])

				IP := models.NewIP()
				IP.Data = ip + ":" + port
				if Type == "HTTPS" {
					IP.Type1 = "https"
					IP.Type2 = "https"
				} else if Type == "HTTP" {
					IP.Type1 = "http"
				}
				IP.Source = "KDL"
				IP.Speed = extractSpeed(speed)
				result = append(result, IP)
			}

			clog.Info("[kuaidaili] done")
		}
	}
	return
}
