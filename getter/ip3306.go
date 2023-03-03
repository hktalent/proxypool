package getter

import (
	"fmt"
	"log"
	clog "unknwon.dev/clog/v2"

	"github.com/hktalent/htmlquery"
	"github.com/hktalent/proxypool/pkg/models"
)

// IP3306 get ip from http://www.ip3366.net/
func IP3306() (result []*models.IP) {
	clog.Info("[IP3306]] start Get IpProxy")
	for i := 1; i < 7; i++ {
		pollURL := fmt.Sprintf("http://www.ip3366.net/free/?stype=3&page=%d", i)
		if doc, err := htmlquery.LoadURL(pollURL); nil != err {
			log.Println("ip3366", err)
		} else {
			trNode := htmlquery.Find(doc, "//div[@id='list']//table//tbody//tr")
			clog.Info("[IP3306] start up")

			//debug begin
			clog.Info("[IP3306] len(trNode) = %d ", len(trNode))
			for i := 1; i < len(trNode); i++ {
				tdNode := htmlquery.Find(trNode[i], "//td")
				ip := htmlquery.InnerText(tdNode[0])
				port := htmlquery.InnerText(tdNode[1])
				Type := htmlquery.InnerText(tdNode[3])
				speed := htmlquery.InnerText(tdNode[5])

				IP := models.NewIP()
				IP.Data = ip + ":" + port

				if Type == "HTTPS" {
					IP.Type1 = "https"
					IP.Type2 = ""

				} else if Type == "HTTP" {
					IP.Type1 = "http"
				}
				IP.Source = "ip3366.net"
				IP.Speed = extractSpeed(speed)

				clog.Info("[IP3306] ip.Data = %s,ip.Type = %s,%s ip.Speed = %d", IP.Data, IP.Type1, IP.Type2, IP.Speed)

				result = append(result, IP)
			}

			clog.Info("IP3306 done.")
		}
	}
	return
}
