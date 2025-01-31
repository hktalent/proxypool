package getter

import (
	"github.com/hktalent/htmlquery"
	"github.com/hktalent/proxypool/pkg/models"
	"golang.org/x/net/html"
	"log"
	clog "unknwon.dev/clog/v2"
)

// FQDL get ip from https://www.fanqieip.com/
func FQDL() (result []*models.IP) {
	pollURL := "https://www.fanqieip.com/free/1"
	if doc, err := htmlquery.LoadURL(pollURL); nil != err {
		log.Println("fanqieip", err)
	} else {
		trNode := htmlquery.Find(doc, "//table[@class='layui-table']//tbody//tr")
		for i := 0; i < len(trNode); i++ {
			tdNode := htmlquery.Find(trNode[i], "//td")
			ip := extractTextFromDivNode(tdNode[0])
			port := extractTextFromDivNode(tdNode[1])
			Type := "http"
			speed := htmlquery.InnerText(tdNode[4])

			IP := models.NewIP()
			IP.Data = ip + ":" + port
			IP.Type1 = Type
			IP.Source = "fanqieip"
			IP.Speed = extractSpeed(speed)
			result = append(result, IP)
		}

		clog.Info("[fanqiedl] done")
	}
	return
}

func extractTextFromDivNode(node *html.Node) string {
	divNode := htmlquery.Find(node, "//div")
	divOut := htmlquery.InnerText(divNode[0])
	return divOut
}
